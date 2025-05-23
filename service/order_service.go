package service

import (
	"context"
	"coresamples/common"
	"coresamples/dbutils"
	"coresamples/ent"
	"coresamples/ent/orderflag"
	"coresamples/ent/orderinfo"
	"coresamples/model"
	pb "coresamples/proto"
	"coresamples/publisher"
	"coresamples/util"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"golang.org/x/exp/slices"
)

const (
	OrderProcessing          = "order_processing"
	OrderCanceled            = "order_canceled_order"
	OrderRedraw              = "order_redraw_order"
	OrderMajorStatusCanceled = "canceled_order"
	OrderCompleted           = "order_completed"
	CategoryOrderStatus      = "order_status"
	CategoryTnpStatus        = "tnp_status"
	CategoryKitStatusDetails = "kit_status_details"
)

type IOrderService interface {
	GetOrder(orderId int, ctx context.Context) (*ent.OrderInfo, error)
	CancelOrder(orderId int, ctx context.Context) (*ent.OrderInfo, error)
	RestoreCanceledOrder(orderId int, ctx context.Context) (*ent.OrderInfo, error)
	AddOrderFlag(orderFlagReq *pb.AddOrderFlagRequest, ctx context.Context) (*ent.OrderFlag, error)
	FlagOrder(orderId int, orderFlagNames []string, ctx context.Context) (*ent.OrderInfo, error)
	// FlagOrderTx flag an order with given flag names using provided tx.
	// This function DOES NOT commit, but WILL roll back on error
	FlagOrderTx(orderId int, orderFlagNames []string, tx *ent.Tx, ctx context.Context) error
	FlagOrdersWithSampleId(sampleId int32, flags []string, ctx context.Context) ([]*ent.OrderInfo, error)
	UnflagOrder(orderId int, orderFlagNames []string, ctx context.Context) (*ent.OrderInfo, error)
	ListOrderFlagTypes(ctx context.Context) ([]*ent.OrderFlag, error)
	ChangeOrderStatus(orderId int, status string, ctx context.Context) (*ent.OrderInfo, error)
	GetOrderStatusForDisplay(ctx context.Context) ([]*ent.OrderFlag, []*ent.PatientFlag, error)
	RerunSampleTests(sampleId int, testIds []int32, ctx context.Context) (*pb.RerunSampleTestsResponse, error)
	RestoreOrderStatus(jwtToken string, sampleId string, ctx context.Context) (int, string, error)
	TriggerOrderTransmissionOnSampleReceiving(sampleId int32, tubeDetails []*model.SampleTubeDetails,
		receivedTime time.Time, isRedraw bool, ctx context.Context) ([]*model.OrderTransmissionResult, error)
	DispatchRemoveSampleOrder(sampleId int, tubeType string, receivedTime time.Time,
		collectionTime time.Time, isRedraw bool, ctx context.Context) error
	DispatchOrderWithReceivedTime(sampleId int, tubeType string,
		receivedTime time.Time, collectionTime time.Time, isRedraw bool, ctx context.Context) error
	UpdateOrderKitStatusByAccessionId(accessionId string, kitStatus string, ctx context.Context) error
}

type OrderService struct {
	Service
}

func NewOrderService(dbClient *ent.Client, redisClient *common.RedisClient) IOrderService {
	s := &OrderService{
		Service: InitService(dbClient, redisClient),
	}
	return s
}

func (s *OrderService) GetOrder(orderId int, ctx context.Context) (*ent.OrderInfo, error) {
	return dbutils.GetOrderById(orderId, s.dbClient, ctx)
}

func (s *OrderService) CancelOrder(orderId int, ctx context.Context) (*ent.OrderInfo, error) {
	return dbutils.CancelOrderById(orderId, s.dbClient, ctx)
}

func (s *OrderService) RestoreCanceledOrder(orderId int, ctx context.Context) (*ent.OrderInfo, error) {
	return dbutils.RestoreOrderById(orderId, s.dbClient, ctx)
}

