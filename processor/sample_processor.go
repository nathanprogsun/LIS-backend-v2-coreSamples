package processor

import (
	"context"
	"fmt"
	"time"

	"coresamples/common"
	"coresamples/dbutils"
	"coresamples/ent"
	"coresamples/ent/labordersendhistory"
	"coresamples/ent/sample"
	"coresamples/ent/tuberequirement"
	pb "coresamples/proto"
	"coresamples/publisher"
	"coresamples/service"
	"coresamples/tasks"
	"coresamples/util"
	"encoding/json"

	"github.com/go-redsync/redsync/v4"
	"github.com/hibiken/asynq"
)

const (
	DeliverMethodSuppliedByClient = "suppliedByClient"
)

type SampleProcessor struct {
	Processor
	OrderService  service.IOrderService
	SampleService service.ISampleService
	rs            *redsync.Redsync
	TestService   service.ITestService
}

func NewSampleProcessor(dbClient *ent.Client, redisClient *common.RedisClient, rs *redsync.Redsync) *SampleProcessor {
	return &SampleProcessor{
		Processor:     InitProcessor(dbClient, redisClient, context.Background()),
		rs:            rs,
		OrderService:  service.GetCurrentOrderService(),
		SampleService: service.GetCurrentSampleService(),
	}
}