func (s *OrderService) AddOrderFlag(orderFlagReq *pb.AddOrderFlagRequest, ctx context.Context) (*ent.OrderFlag, error) {
	var flag *ent.OrderFlag
	var err error
	if flag, err = dbutils.GetOrderFlagByName(orderFlagReq.OrderFlagName, s.dbClient, ctx); flag == nil || err != nil {
		flag, err = dbutils.CreateOrderFlag(orderFlagReq, s.dbClient, ctx)
	}
	if err != nil {
		return nil, err
	}
	return flag, err
}

// TODO: not sure transaction is used right, test this thoroughly
func (s *OrderService) FlagOrder(orderId int, orderFlagNames []string, ctx context.Context) (*ent.OrderInfo, error) {
	tx, err := s.dbClient.Tx(ctx)
	if err != nil {
		return nil, err
	}
	err = s.FlagOrderTx(orderId, orderFlagNames, tx, ctx)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, dbutils.Rollback(tx, err)
	}
	return dbutils.GetOrderById(orderId, s.dbClient, ctx)
}

func (s *OrderService) FlagOrdersWithSampleId(sampleId int32, flags []string, ctx context.Context) ([]*ent.OrderInfo, error) {
	var flagged []*ent.OrderInfo
	orders, err := dbutils.GetOrdersWithSampleId(int(sampleId), s.dbClient, ctx)
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		o, err := s.FlagOrder(order.ID, flags, ctx)
		if err != nil {
			common.Error(err)
			sentry.CaptureMessage(err.Error())
			continue
		}
		flagged = append(flagged, o)
	}
	return flagged, nil
}

// FlagOrderTx flag an order with given flag names using provided tx.
// This function DOES NOT commit, but WILL roll back on error
func (s *OrderService) FlagOrderTx(orderId int, orderFlagNames []string, tx *ent.Tx, ctx context.Context) error {
	order, err := tx.OrderInfo.Get(ctx, orderId)
	if err != nil {
		return dbutils.Rollback(tx, err)
	}
	orderUpdate := order.Update()
	for _, flagName := range orderFlagNames {
		flag, err := dbutils.GetOrderFlagByName(flagName, tx.Client(), ctx)
		if flag == nil || err != nil {
			continue
		}
		// Update the Primary Order Record, then update the Connection which is used as history
		switch flagName {
		case OrderProcessing:
			orderUpdate = orderUpdate.SetOrderProcessingTime(time.Now())
		case OrderCanceled:
			orderUpdate = orderUpdate.SetOrderCancelTime(time.Now()).SetOrderMajorStatus(OrderMajorStatusCanceled).SetOrderCanceled(true).SetOrderStatus(OrderCanceled)

		case OrderRedraw:
			orderUpdate = orderUpdate.SetOrderRedrawTime(time.Now())
		}

		orderUpdate = orderUpdate.SetOrderFlagged(true)
		orderUpdate = orderUpdate.AddOrderFlags(flag)
		order, err = orderUpdate.Save(ctx)
		if err != nil {
			return dbutils.Rollback(tx, err)
		}

		order, err = s.fillInOrderStatus(order, ctx)
		if err != nil {
			return dbutils.Rollback(tx, err)
		}
		orderUpdate = order.Update()
		// Order status update start
		if flag.OrderFlagCategory == CategoryKitStatusDetails {
			switch flag.OrderFlagName {
			case "kit_sample_shipped_back", "kit_lab_shipped_kit":
				orderUpdate = orderUpdate.
					SetOrderMinorStatus("in_transit").
					SetOrderMajorStatus("awaiting_sample")
			case "kit_patient_received_kit":
				orderUpdate = orderUpdate.
					SetOrderMinorStatus("kit_patient_received_kit").
					SetOrderMajorStatus("awaiting_sample")
			case "kit_lab_received":
				orderUpdate = orderUpdate.
					SetOrderMinorStatus("kit_lab_received").
					SetOrderMajorStatus("analyzing_sample")
			}
		}
		// Order status update stop

		order, err = orderUpdate.Save(ctx)
		if err != nil {
			return dbutils.Rollback(tx, err)
		}

		if flag.OrderFlagCategory == CategoryKitStatusDetails &&
			order.OrderKitStatus != "kit_lab_received" {
			order, err = order.Update().SetOrderKitStatus(flagName).Save(ctx)
		}
	}
	return nil
}