func (p *SampleProcessor) HandlePostSampleOrder(ctx context.Context, t *asynq.Task) error {
	task := &tasks.PostSampleOrderTask{}
	if err := json.Unmarshal(t.Payload(), task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	if common.Env.DryRun {
		return nil
	}
	return p.handlePostSampleOrder(task, ctx)
}

func (p *SampleProcessor) HandleFlagOrderOnReceiving(ctx context.Context, t *asynq.Task) error {
	task := &tasks.SampleTubeReceiveTask{}
	if err := json.Unmarshal(t.Payload(), task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	if common.Env.DryRun {
		return nil
	}
	_, err := p.OrderService.FlagOrdersWithSampleId(task.SampleId, []string{
		"kit_received",
		"order_processing",
		"report_pending",
		"kit_lab_received",
	}, ctx)
	return err
}

// HandleSendOrderOnReceiving Send order test to DI to trigger lab testing, this needs to be idempotent
// for each tube type by looking up in the lab order send log and corresponding redis
func (p *SampleProcessor) HandleSendOrderOnReceiving(ctx context.Context, t *asynq.Task) error {
	task := &tasks.SampleTubeReceiveTask{}
	if err := json.Unmarshal(t.Payload(), task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	if common.Env.DryRun {
		return nil
	}
	_, err := p.OrderService.TriggerOrderTransmissionOnSampleReceiving(task.SampleId, task.TubeDetails, task.ReceivedTime, task.IsRedraw, ctx)
	return err
}

func (p *SampleProcessor) HandleSendSampleReceiveGeneralEvent(ctx context.Context, t *asynq.Task) error {
	event := &pb.GeneralEvent{}
	if err := json.Unmarshal(t.Payload(), event); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	if common.Env.DryRun {
		return nil
	}
	return publisher.GetPublisher().SendGeneralEvent(event)
}

// v1_Zhibin: LIS-Kafka General events
// HandleGeneralEvent processes a general event task.
func (p *SampleProcessor) HandleSampleOrderGeneralEvent(ctx context.Context, t *asynq.Task) error {
	var task tasks.GeneralEventTask
	if err := json.Unmarshal(t.Payload(), &task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	if common.Env.DryRun {
		return nil
	}

	switch task.Event.EventProvider {
	//todo: issue system
	// case "lis-issue-system":
	// 	return p.handleIssueSystemEvent(task.Event)
	case "lis-order":
		return p.handleOrderEvent(task.Event, ctx)
	case "lis-accessioning":
		return p.handleAccessioningEvent(task.Event, ctx)
	case "lis-report":
		return p.handleReportEvent(task.Event, ctx)
	default:
		return nil
	}
}

func (p *SampleProcessor) HandleCancelOrderEvent(ctx context.Context, t *asynq.Task) error {
	var task tasks.CancelOrderTask
	if err := json.Unmarshal(t.Payload(), &task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	if common.Env.DryRun {
		return nil
	}
	sampleId := task.Event.GetSampleId()
	_, err := p.OrderService.FlagOrdersWithSampleId(sampleId, []string{
		"order_canceled_order",
	}, ctx)
	return err
}

func (p *SampleProcessor) HandleRedrawOrderEvent(ctx context.Context, t *asynq.Task) error {
	// todo: it has parse in func (eh *EventHandler) handleRedrawOrderEvent()
	var task tasks.RedrawOrderInfoTask
	if err := json.Unmarshal(t.Payload(), &task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	if task.Event == nil {
		return fmt.Errorf("task.Event is nil: %v", asynq.SkipRetry)
	}

	if common.Env.DryRun {
		return nil
	}

	if task.Event.GetType() != "insert" {
		return fmt.Errorf("task.Event.GetType() is not 'insert': %v", asynq.SkipRetry)
	}
	sampleID := task.Event.GetData().GetSampleId()

	// Query sample table by sampleID
	_, err := dbutils.GetSampleById(int(sampleID), p.dbClient, p.ctx)
	if err != nil {
		return fmt.Errorf("sampleID:%d is not in sample table, err:%w", sampleID, err)
	}

	// Fetch order info in order table by sampleID
	orderInfo, err := dbutils.FindOrderWithSampleId(sampleID, p.dbClient, p.ctx)
	if err != nil {
		return fmt.Errorf("failed to query order info by sampleID:%d, err:%w", sampleID, err)
	}

	// If report is ready, flag the order
	_, err = p.OrderService.FlagOrder(orderInfo.ID, []string{"order_redraw_order"}, ctx)
	if err != nil {
		return fmt.Errorf("failed to flag order: %w", err)
	}

	return nil
}

func (p *SampleProcessor) HandleEditOrderEvent(ctx context.Context, t *asynq.Task) error {
	var task tasks.EditOrderTask
	if err := json.Unmarshal(t.Payload(), &task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	if common.Env.DryRun {
		return nil
	}

	sampleId := task.Event.GetSampleId()
	_, err := dbutils.GetSampleById(int(sampleId), p.dbClient, p.ctx)
	if err != nil {
		return fmt.Errorf("sampleID:%d is not found, err:%w", sampleId, err)
	}

	// Fetch order info in order table by sampleID
	orderInfo, err := dbutils.FindOrderWithSampleId(sampleId, p.dbClient, p.ctx)
	if err != nil {
		return fmt.Errorf("failed to query order info by sampleID:%d, err:%w", sampleId, err)
	}

	if len(task.Event.AddOnTestList) > 0 {
		var testsToAdd []int
		for _, testID := range task.Event.AddOnTestList {
			testsToAdd = append(testsToAdd, int(testID))
		}

		// Update database: add tests
		err = p.dbClient.OrderInfo.UpdateOne(orderInfo).
			AddTestIDs(testsToAdd...).
			Exec(p.ctx)
		if err != nil {
			return fmt.Errorf("failed to add tests to order, err: %w", err)
		}

		// Clear Redis cache
		getSampleTestsRedisKey := dbutils.KeyGetSampleTests(int(sampleId))
		p.redisClient.Del(p.ctx, getSampleTestsRedisKey)
		getLabTestsRedisKey := dbutils.KeyLabTestsBySampleID(int(sampleId))
		p.redisClient.Del(p.ctx, getLabTestsRedisKey)

		var receivedTubeTypes []string
		orderSentHistory, err := p.dbClient.LabOrderSendHistory.Query().Where(
			labordersendhistory.SampleIDEQ(int(sampleId)),
		).All(p.ctx)
		if err != nil {
			return fmt.Errorf("failed to query lab_order_send_history for sampleID:%d, err: %w", sampleId, err)
		}

		if len(orderSentHistory) > 0 {
			for _, orderSent := range orderSentHistory {
				receivedTubeTypes = append(receivedTubeTypes, orderSent.TubeType)
			}
		}

		err = p.SendAddOnOrderToLIS(
			task.Event.AddOnTestList,
			receivedTubeTypes,
			"+",
			int(sampleId),
			p.ctx,
		)
		if err != nil {
			return fmt.Errorf("failed to send add-on order to LIS, err: %w", err)
		}
	}

	// Process Remove Tests
	if len(task.Event.DeleteTestList) > 0 {
		var testsToDelete []int
		for _, testID := range task.Event.DeleteTestList {
			testsToDelete = append(testsToDelete, int(testID))
		}

		// 🔹 Step 1: Remove Tests from Order
		err = p.dbClient.OrderInfo.UpdateOne(orderInfo).
			RemoveTestIDs(testsToDelete...).
			Exec(p.ctx)
		if err != nil {
			return fmt.Errorf("failed to remove tests from order, err: %w", err)
		}

		// 🔹 Step 2: Clear Redis Cache
		getSampleTestsRedisKey := dbutils.KeyGetSampleTests(int(sampleId))
		p.redisClient.Del(p.ctx, getSampleTestsRedisKey)
		getLabTestsRedisKey := dbutils.KeyLabTestsBySampleID(int(sampleId))
		p.redisClient.Del(p.ctx, getLabTestsRedisKey)

		// 🔹 Step 3: Fetch Order Sent History
		var receivedTubeTypes []string
		orderSentHistory, err := p.dbClient.LabOrderSendHistory.Query().Where(
			labordersendhistory.SampleIDEQ(int(sampleId)),
		).All(p.ctx)
		if err != nil {
			return fmt.Errorf("failed to query lab_order_send_history for sampleID:%d, err: %w", sampleId, err)
		}

		// 🔹 Step 4: Delete Tube Requirements if Order Was Sent
		if len(orderSentHistory) > 0 {
			_, err := p.dbClient.TubeRequirement.Delete().
				Where(tuberequirement.SampleIDEQ(int(sampleId))).
				Exec(p.ctx)
			if err != nil {
				return fmt.Errorf("failed to delete tube requirements for sampleID:%d, err: %w", sampleId, err)
			}

			// Extract received tube types
			for _, orderSent := range orderSentHistory {
				receivedTubeTypes = append(receivedTubeTypes, orderSent.TubeType)
			}
		}

		// 🔹 Step 5: Send Remove Order to LIS
		err = p.SendAddOnOrderToLIS(
			task.Event.DeleteTestList,
			receivedTubeTypes,
			"-",
			int(sampleId),
			p.ctx,
		)
		if err != nil {
			return fmt.Errorf("failed to send remove order to LIS, err: %w", err)
		}
	}

	// 🔹 Step 1: Check if `new_tube_info.noOfTubes` exists
	if task.Event.NewTubeInfo != nil && len(task.Event.NewTubeInfo.NoOfTubes) > 0 {
		// 🔹 Step 2: Delete old tube requirements
		_, err := p.dbClient.TubeRequirement.Delete().
			Where(tuberequirement.SampleIDEQ(int(sampleId))).
			Exec(p.ctx)
		if err != nil {
			return fmt.Errorf("failed to delete tube requirements for sampleID:%d, err: %w", sampleId, err)
		}

		// 🔹 Step 3: Iterate over `noOfTubes` and update counts
		for tubeType, requiredCount := range task.Event.NewTubeInfo.NoOfTubes {
			// 🔹 Step 4: Call `updateSampleRequiredTubeCount()`
			err := p.SampleService.UpdateSampleRequiredTubeCount(
				sampleId,
				tubeType,
				requiredCount,
				"Order Service Update",
				p.ctx,
			)
			if err != nil {
				return fmt.Errorf("failed to update required tube count for sample_id=%d, tube_type=%s, err: %w", sampleId, tubeType, err)
			}
		}
	}

	return nil
}

func (p *SampleProcessor) HandleClientTransactionShipping(ctx context.Context, t *asynq.Task) error {
	var task tasks.ClientTransactionShippingTask
	if err := json.Unmarshal(t.Payload(), &task); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	kitStatus := task.Event.Data.KitStatus
	// TODO: set kitTrackingID when shipping service is ready
	// kitTrackingID := task.Event.Data.TrackingId
	var oldKitStatus *string
	if task.Event.Old != nil {
		oldKitStatus = task.Event.Old.KitStatus
	}

	if common.Env.DryRun {
		return nil
	}

	eventType := task.Event.Type

	if (eventType == "update" && oldKitStatus != nil && *oldKitStatus != kitStatus) ||
		(eventType == "insert" && (kitStatus == "READY_FOR_SHIPMENT" || kitStatus == "READY_FOR_RETURN_SHIPMENT")) {

		// TODO: Need shipping service to get accession IDs by kitTrackingID
	}

	return nil
}

func (p *SampleProcessor) handlePostSampleOrder(task *tasks.PostSampleOrderTask, ctx context.Context) error {
	//redis lock, because job queue does tasks at least once
	lockName := fmt.Sprintf("%s_%d", common.CreateOrderLock, task.SampleId)
	lock, err := common.RedLock(lockName, p.rs)
	if err != nil {
		// if lock fails, it means some other instance is handling this task, so no need to do it again
		return nil
	}
	defer func() {
		ok, err := common.RedUnlock(lock)
		if !ok || err != nil {
			common.Error(err)
		}
	}()
	if s, err := dbutils.GetSampleById(int(task.SampleId), p.dbClient, p.ctx); err == nil && s != nil {
		// sample already exists, don't do task twice
		return nil
	}
	//TODO: add patient stuff once it's in
	//TODO: add customer stuff once it's in
	tx, err := p.dbClient.Tx(p.ctx)
	if err != nil {
		return err
	}

	if err = p.createAndFlagSampleOrder(task, tx, ctx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return dbutils.Rollback(tx, err)
	}
	//update sample required tube
	for tubeType, requiredCnt := range task.RequiredNumberOfTubes {
		err = p.SampleService.UpdateSampleRequiredTubeCount(task.SampleId, tubeType, requiredCnt, "order_service", ctx)
		if err != nil {
			common.Error(err)
		}
	}
	return nil
}

// createAndFlagSampleOrder Create order and sample and flag the order accordingly.
// This function does not commit, but will roll back on error
func (p *SampleProcessor) createAndFlagSampleOrder(task *tasks.PostSampleOrderTask, tx *ent.Tx, ctx context.Context) error {
	_, err := tx.Sample.Create().
		SetID(int(task.SampleId)).
		SetAccessionID(task.AccessionId).
		SetSampleDescription("Added by Kafka Sample Order Processing Queue").
		SetDelayedHours(0).
		SetPatientID(int(task.PatientId)).
		SetCustomerID(int(task.CustomerId)).
		SetSampleOrderMethod("ONLINE").
		Save(p.ctx)
	if err != nil {
		return dbutils.Rollback(tx, err)
	}
	//TODO: delete get_patient_info_by_id record from redis when patient is in
	//TODO: get patient info when patient is in
	//TODO: connect patient, clinic, customer
	order, err := tx.OrderInfo.Create().
		//SetPatientFirstName("TO BE SET")
		//SetPatientLastName("TO BE SET")
		SetOrderTitle("Sample Order Created By Order Kafka").
		SetOrderDescription("Order Generated Based on Order Kafka Message Content").
		SetOrderConfirmationNumber(task.OrderConfirmationNumber).
		SetBillingOrderID(task.BillingOrderId).
		AddTestIDs(task.Tests...).
		SetSampleID(int(task.SampleId)).
		Save(p.ctx)
	if err != nil {
		return dbutils.Rollback(tx, err)
	}
	//TODO: update order sales id after customer is in
	err = p.OrderService.FlagOrderTx(order.ID, []string{"order_received", "kit_wait_for_shipment", "report_not_ready"}, tx, ctx)
	if err != nil {
		// should have been rolled back in order service
		return err
	}

	if task.NonBloodKitDeliverMethod == DeliverMethodSuppliedByClient ||
		task.BloodKitDeliverMethod == DeliverMethodSuppliedByClient {
		err = p.OrderService.FlagOrderTx(order.ID, []string{"kit_patient_received_kit", "order_processing"}, tx, ctx)
		if err != nil {
			// should have been rolled back in order service
			return err
		}
	}
	return nil
}

// v1_Zhibin: LIS-Kafka General events
// handleIssueSystemEvent processes events from the issue system.
func (p *SampleProcessor) handleIssueSystemEvent(event *pb.GeneralEvent) error {
	if event.AddonColumn == nil {
		return fmt.Errorf("addon_column is nil in event: %v", event)
	}

	issueIDs := event.AddonColumn.GetIssueIds()
	if len(issueIDs) == 0 {
		return fmt.Errorf("no issue IDs found in event: %+v", event)
	}

	//TODO: issue system
	// issueInfo, err := p.fetchIssueDetails(issueIDs)
	// if err != nil {
	// 	return fmt.Errorf("failed to fetch issue details: %w", err)
	// }

	// for _, issue := range issueInfo.Issues {
	// 	var sampleID, patientID int32
	// 	foreignLinks := issue.fk.ForeignLinks
	// 	for _, foreignLink := range foreignLinks {
	// 		foreignObjectId := foreignLink.fk.ForeignObjectId
	// 		foreignObjectInfos, _ := p.getIssueForeignObject(foreignObjectId)

	// 		for _, foreignObjectInfo := range foreignObjectInfos.ForeignObjects {
	// 			if foreignObjectInfo.Name == "sample_id" {
	// 				sampleID = foreignLink.objectId
	// 			}

	// 			if foreignObjectInfo.Name == "patient_id" {
	// 				patientID = foreignLink.ObjectId
	// 			}
	// 		}

	// 	}

	// 	// Process sample-based updates
	// 	if sampleID != 0 {
	// 		//// TODO: issue system

	// 		// err = p.processIssueActionsSampleID(sampleID, event, issue)
	// 		// if err != nil {
	// 		// 	return err
	// 		// }
	// 	}

	// 	// Process patient-based updates
	// 	if patientID != 0 {
	// 		// TODO: issue system

	// 		// err = p.processIssueActionsPatientID(patientID, event, issue)
	// 		// if err != nil {
	// 		// 	return err
	// 		// }
	// 	}

	// 	// Handle case where both sampleID and patientID are missing
	// 	if sampleID == 0 && patientID == 0 {
	// 		return fmt.Errorf("no valid foreign objects found for issue %d", issueIDs)
	// 	}

	// 	// if err := p.processIssueActions(sampleID, patientID, event, issue); err != nil {
	// 	// 	common.Error(fmt.Errorf("failed to process issue: %w", err))
	// 	// }
	// }
	return nil
}

// handleOrderEvent processes events related to orders.
func (p *SampleProcessor) handleOrderEvent(event *pb.GeneralEvent, ctx context.Context) error {
	if event.AddonColumn == nil {
		return fmt.Errorf("addon_column is nil in event: %v", event)
	}

	if event.AddonColumn.GetOrderStatus() == "order_redraw" {
		sampleID := event.SampleId
		return p.handleOrderRedraw(sampleID, ctx)
	}

	return nil
}

// handleAccessioningEvent processes accessioning-related events.
func (p *SampleProcessor) handleAccessioningEvent(event *pb.GeneralEvent, ctx context.Context) error {
	if event.EventName == "change_tube_count_to_zero" {
		sampleID := event.SampleId
		tubeType := event.TubeTypes[0]
		receiveTime, _ := util.ParseEventTime(event.EventTime)
		// Delete record from Redis
		redisKey := dbutils.KeyGetSampleTests(int(sampleID))
		p.redisClient.Del(p.ctx, redisKey)
		return p.handleChangeTubeCountToZero(sampleID, tubeType, receiveTime, ctx)
	}

	return nil
}

// handleReportEvent processes events related to reports.
func (p *SampleProcessor) handleReportEvent(event *pb.GeneralEvent, ctx context.Context) error {
	switch event.EventName {
	case "old_report_opened", "new_report_viewed":
		return p.handleOldNewReportsOpenedViewed(event, ctx)
	case "personalized_report_ready":
		return p.handlePersonalizedReportReady(event, ctx)
	case "redraw_personalized_report_ready":
		return p.handleRedrawPersonalizedReportReady(event, ctx)
	case "personalized_report_updated":
		return p.handlePersonalizedReportUpdated(event, ctx)
	default:
		return nil
	}
}

// TODO: issue system
// processIssueActions applies actions for a specific issue.
// func (p *SampleProcessor) processIssueActionsSampleID(sampleID int32, event *pb.GeneralEvent, issue *pb.Issue) error {
// 	issueType := strconv.Itoa(int(issue.Fk.IssueTypeID))
// 	action := event.EventAction

// 	switch action {
// 	case "create":
// 		return p.handleIssueCreation(issueType, sampleID)
// 	case "resolve", "delete":
// 		return p.handleIssueResolution(issueType, sampleID)
// 	default:
// 		return fmt.Errorf("unsupported action: %s", action)
// 	}
// }

// func (p *SampleProcessor) processIssueActionsPatientID(patientID int32, event *pb.GeneralEvent, issue *pb.Issue) error {
// 	// issueType := strconv.Itoa(int(issue.Fk.IssueTypeID))
// 	action := event.EventAction

// 	switch action {
// 	case "create":
// 		return p.handleIssueCreation("67", patientID)
// 	case "resolve", "delete":
// 		return p.handleIssueResolution("67", patientID)
// 	default:
// 		return fmt.Errorf("unsupported action: %s", action)
// 	}
// }

// handleIssueCreation handles issue creation logic.
func (p *SampleProcessor) handleIssueCreation(issueType string, sampleID int32) error {
	switch issueType {
	case "85": // TNP issue
		return p.updateOrderStatusBySampleID(sampleID, "order_tnp_issue_status", "tnp_issue")
	case "67": // NY Waive Form issue
		return p.updateOrderStatusBySampleID(sampleID, "order_ny_waive_form_issue_status", "new_ny_waive_form_issue")
	case "92": // Incomplete questionnaire issue
		return p.updateOrderStatusBySampleID(sampleID, "order_incomplete_questionnaire_issue_status", "incomplete_questionnaire_issue")
	case "84": // Billing issue
		return p.handleBillingIssue(sampleID)
	case "86": // Shipping issue
		return p.updateOrderStatusBySampleID(sampleID, "order_shipping_issue_status", "shipping_issue")
	case "94", "100", "101":
		return p.updateOrderStatusBySampleID(sampleID, "order_missing_info_issue_status", "missing_info_issue")
	case "87":
		return p.updateOrderStatusBySampleID(sampleID, "order_receive_issue_status", "receive_issue")
	case "59", "61", "62", "63", "65":
		return p.updateOrderStatusBySampleID(sampleID, "order_lab_issue_status", "lab_issue")
	case "60", "64":
		_ = p.updateOrderStatusBySampleID(sampleID, "order_missing_info_issue_status", "missing_info_issue")
		return p.updateOrderStatusBySampleID(sampleID, "order_lab_issue_status", "lab_issue")
	default:
		return nil
	}
}

// handleIssueResolution handles issue resolution or deletion logic.
func (p *SampleProcessor) handleIssueResolution(issueType string, sampleID int32) error {
	switch issueType {
	case "85": // TNP issue
		return p.updateOrderStatusBySampleID(sampleID, "order_tnp_issue_status", "no_tnp_issue")
	case "92": // Incomplete questionnaire issue
		return p.updateOrderStatusBySampleID(sampleID, "order_incomplete_questionnaire_issue_status", "no_incomplete_questionnaire_issue")
	case "84": // Billing issue
		return p.handleBillingIssueResolution(sampleID)
	default:
		return nil
	}
}

// TODO: issue service
// fetchIssueDetails fetches details for a list of issue IDs.
// func (p *SampleProcessor) fetchIssueDetails(issueIDs []int64) (*pb.GetIssuesResponse, error) {
// 	if len(issueIDs) == 0 {
// 		return nil, fmt.Errorf("no issue IDs provided")
// 	}

// 	issuesResponse, err := p.IssueService.GetIssues(p.ctx, issueIDs)
// 	if err != nil {
// 		common.Error(fmt.Errorf("failed to fetch issues: %w", err))
// 		return nil, err
// 	}

// 	return issuesResponse, nil
// }

// // extractForeignLinks processes foreign links to find related IDs.
// func (p *SampleProcessor) getIssueForeignObject(foreignObjectIDs []int64) (*pb.GetForeignObjectsResponse, error) {
// 	if len(foreignObjectIDs) == 0 {
// 		return nil, fmt.Errorf("no foreign_object_ids provided")
// 	}

// 	foreignObjectsResponse, err := p.IssueService.GetForeignObjects(p.ctx, foreignObjectIDs)
// 	if err != nil {
// 		common.Error(fmt.Errorf("failed to fetch foreign objects: %w", err))
// 		return nil, err
// 	}

// 	return foreignObjectsResponse, nil
// }

// updateOrderStatusBySampleID updates a specific status for an order based on a sample ID.
func (p *SampleProcessor) updateOrderStatusBySampleID(sampleID int32, field, value string) error {
	return dbutils.UpdateOrderStatusBySampleID(sampleID, field, value, p.ctx, p.dbClient)
}

// handleBillingIssue handles billing-related issues.
func (p *SampleProcessor) handleBillingIssue(sampleID int32) error {
	err := p.updateOrderStatusBySampleID(sampleID, "order_billing_issue_status", "billing_issue")
	if err != nil {
		return err
	}

	return dbutils.UpdateMajorOrderStatusBySampleID(sampleID, "pending_payment_order", p.ctx, p.dbClient)
}

// handleBillingIssueResolution resolves billing-related issues.
func (p *SampleProcessor) handleBillingIssueResolution(sampleID int32) error {
	err := p.updateOrderStatusBySampleID(sampleID, "order_billing_issue_status", "no_billing_issue")
	if err != nil {
		return err
	}
	return dbutils.UpdateMajorOrderStatusBySampleID(sampleID, "awaiting_sample", p.ctx, p.dbClient)
}

// handleOrderRedraw handles order redraw logic.
func (p *SampleProcessor) handleOrderRedraw(sampleID int32, ctx context.Context) error {

	_, err := p.OrderService.FlagOrdersWithSampleId(sampleID, []string{"order_redraw_order"}, ctx)
	return err
}

// handleChangeTubeCountToZero handles events where tube count is changed to zero.
func (p *SampleProcessor) handleChangeTubeCountToZero(sampleID int32, tubeType string, receiveTime time.Time, ctx context.Context) error {

	if tubeType == "METAL_FREE_URINE" {
		return p.OrderService.DispatchRemoveSampleOrder(int(sampleID), "URINE", receiveTime, receiveTime, false, ctx)
	}

	return p.OrderService.DispatchRemoveSampleOrder(int(sampleID), tubeType, receiveTime, receiveTime, false, ctx)
}

func (p *SampleProcessor) handleOldNewReportsOpenedViewed(event *pb.GeneralEvent, ctx context.Context) error {

	sampleID := event.GetSampleId()
	// userID := event.GetInternalUserId()

	// Fetch order info
	orderInfo, err := dbutils.FindOrderWithSampleId(sampleID, p.dbClient, p.ctx)
	if err != nil {
		return fmt.Errorf("failed to query order info: %w", err)
	}

	// If report is ready, flag the order
	if orderInfo.OrderReportStatus == "report_ready" {
		_, err = p.OrderService.FlagOrder(orderInfo.ID, []string{"report_delivered", "order_completed"}, ctx)
		if err != nil {
			return fmt.Errorf("failed to flag order: %w", err)
		}
	}

	// TODO: Record audit log
	snapshot := "All Reports"

	//todo: when addon column includes report name
	// if event.AddonColumn != nil {
	// 	reportName := event.AddonColumn.ReportName // Adjust this line based on your protobuf definition
	// 	if reportName != "" {
	// 		snapshot = reportName
	// 	}
	// }

	auditLogMessage := common.AuditLogEntry{
		EventID:        event.EventId,
		ServiceName:    common.ServiceName,
		ServiceType:    "backend",
		EventName:      "report_opened",
		EntityType:     "sample_id",
		EntityID:       string(sampleID),
		User:           string(event.UserId),
		Entrypoint:     "Kafka",
		EntitySnapshot: snapshot,
	}

	common.RecordAuditLog(auditLogMessage)
	//
	return nil
}

// handlePersonalizedReportReady handles personalized report readiness events.
func (p *SampleProcessor) handlePersonalizedReportReady(event *pb.GeneralEvent, ctx context.Context) error {

	accessionID := event.GetAccessionId()

	// Fetch sample info
	sampleInfo, err := p.dbClient.Sample.
		Query().
		Where(sample.AccessionIDEQ(accessionID)).
		WithOrder().
		Only(p.ctx)
	if err != nil {
		return fmt.Errorf("failed to query sample info: %w", err)
	}

	// TODO: Add customer to get beta program
	//Fetch customer info
	// customerInfo, err := p.dbClient.Customer.
	// 	Query().
	// 	Where(Customer.IDEQ(sampleInfo.Edges.Order.CustomerID)).
	// 	WithParticipatingBetaProgram().
	// 	Only(ctx)
	// if err != nil {
	// 	return fmt.Errorf("failed to query customer info: %w", err)
	// }

	// Flag the order
	_, err = p.OrderService.FlagOrder(sampleInfo.OrderID, []string{
		"order_completed", "report_ready", "kit_received", "kit_lab_received",
	}, ctx)
	if err != nil {
		return fmt.Errorf("failed to flag order: %w", err)
	}

	//TODO: beta and kafka publisher
	// Prepare beta programs
	// betaPrograms := util.ExtractActiveBetaPrograms(customerInfo.ParticipatingBetaProgram)

	// Publish dashboard message

	err = publisher.GetPublisher().SendGeneralEvent(&pb.GeneralEvent{
		EventId:       event.EventId,
		EventProvider: "lis-report",
		EventName:     "report_finished",
		EventAction:   "notify",
		EventComment:  "This event is triggered when the report status is finished",
		EventTime:     event.EventTime,
		SampleId:      int32(sampleInfo.ID),
		PatientId:     int32(sampleInfo.PatientID),
		AccessionId:   sampleInfo.AccessionID,
		CustomerId:    int32(sampleInfo.CustomerID),
		// AddonColumn: map[string]interface{}{
		// 	"beta_program_enabled": customerInfo.BetaProgramEnabled,
		// 	"beta_programs":        betaPrograms,
		// },
	})
	if err != nil {
		return fmt.Errorf("failed to publish dashboard message: %w", err)
	}

	// Update sample report time
	parsedTime, _ := util.ParseEventTime(event.EventTime)
	_, err = sampleInfo.Update().
		SetSampleReportTime(parsedTime).
		Save(p.ctx)
	if err != nil {
		return fmt.Errorf("failed to update sample report time: %w", err)
	}

	return nil
}

// handleRedrawPersonalizedReportReady handles redraw report readiness events.
func (p *SampleProcessor) handleRedrawPersonalizedReportReady(event *pb.GeneralEvent, ctx context.Context) error {

	accessionID := event.GetAccessionId()

	// Fetch sample info
	sampleInfo, err := p.dbClient.Sample.
		Query().
		Where(sample.AccessionIDEQ(accessionID)).
		WithOrder().
		Only(p.ctx)
	if err != nil {
		return fmt.Errorf("failed to query sample info: %w", err)
	}

	// Flag the order
	_, err = p.OrderService.FlagOrder(sampleInfo.OrderID, []string{
		"order_completed", "report_ready", "kit_received", "kit_lab_received",
	}, ctx)
	if err != nil {
		return fmt.Errorf("failed to flag order: %w", err)
	}

	//TODO: kafka publisher
	// Publish dashboard message
	err = publisher.GetPublisher().SendGeneralEvent(&pb.GeneralEvent{
		EventId:       event.EventId,
		EventProvider: "lis-report",
		EventName:     "redraw_report_finished",
		EventAction:   "notify",
		EventComment:  "This event is triggered when the redraw report finished",
		EventTime:     event.EventTime,
		SampleId:      int32(sampleInfo.ID),
		PatientId:     int32(sampleInfo.PatientID),
		AccessionId:   sampleInfo.AccessionID,
		CustomerId:    int32(sampleInfo.CustomerID),
		// AddonColumn: map[string]interface{}{
		// 	"beta_program_enabled": false,
		// 	"beta_programs":        []string{},
		// },
	})

	if err != nil {
		return fmt.Errorf("failed to publish dashboard message: %w", err)
	}

	// Update sample report time
	parsedTime, _ := util.ParseEventTime(event.EventTime)
	_, err = sampleInfo.Update().
		SetSampleReportTime(parsedTime).
		Save(p.ctx)
	if err != nil {
		return fmt.Errorf("failed to update sample report time: %w", err)
	}

	return nil
}

// handlePersonalizedReportUpdated handles personalized report updated events.
func (p *SampleProcessor) handlePersonalizedReportUpdated(event *pb.GeneralEvent, ctx context.Context) error {
	accessionID := event.GetAccessionId()

	// Fetch sample info
	sampleInfo, err := p.dbClient.Sample.
		Query().
		Where(sample.AccessionIDEQ(accessionID)).
		Only(p.ctx)
	if err != nil {
		return fmt.Errorf("failed to query sample info: %w", err)
	}

	// Flag the order
	_, err = p.OrderService.FlagOrder(sampleInfo.OrderID, []string{"order_processing", "report_not_ready"}, ctx)
	if err != nil {
		return fmt.Errorf("failed to flag order: %w", err)
	}

	// Update sample report time to null
	parsedTime, _ := util.ParseEventTime(event.EventTime)
	_, err = sampleInfo.Update().
		SetSampleReportTime(parsedTime).
		Save(p.ctx)
	if err != nil {
		return fmt.Errorf("failed to clear sample report time: %w", err)
	}

	return nil
}

func (p *SampleProcessor) SendAddOnOrderToLIS(testList []int32, tubeList []string, action string, sampleID int, ctx context.Context) error {
	var sampleTypeList []string
	var testToSend []int32
	if common.Env.DryRun {
		return nil
	}
	// 🔹 Step 1: Get Potential Sample Type List
	if len(tubeList) > 0 {
		for _, tubeType := range tubeList {
			sampleTypeInfo, err := p.SampleService.GetSampleTypeViaTubeType(tubeType, p.ctx)
			if err != nil {
				return fmt.Errorf("failed to get sample type via tube: %s, err: %w", tubeType, err)
			}

			if len(sampleTypeInfo) > 0 && len(sampleTypeInfo[0].Edges.SampleTypes) > 0 {
				for _, sampleType := range sampleTypeInfo[0].Edges.SampleTypes {
					sampleTypeList = append(sampleTypeList, sampleType.SampleTypeEnum)
					if sampleType.SampleTypeEnum == "Urine" {
						sampleTypeList = append(sampleTypeList, "METAL_FREE_URINE")
					}
					if sampleType.SampleTypeEnum == "METAL_FREE_URINE" {
						sampleTypeList = append(sampleTypeList, "Urine")
					}
				}
			}
		}
	}

	for _, testID := range testList {
		if p.TestService == nil {
			return fmt.Errorf("testService is not initialized")
		}
		testInfo, err := p.TestService.GetTestField([]int{int(testID)}, []string{"test_sample_type"}, p.ctx)
		if err != nil || len(testInfo) == 0 {
			return fmt.Errorf("error getting Test: %d, err: %w", testID, err)
		}

		if len(testInfo[0].Edges.TestDetails) > 0 {
			testSampleType := testInfo[0].Edges.TestDetails[0].TestDetailsValue
			if util.Contains(sampleTypeList, testSampleType) {
				testToSend = append(testToSend, int32(testID))
			}
		}

	}

	orderMessage := &pb.OrderMessage{
		SampleId:         int32(sampleID),
		Action:           action,
		TestId:           testToSend,
		ReceiveTime:      time.Now().Format("2006-01-02 15:04:05-0700"),
		CollectionTime:   time.Now().Format("2006-01-02 15:04:05-0700"),
		IsRerun:          false,
		IsRedraw:         false,
		IsLabDirectOrder: true,
		Destination:      "DI",
		IsAddon:          true,
	}

	err := publisher.GetPublisher().SendOrderMessage(orderMessage)
	if err != nil {
		return err
	}

	return nil
}