func (s *OrderService) UnflagOrder(orderId int, orderFlagNames []string, ctx context.Context) (*ent.OrderInfo, error) {
	tx, err := s.dbClient.Tx(ctx)
	if err != nil {
		return nil, err
	}
	order, err := tx.OrderInfo.Get(ctx, orderId)
	if err != nil {
		return nil, dbutils.Rollback(tx, err)
	}
	orderUpdate := order.Update()
	for _, flagName := range orderFlagNames {
		flag, err := dbutils.GetOrderFlagByName(flagName, tx.Client(), ctx)
		if err != nil || flag == nil {
			continue
		}
		orderUpdate = orderUpdate.RemoveOrderFlags(flag)
	}
	order, err = orderUpdate.Save(ctx)
	if err != nil {
		return nil, dbutils.Rollback(tx, err)
	}
	flags, err := order.QueryOrderFlags().All(ctx)
	if err != nil {
		return nil, dbutils.Rollback(tx, err)
	}
	if flags == nil || len(flags) == 0 {
		// unset flagged status
		order, err = order.Update().SetOrderFlagged(false).Save(ctx)
		if err != nil {
			return nil, dbutils.Rollback(tx, err)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, dbutils.Rollback(tx, err)
	}
	return order, nil
}

func (s *OrderService) ListOrderFlagTypes(ctx context.Context) ([]*ent.OrderFlag, error) {
	return dbutils.GetAllOrderFlags(s.dbClient, ctx)
}

func (s *OrderService) ChangeOrderStatus(orderId int, status string, ctx context.Context) (*ent.OrderInfo, error) {
	err := dbutils.UpdateOrderStatus(orderId, status, s.dbClient, ctx)
	if err != nil {
		return nil, err
	}
	return dbutils.GetOrderById(orderId, s.dbClient, ctx)
}

func (s *OrderService) GetOrderStatusForDisplay(ctx context.Context) ([]*ent.OrderFlag, []*ent.PatientFlag, error) {
	patientFlags, err := dbutils.GetAllPatientFlags(s.dbClient, ctx)
	if err != nil {
		return nil, nil, err
	}
	orderFlags, err := dbutils.GetAllOrderFlags(s.dbClient, ctx)
	if err != nil {
		return nil, nil, err
	}

	return orderFlags, patientFlags, nil
}

func (s *OrderService) RerunSampleTests(sampleId int, testIds []int32, ctx context.Context) (*pb.RerunSampleTestsResponse, error) {
	msg := &pb.OrderMessage{
		SampleId:         int32(sampleId),
		Action:           "lis",
		TestId:           testIds,
		IsRerun:          true,
		IsRedraw:         false,
		IsLabDirectOrder: false,
		Destination:      "DI",
	}
	err := publisher.GetPublisher().SendOrderMessage(msg)
	if err != nil {
		return &pb.RerunSampleTestsResponse{
			SendStatus: "failed",
			SendTime:   time.Now().Format("2006-01-02 15:04:05-0700"),
			SendLog:    err.Error(),
		}, err
	}
	return &pb.RerunSampleTestsResponse{
		SendStatus: "success",
		SendTime:   time.Now().Format("2006-01-02 15:04:05-0700"),
		SendLog:    fmt.Sprintf("%v", msg),
	}, nil
}

func (s *OrderService) RestoreOrderStatus(jwtToken string, sampleId string, ctx context.Context) (int, string, error) {
	var sample *ent.Sample
	var err error
	var claim *util.TokenClaim
	if claim, err = util.ParseJWTToken(jwtToken, common.Secrets.JWTSecret); err != nil {
		return 401, "invalid user token", err
	}
	if len(sampleId) < 10 {
		id, _ := strconv.Atoi(sampleId)
		sample, err = dbutils.GetSampleById(id, s.dbClient, ctx)
	} else {
		sample, err = dbutils.GetSampleByAccessionId(sampleId, s.dbClient, ctx)
	}
	if err != nil {
		return 404, "sample not found", err
	}

	order, err := s.dbClient.OrderInfo.Query().Where(
		orderinfo.ID(sample.OrderID),
		orderinfo.OrderStatusEQ(OrderCanceled),
	).WithOrderFlags(
		func(query *ent.OrderFlagQuery) {
			query.Where(
				orderflag.OrderFlagCategoryEQ("order_status"),
				orderflag.OrderFlagLevelLTE(4),
			).Limit(1)
		},
	).First(ctx)
	if err != nil {
		return 404, "order_canceled_order info not found", err
	}
	orderFlagName := order.Edges.OrderFlags[0].OrderFlagName
	common.RecordAuditLog(
		common.AuditLogEntry{
			EventID:             uuid.NewString(),
			ServiceName:         common.ServiceName,
			ServiceType:         "backend",
			EventName:           "restore_order_status",
			EntityType:          "order_info",
			EntityID:            strconv.Itoa(order.ID),
			AttributeValuePrior: "order_canceled_order",
			AttributeName:       "order_status",
			AttributeValuePost:  orderFlagName,
			User:                strconv.Itoa(claim.UserId),
			Entrypoint:          fmt.Sprintf("GRPC, order %d restored by user %d", order.ID, claim.UserId),
		})
	err = order.Update().SetOrderStatus(orderFlagName).Exec(ctx)
	if err != nil {
		return 500, fmt.Sprintf("error occured when restoring order status: %s", err.Error()), err
	}
	return 200, "successfully restored order status", nil
}

// TriggerOrderTransmissionOnSampleReceiving Send order to lab and trigger testing.
// The function needs to be executed in task queue to ensure the message being sent successfully,
// and also avoid blocking for too long
func (s *OrderService) TriggerOrderTransmissionOnSampleReceiving(sampleId int32, tubeDetails []*model.SampleTubeDetails,
	receivedTime time.Time, isRedraw bool, ctx context.Context) ([]*model.OrderTransmissionResult, error) {
	// bookkeeping for sent order tests
	var sentOrders []*model.OrderTransmissionResult
	var kafkaSendErr error
	kafkaSendErr = nil
	for _, detail := range tubeDetails {
		recordFound := true
		val, err := s.redisClient.Get(ctx, dbutils.KeyLabOrderSendRecord(int(sampleId), detail.TubeType)).Result()
		if err == nil {
			// found the record
			sendRecord := &ent.LabOrderSendHistory{}
			err = json.Unmarshal([]byte(val), sendRecord)
			if err != nil {
				return nil, err
			}
		} else {
			// fail to find the record in redis, try searching in the database
			_, err = dbutils.GetLabOrderSendRecord(int(sampleId), detail.TubeType, false, true, s.dbClient, ctx)
			if err != nil {
				// fail to find the record in both redis and database
				recordFound = false
			}
		}
		// if we find the record and it's not redraw order, do not send the order to DI
		// record it as sent
		if recordFound && !isRedraw {
			sentOrders = append(sentOrders, &model.OrderTransmissionResult{
				SampleId: int(sampleId),
				TubeType: detail.TubeType,
				Status:   "sent",
			})
			continue
		}
		// for order not yet sent or needs redraw
		err = s.DispatchOrderWithReceivedTime(int(sampleId), detail.TubeType, receivedTime, detail.CollectionTime, isRedraw, ctx)
		if err != nil {
			kafkaSendErr = err
			common.Error(err)
			continue
		}
		if detail.TubeType == "METAL_FREE_URINE" {
			err = s.DispatchOrderWithReceivedTime(int(sampleId), "URINE", receivedTime, detail.CollectionTime, isRedraw, ctx)
			if err != nil {
				common.Error(err)
				kafkaSendErr = err
				continue
			}
			dbutils.CreateLabOrderSendRecord(int(sampleId), "URINE", isRedraw, false, "+", s.dbClient, ctx)
		} else if detail.TubeType == "COVID19_STOOL" {
			err = s.DispatchOrderWithReceivedTime(int(sampleId), "COVID_STOOL", receivedTime, detail.CollectionTime, isRedraw, ctx)
			if err != nil {
				common.Error(err)
				kafkaSendErr = err
				continue
			}
			dbutils.CreateLabOrderSendRecord(int(sampleId), "COVID_STOOL_Triggered_By_COVID19_STOOL", isRedraw, false, "+", s.dbClient, ctx)
		}
		record, _ := dbutils.CreateLabOrderSendRecord(int(sampleId), detail.TubeType, isRedraw, false, "+", s.dbClient, ctx)
		recb, err := json.Marshal(record)
		if err != nil {
			common.Error(err)
			continue
		}
		s.redisClient.SetEX(ctx, dbutils.KeyLabOrderSendRecord(int(sampleId), detail.TubeType), recb, time.Second*1000)
	}
	// if some kafka sending fails return an error to trigger resend
	return sentOrders, kafkaSendErr
}

func (s *OrderService) DispatchRemoveSampleOrder(sampleId int, tubeType string, receivedTime time.Time,
	collectionTime time.Time, isRedraw bool, ctx context.Context) error {
	sampleType := util.TubeTypeToSampleType(tubeType)
	if tubeType == "SALIVA" {
		sampleType = "COVID_STOOL"
	}

	tests, err := s.getSampleTestIdsWithSampleType(int32(sampleId), sampleType, ctx)
	if err != nil {
		return err
	}
	msg := &pb.OrderMessage{
		SampleId:         int32(sampleId),
		Action:           "-",
		TestId:           tests,
		ReceiveTime:      receivedTime.Format("2006-01-02 15:04:05-0700"),
		CollectionTime:   collectionTime.Format("2006-01-02 15:04:05-0700"),
		IsRerun:          false,
		IsRedraw:         isRedraw,
		IsLabDirectOrder: false,
		Destination:      "DI",
	}
	err = publisher.GetPublisher().SendOrderMessage(msg)
	if err != nil {
		return err
	}
	if sampleType == "Serum" {
		customerIds := []int{999990, 999993, 999995, 3194}
		sample, err := dbutils.GetSampleById(sampleId, s.dbClient, ctx)
		if err != nil {
			return err
		}
		if !slices.Contains(customerIds, sample.CustomerID) {
			_, err := dbutils.CreateLabOrderSendRecord(sampleId, tubeType, isRedraw, true, "-", s.dbClient, ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// UpdateOrderKitStatusByAccessionId updates the order_kit_status for the given accession ID
func (s *OrderService) UpdateOrderKitStatusByAccessionId(accessionId string, kitStatus string, ctx context.Context) error {
	// Find the order by accession ID
	order, err := dbutils.GetOrderByAccessionId(accessionId, s.dbClient, ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch order: %w", err)
	}

	// Update the order_kit_status field
	err = dbutils.UpdateOrderKitStatus(order.ID, kitStatus, s.dbClient, ctx)
	if err != nil {
		return fmt.Errorf("failed to update order kit status: %w", err)
	}

	return nil
}

// private helper functions

// this function does not commit the transaction
func (s *OrderService) fillInOrderStatus(order *ent.OrderInfo, ctx context.Context) (*ent.OrderInfo, error) {
	flags, err := order.QueryOrderFlags().All(ctx)
	if err != nil {
		return order, err
	}
	orderUpdate := order.Update()
	categoryStatus := s.findHighestLevelOrderFlags(flags, ctx)
	if categoryStatus[CategoryOrderStatus] != "" && (order.OrderStatus != OrderCompleted ||
		categoryStatus[CategoryOrderStatus] == OrderCanceled ||
		categoryStatus[CategoryOrderStatus] == OrderRedraw ||
		categoryStatus[CategoryOrderStatus] == OrderCompleted) {
		orderUpdate = orderUpdate.SetOrderStatus(categoryStatus[CategoryOrderStatus])
	}

	if categoryStatus[CategoryTnpStatus] != "" {
		orderUpdate = orderUpdate.SetOrderTnpIssueStatus(categoryStatus[CategoryTnpStatus])
	}
	return orderUpdate.Save(ctx)
}

func (s *OrderService) findHighestLevelOrderFlags(flags []*ent.OrderFlag, ctx context.Context) map[string]string {
	res := make(map[string]string)
	highestLevelFlags := make(map[string]*ent.OrderFlag)
	for _, flag := range flags {
		category := flag.OrderFlagCategory
		_, exist := highestLevelFlags[category]
		if !exist || flag.OrderFlagLevel > highestLevelFlags[category].OrderFlagLevel {
			highestLevelFlags[category] = flag
			res[category] = flag.OrderFlagName
		}
	}
	return res
}

// dispatchOrderWithReceivedTime Dispatch the order to lab through kafka to trigger lab testing.
func (s *OrderService) DispatchOrderWithReceivedTime(sampleId int, tubeType string,
	receivedTime time.Time, collectionTime time.Time, isRedraw bool, ctx context.Context) error {
	sampleType := util.TubeTypeToSampleType(tubeType)
	if tubeType == "SALIVA" {
		//weird hard code
		sampleType = "COVID_STOOL"
	}

	testIds, err := s.getSampleTestIdsWithSampleType(int32(sampleId), sampleType, ctx)
	if err != nil {
		return err
	}
	orderMessage := &pb.OrderMessage{
		SampleId:         int32(sampleId),
		Action:           "+",
		TestId:           testIds,
		ReceiveTime:      receivedTime.Format("2006-01-02 15:04:05-0700"),
		CollectionTime:   collectionTime.Format("2006-01-02 15:04:05-0700"),
		IsRerun:          false,
		IsRedraw:         isRedraw,
		IsLabDirectOrder: false,
		Destination:      "DI",
	}
	err = publisher.GetPublisher().SendOrderMessage(orderMessage)
	if err != nil {
		return err
	}

	if sampleType == "Serum" {
		customerIds := []int{999990, 999993, 999995, 3194}
		sample, err := dbutils.GetSampleById(sampleId, s.dbClient, ctx)
		if err != nil {
			return err
		}
		// check customer first
		// TODO: this may only work when customer is in
		if !slices.Contains(customerIds, sample.CustomerID) {
			_, err = dbutils.CreateLabOrderSendRecord(sampleId, tubeType, isRedraw, true, "+", s.dbClient, ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *OrderService) getSampleTestIdsWithSampleType(sampleId int32, sampleType string, ctx context.Context) ([]int32, error) {
	var ans []int32
	samples, err := GetCurrentSampleService().GetSampleTests([]int32{sampleId}, ctx)
	if err != nil {
		return nil, err
	}
	if len(samples) == 0 {
		return ans, err
	}
	sample := samples[0]
	for _, test := range sample.Tests {
		if strings.ToUpper(test.TestType) == strings.ToUpper(sampleType) {
			ans = append(ans, test.TestId)
		}
	}
	return ans, nil
}
