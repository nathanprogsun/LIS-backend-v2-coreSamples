package processor

import (
	"context"
	"coresamples/common"
	"coresamples/dbutils"
	"coresamples/ent/enttest"
	"coresamples/ent/orderinfo"
	"coresamples/ent/sample"
	"coresamples/ent/tuberequirement"
	"coresamples/model"
	"coresamples/publisher"
	"coresamples/service"
	"coresamples/tasks"
	"coresamples/util"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/hibiken/asynq"
	"github.com/stretchr/testify/assert"

	pb "coresamples/proto"

	"github.com/go-redis/redis/v8"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stvp/tempredis"
	"google.golang.org/protobuf/proto"
)

var existedOrderFlags = []string{
	"analyzing_sample",
	"awaiting_sample",
	"billing_issue",
	"canceled_order",
	"incomplete_questionnaire_issue",
	"kit_delivery_exception",
	"kit_lab_received",
	"kit_lab_shipped_kit",
	"kit_patient_received_kit",
	"kit_pending",
	"kit_ready_for_shipment",
	"kit_received",
	"kit_sample_shipped_back",
	"kit_voided_shipment",
	"kit_wait_for_shipment",
	"lab_issue",
	"missing_info_issue",
	"new_ny_waive_form_issue",
	"no_billing_issue",
	"no_incomplete_questionnaire_issue",
	"no_lab_issue",
	"no_missing_info_issue",
	"no_new_ny_waive_form_issue",
	"no_ny_waive_form_issue",
	"no_receive_issue",
	"no_shipping_issue",
	"no_tnp_issue",
	"ny_waive_form_issue",
	"order_canceled_order",
	"order_completed",
	"order_modified",
	"order_pending_payment",
	"order_processing",
	"order_received",
	"order_redraw_order",
	"order_tnp",
	"order_waiting_redraw",
	"pending_payment_order",
	"receive_issue",
	"report_amended_report",
	"report_delivered",
	"report_not_ready",
	"report_pending",
	"report_ready",
	"report_status_ready",
	"scheduled_order",
	"shipping_issue",
	"tnp_issue",
}

func setupSampleProcessorTest(t *testing.T) (*SampleProcessor, *tempredis.Server, tasks.AsynqClient) {
	dataSource := "file:ent?mode=memory&_fk=1"
	dbClient := enttest.Open(t, "sqlite3", dataSource)
	err := dbClient.Schema.Create(context.Background())
	if err != nil {
		common.Fatalf("failed opening connection to MySQL", err)
	}
	server, err := tempredis.Start(tempredis.Config{
		"port": "0",
	})
	if err != nil {
		common.Fatalf("Failed to start tempredis: %v", err)
	}

	common.InitZapLogger("debug")

	redisClient := redis.NewClient(&redis.Options{
		Network: "unix",
		Addr:    server.Socket(),
	})
	asynqClient := tasks.NewMockAsynqClient()

	// This will start a mock publisher in the memory
	publisher.InitMockPublisher()

	client := common.NewRedisClient(redisClient, redisClient)

	service.SampleSvc = service.NewSampleService(dbClient, common.NewRedisClient(redisClient, redisClient), asynqClient)
	processor := &SampleProcessor{
		Processor: InitProcessor(dbClient,
			client,
			context.Background()),
		OrderService:  service.NewOrderService(dbClient, common.NewRedisClient(redisClient, redisClient)),
		SampleService: service.SampleSvc,
		TestService:   service.NewTestService(dbClient, common.NewRedisClient(redisClient, redisClient)),
		rs:            nil,
	}

	return processor, server, asynqClient
}

func cleanUpSampleProcessorTest(p *SampleProcessor, s *tempredis.Server, asynqClient tasks.AsynqClient) {
	var err error
	if err = s.Kill(); err != nil {
		common.Error(err)
	}
	if p.dbClient != nil {
		if err = p.dbClient.Close(); err != nil {
			common.Error(err)
		}
	}

	publisher.GetPublisher().GetWriter().Close()
	asynqClient.Close()
	service.SampleSvc = nil
}

func TestHandlePostSampleOrder(t *testing.T) {
	processor, server, asynqClient := setupSampleProcessorTest(t)
	defer cleanUpSampleProcessorTest(processor, server, asynqClient)

	// This test validates the functionality of the HandlePostSampleOrder method. It ensures that a PostSampleOrder task
	// is correctly processed by enqueueing it into a mock task queue, verifying its presence in the queue, and processing
	// it to create the associated sample and order entities in the database. The test also checks that the correct flags
	// are applied to the order, the tube requirements are updated as expected, and all relevant database changes are
	// accurately reflected.

	// Create the test task
	testTask := &tasks.PostSampleOrderTask{
		SampleId:                 int32(12345),
		Tests:                    []int{1, 2, 3},
		PatientId:                int32(67890),
		OrderConfirmationNumber:  "TestOrderConfirmationNumber",
		CustomerId:               int32(54321),
		ClinicId:                 int32(101),
		BillingOrderId:           "TestBillingOrderId",
		AccessionId:              "0000000000",
		BloodKitDeliverMethod:    "suppliedByClient",
		NonBloodKitDeliverMethod: "suppliedByClient",
		RequiredNumberOfTubes:    map[string]int32{"STOOL": int32(1), "UNPRESERVED_STOOL": int32(1)},
	}
	task, err := tasks.NewPostSampleOrderTask(testTask)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	err = processor.dbClient.Customer.Create().
		SetID(54321).
		Exec(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = processor.dbClient.Patient.Create().
		SetID(67890).
		SetCustomerID(54321).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, testID := range testTask.Tests {
		_, err = processor.dbClient.Test.Create().
			SetID(testID).
			SetTestName("Test test").
			SetTestCode("Test testcode").
			SetDisplayName("Test testDisplayName").
			SetTestDescription("Test testDescription").
			SetAssayName("Test assayName").
			SetIsActive(true).
			Save(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}

	for _, flagName := range existedOrderFlags {
		_, err = processor.dbClient.OrderFlag.Create().
			SetOrderFlagName(flagName).
			Save(ctx)

		if err != nil {
			t.Fatal(err)
		}
	}

	err = processor.HandlePostSampleOrder(ctx, task)
	if err != nil {
		t.Fatalf("HandlePostSampleOrder failed: %v", err)
	}

	// Validate processing outcomes (e.g., database changes)
	// Check sample creation
	createdSample, err := dbutils.GetSampleById(int(testTask.SampleId), processor.dbClient, ctx)
	if err != nil {
		t.Fatalf("failed to fetch created sample: %v", err)
	}
	if createdSample.AccessionID != testTask.AccessionId {
		t.Fatalf("expected accession ID %s, but got %s", testTask.AccessionId, createdSample.AccessionID)
	}

	// Check order creation
	createdOrder, err := processor.dbClient.OrderInfo.Query().Where(
		orderinfo.HasSampleWith(sample.IDEQ(int(testTask.SampleId))),
	).WithOrderFlags().First(ctx)

	if err != nil {
		t.Fatalf("failed to fetch created order: %v", err)
	}
	if createdOrder.OrderConfirmationNumber != testTask.OrderConfirmationNumber {
		t.Fatalf("expected order confirmation number %s, but got %s", testTask.OrderConfirmationNumber, createdOrder.OrderConfirmationNumber)
	}

	// Check order flag
	expectedFlags := []string{
		"order_received",
		"kit_wait_for_shipment",
		"report_not_ready",
		"kit_patient_received_kit",
		"order_processing",
	}
	expectedFlagSet := make(map[string]bool, len(expectedFlags))
	for _, flagName := range expectedFlags {
		expectedFlagSet[flagName] = true
	}
	actualFlagSet := make(map[string]bool, len(createdOrder.Edges.OrderFlags))
	for _, flag := range createdOrder.Edges.OrderFlags {
		actualFlagSet[flag.OrderFlagName] = true
	}

	// Problem: flag not successful
	for flagName := range expectedFlagSet {
		if !actualFlagSet[flagName] {
			t.Fatalf("expected order to have flag %s, but it was not found", flagName)
		}
	}

	// Validate tube requirements
	for tubeType, requiredCount := range testTask.RequiredNumberOfTubes {
		tubeReq, err := dbutils.FindTubeRequirement(int(testTask.SampleId), tubeType, processor.dbClient, ctx)
		if err != nil {
			t.Fatalf("failed to fetch tube requirement for %s: %v", tubeType, err)
		}
		if tubeReq.RequiredCount != int(requiredCount) {
			t.Fatalf("expected required count %d for tube type %s, but got %d", requiredCount, tubeType, tubeReq.RequiredCount)
		}
	}
}

func TestHandleFlagOrderOnReceiving(t *testing.T) {
	processor, server, asynqClient := setupSampleProcessorTest(t)
	defer cleanUpSampleProcessorTest(processor, server, asynqClient)

	// This test ensures that a SampleTubeReceiveTask is correctly processed,
	// including its enqueueing into a mock task queue, handling by
	// the processor to apply appropriate flags to an order, and verifying
	// that all expected flags are associated with the order in the database.

	// Create the test task
	testSampleId := int32(12345)

	flagOrderTask, err := tasks.NewFlagOrderOnReceivingTask(&tasks.SampleTubeReceiveTask{SampleId: testSampleId})
	if err != nil {
		t.Fatal(err)
	}

	// Create test example and process the task
	ctx := context.Background()
	_, err = processor.dbClient.Sample.Create().
		SetID(int(testSampleId)).
		SetAccessionID("0000000000").
		SetDelayedHours(0).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = processor.dbClient.OrderInfo.Create().
		SetSampleID(int(testSampleId)).
		SetOrderConfirmationNumber("test_order").
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, flagName := range existedOrderFlags {
		_, err = processor.dbClient.OrderFlag.Create().
			SetOrderFlagName(flagName).
			Save(ctx)

		if err != nil {
			t.Fatal(err)
		}
	}

	err = processor.HandleFlagOrderOnReceiving(ctx, flagOrderTask)
	if err != nil {
		t.Fatalf("HandleFlagOrderOnReceiving failed: %v", err)
	}

	// Check flagged order
	order, err := processor.dbClient.OrderInfo.Query().Where(
		orderinfo.HasSampleWith(sample.IDEQ(int(testSampleId))),
	).WithOrderFlags().First(ctx)

	if err != nil {
		t.Fatalf("failed to fetch created order: %v", err)
	}

	expectedFlags := []string{
		"kit_received",
		"order_processing",
		"report_pending",
		"kit_lab_received",
	}

	expectedFlagSet := make(map[string]bool, len(expectedFlags))
	for _, flagName := range expectedFlags {
		expectedFlagSet[flagName] = true
	}

	actualFlagSet := make(map[string]bool, len(order.Edges.OrderFlags))
	for _, flag := range order.Edges.OrderFlags {
		actualFlagSet[flag.OrderFlagName] = true
	}
	// Problem: flag not successful
	for flagName := range expectedFlagSet {
		if !actualFlagSet[flagName] {
			t.Fatalf("expected order to have flag %s, but it was not found. Order flag length %d.", flagName, len(order.Edges.OrderFlags))
		}
	}
}

func TestHandleSendOrderOnReceiving(t *testing.T) {
	processor, server, asynqClient := setupSampleProcessorTest(t)
	defer cleanUpSampleProcessorTest(processor, server, asynqClient)

	receivedTime := time.Now()
	isRedraw := false
	testSampleId := int32(12345)
	tubeDetails := []*model.SampleTubeDetails{
		{
			TubeType:       "METAL_FREE_URINE",
			CollectionTime: receivedTime.Add(-1 * time.Hour),
			ReceiveCount:   2,
		},
		{
			TubeType:       "COVID19_STOOL",
			CollectionTime: receivedTime.Add(-2 * time.Hour),
			ReceiveCount:   1,
		},
	}
	sendOrderTask, err := tasks.NewSendOrderOnReceivingTask(&tasks.SampleTubeReceiveTask{
		SampleId:     testSampleId,
		TubeDetails:  tubeDetails,
		ReceivedTime: receivedTime,
		IsRedraw:     isRedraw,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	_, err = processor.dbClient.Sample.Create().
		SetID(int(testSampleId)).
		SetAccessionID("0000000000").
		SetDelayedHours(0).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = processor.dbClient.OrderInfo.Create().
		SetSampleID(int(testSampleId)).
		SetOrderConfirmationNumber("test_order").
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = processor.HandleSendOrderOnReceiving(ctx, sendOrderTask)
	if err != nil {
		t.Fatalf("HandleSendOrderOnReceiving failed: %v", err)
	}

	// Validate the process

	// Problem: redis cached but no record found in LabOrderSendRecord
	for _, detail := range tubeDetails {
		tubeType := detail.TubeType

		// Validate whether CreateLabOrderSendRecord for each TubeType
		_, err := dbutils.GetLabOrderSendRecord(int(testSampleId), tubeType, false, true, processor.dbClient, ctx)
		if err != nil {
			t.Fatalf("expected lab order send record for tube type %s, but not found: %v", tubeType, err)
		}
	}

	// Validate CreateLabOrderSendRecord for METAL_FREE_URINE and COVID19_STOOL special handling
	specialTubeTypes := []string{"URINE", "COVID_STOOL_Triggered_By_COVID19_STOOL"}
	for _, tubeType := range specialTubeTypes {
		_, err := dbutils.GetLabOrderSendRecord(int(testSampleId), tubeType, false, true, processor.dbClient, ctx)
		if err != nil {
			t.Fatalf("expected special lab order send record for tube type %s, but not found: %v", tubeType, err)
		}
	}

	// Validate DispatchOrderWithReceivedTime
	mockWriter := publisher.GetPublisher().GetWriter().(*publisher.MockKafkaWriter)
	if len(mockWriter.MessageQueue) < len(tubeDetails) {
		t.Fatalf("expected at least %d Kafka messages, but got %d", len(tubeDetails), len(mockWriter.MessageQueue))
	}

	for _, msg := range mockWriter.MessageQueue {
		orderMessage := &pb.OrderMessage{}
		if err := json.Unmarshal(msg.Value, orderMessage); err != nil {
			t.Fatalf("failed to unmarshal Kafka message payload: %v", err)
		}

		// Validate the OrderMessage fields
		if orderMessage.SampleId != testSampleId {
			t.Fatalf("expected SampleId %d, but got %d", testSampleId, orderMessage.SampleId)
		}
		if orderMessage.Action != "+" {
			t.Fatalf("expected Action '+', but got '%s'", orderMessage.Action)
		}
		if orderMessage.IsRedraw != isRedraw {
			t.Fatalf("expected IsRedraw %v, but got %v", isRedraw, orderMessage.IsRedraw)
		}
		if orderMessage.Destination != "DI" {
			t.Fatalf("expected Destination 'DI', but got '%s'", orderMessage.Destination)
		}

	}
}

func TestHandleSendSampleReceiveGeneralEvent(t *testing.T) {
	processor, server, asynqClient := setupSampleProcessorTest(t)
	defer cleanUpSampleProcessorTest(processor, server, asynqClient)

	testSampleId := int32(12345)
	receivedTime := time.Now()
	isRedraw := false
	tubeDetails := []*pb.TubeDetail{
		{
			TubeType:       "METAL_FREE_URINE",
			CollectionTime: receivedTime.Add(-1 * time.Hour).Format(time.RFC3339),
			ReceiveCount:   2,
		},
		{
			TubeType:       "COVID19_STOOL",
			CollectionTime: receivedTime.Add(-2 * time.Hour).Format(time.RFC3339),
			ReceiveCount:   1,
		},
	}
	event := &pb.GeneralEvent{
		SampleId: testSampleId,
		AddonColumn: &pb.EventAddonColumn{
			TubeDetails:  tubeDetails,
			ReceivedBy:   "test.test",
			ReceivedTime: receivedTime.Format(time.RFC3339),
			IsRedraw:     isRedraw,
		},
		EventId:       "TestEventId",
		EventProvider: "lis-shipping",
		EventName:     "receive_sample_tubes",
	}
	sendSampletask, err := tasks.NewSendSampleReceiveGeneralEventTask(event)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	err = processor.HandleSendSampleReceiveGeneralEvent(ctx, sendSampletask)
	if err != nil {
		t.Fatalf("HandleSendSampleReceiveGeneralEvent failed: %v", err)
	}

	// Validate publisher
	mockWriter := publisher.GetPublisher().GetWriter().(*publisher.MockKafkaWriter)
	if len(mockWriter.MessageQueue) < 1 {
		t.Fatalf("expected at least %d Kafka messages, but got %d", 1, len(mockWriter.MessageQueue))
	}

	// Validate the published message
	publishedMessage := mockWriter.MessageQueue[len(mockWriter.MessageQueue)-1]
	publishedEvent := &pb.GeneralEvent{}
	if err := json.Unmarshal(publishedMessage.Value, publishedEvent); err != nil {
		t.Fatalf("failed to unmarshal published Kafka message: %v", err)
	}

	if !proto.Equal(publishedEvent, event) {
		t.Fatalf("published event does not match the original event, got: %+v, expected: %+v", publishedEvent, event)
	}
}

func TestHandleSampleOrderGeneralEvent1(t *testing.T) {
	processor, server, asynqClient := setupSampleProcessorTest(t)
	defer cleanUpSampleProcessorTest(processor, server, asynqClient)

	testKafkaMessage := []byte(`{
		"event_id": "4b32e2f7-9d5c-4378-8d1a-e6ed09e0532e",
		"event_provider": "lis-order",
		"event_name": "order_status_updates",
		"event_action": "create",
		"event_comment": "This event is triggered when the order status is updated",
		"event_time": "2025-01-07T16:12:06.148Z",
		"addon_column": {
			"orderStatus": "order_redraw",
			"sampleId": "2104752",
			"patientId": "607752",
			"customerId": "6110",
			"orderConfirmationNumber": "VA-0006110-607752-002",
			"sampleOrderTime": "2024-06-10 18:29:41",
			"clinicId": "79972",
			"is_patient_has_address": false,
			"email_address": "test@test.com"
		},
		"sample_types": [],
		"sample_id": 2104752,
		"accession_id": "2406106419",
		"customer_id": 6110,
		"patient_id": 607752,
		"clinic_id": 79972,
		"order_confirmation_number": "VA-0006110-607752-002"
	}`)

	event := &pb.GeneralEvent{}
	err := json.Unmarshal(testKafkaMessage, event)
	if err != nil {
		t.Fatalf("failed to unmarshal message: %v", err)
	}

	task, _ := tasks.NewSampleOrderGeneralEvent(&tasks.GeneralEventTask{
		Event: event,
	})

	ctx := context.Background()
	_, err = processor.dbClient.Sample.Create().
		SetID(2104752).
		SetAccessionID("2406106419").
		SetDelayedHours(0).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = processor.dbClient.OrderInfo.Create().
		SetSampleID(2104752).
		SetOrderConfirmationNumber("VA-0006110-607752-002").
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, flagName := range existedOrderFlags {
		_, err = processor.dbClient.OrderFlag.Create().
			SetOrderFlagName(flagName).
			Save(ctx)

		if err != nil {
			t.Fatal(err)
		}
	}

	err = processor.HandleSampleOrderGeneralEvent(ctx, task)

	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	order, err := processor.dbClient.OrderInfo.Query().Where(
		orderinfo.HasSampleWith(sample.IDEQ(int(event.SampleId))),
	).WithOrderFlags().First(ctx)
	if err != nil {
		t.Fatalf("failed to query order info: %v", err)
	}

	//then check whether order has the flag "order_redraw_order"
	// Assert: Check if the order has the "order_redraw_order" flag
	hasRedrawFlag := false
	for _, flag := range order.Edges.OrderFlags {
		if flag.OrderFlagName == "order_redraw_order" {
			hasRedrawFlag = true
			break
		}
	}

	if !hasRedrawFlag {
		t.Fatal("expected order to have the 'order_redraw_order' flag, but it was not found")
	}
}

func TestHandleSampleOrderGeneralEvent2(t *testing.T) {
	processor, server, asynqClient := setupSampleProcessorTest(t)
	defer cleanUpSampleProcessorTest(processor, server, asynqClient)

	testKafkaMessage := []byte(`{
		"event_id": "5522057e-9462-4ec4-9c34-9a0f217452ee",
		"event_provider": "lis-accessioning",
		"event_name": "change_tube_count_to_zero",
		"event_action": "update",
		"event_comment": "This event is triggered when the received tube count in a tube receiving record is updated to zero",
		"event_time": "2025-01-16T23:25:47.378Z",
		"tube_types": [
			"SST"
		],
		"internal_user_id": 908,
		"sample_id": 2220471
	}`)

	event := &pb.GeneralEvent{}
	err := json.Unmarshal(testKafkaMessage, event)
	if err != nil {
		t.Fatalf("failed to unmarshal message: %v", err)
	}

	sampleID := event.SampleId
	tubeType := event.TubeTypes[0]
	task, _ := tasks.NewSampleOrderGeneralEvent(&tasks.GeneralEventTask{
		Event: event,
	})

	ctx := context.Background()
	err = processor.dbClient.Customer.Create().
		SetID(54321).
		Exec(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = processor.dbClient.Sample.Create().
		SetID(int(sampleID)).
		SetAccessionID("testAccessionId").
		SetDelayedHours(0).
		SetCustomerID(54321).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = processor.dbClient.OrderInfo.Create().
		SetSampleID(int(sampleID)).
		SetOrderConfirmationNumber("test_order").
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = processor.HandleSampleOrderGeneralEvent(ctx, task)

	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	labOrderSendRecord, err := dbutils.GetLabOrderSendRecord(int(sampleID), tubeType, true, true, processor.dbClient, ctx)

	if err != nil {
		t.Fatalf("expected special lab order send record for tube type %s, but not found: %v", tubeType, err)
	}

	if labOrderSendRecord.Action != "-" {
		t.Fatalf("expected lab order send record's action '-', but got: %v", labOrderSendRecord.Action)
	}

	mockWriter := publisher.GetPublisher().GetWriter().(*publisher.MockKafkaWriter)
	if len(mockWriter.MessageQueue) < 1 {
		t.Fatalf("expected at least %d Kafka messages, but got %d", 1, len(mockWriter.MessageQueue))
	}

	// Validate the published message
	publishedMessage := mockWriter.MessageQueue[len(mockWriter.MessageQueue)-1]
	msg := &pb.OrderMessage{}
	if err := json.Unmarshal(publishedMessage.Value, msg); err != nil {
		t.Fatalf("failed to unmarshal published Kafka message: %v", err)
	}

	// Validate the OrderMessage fields
	if msg.SampleId != sampleID {
		t.Fatalf("expected SampleId %d, but got %d", sampleID, msg.SampleId)
	}
	if msg.Action != "-" {
		t.Fatalf("expected Action '-', but got '%s'", msg.Action)
	}
	if msg.IsRedraw != false {
		t.Fatalf("expected IsRedraw %v, but got %v", false, msg.IsRedraw)
	}
	if msg.Destination != "DI" {
		t.Fatalf("expected Destination 'DI', but got '%s'", msg.Destination)
	}
}

func TestHandleSampleOrderGeneralEvent3(t *testing.T) {
	processor, server, asynqClient := setupSampleProcessorTest(t)
	defer cleanUpSampleProcessorTest(processor, server, asynqClient)

	testKafkaMessage := []byte(`{
		"event_time": "2025-01-21T22:50:22.797Z",
		"event_provider": "lis-report",
		"event_name": "new_report_viewed",
		"event_action": "update",
		"event_comment": "",
		"addon_column": {},
		"tube_types": [],
		"test_ids": [],
		"sample_types": [],
		"products": [
			"AllReport"
		],
		"internal_user_id": null,
		"sample_id": 2206575,
		"accession_id": "2412126599",
		"customer_id": 9843,
		"patient_id": 3004134,
		"clinic_id": null,
		"order_id": null,
		"order_confirmation_number": null,
		"event_id": "d27bbfff-8ca7-427e-9965-5d434a4c8a27",
		"user_id": 26489
	}`)

	event := &pb.GeneralEvent{}
	err := json.Unmarshal(testKafkaMessage, event)
	if err != nil {
		t.Fatalf("failed to unmarshal message: %v", err)
	}

	task, _ := tasks.NewSampleOrderGeneralEvent(&tasks.GeneralEventTask{
		Event: event,
	})

	ctx := context.Background()
	_, err = processor.dbClient.Customer.Create().
		SetID(9843).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = processor.dbClient.Sample.Create().
		SetID(2206575).
		SetAccessionID("2412126599").
		SetDelayedHours(0).
		SetCustomerID(9843).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = processor.dbClient.OrderInfo.Create().
		SetSampleID(2206575).
		SetOrderConfirmationNumber("test_order").
		SetOrderReportStatus("report_ready").
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, flagName := range existedOrderFlags {
		_, err = processor.dbClient.OrderFlag.Create().
			SetOrderFlagName(flagName).
			Save(ctx)

		if err != nil {
			t.Fatal(err)
		}
	}

	err = processor.HandleSampleOrderGeneralEvent(ctx, task)

	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	order, err := processor.dbClient.OrderInfo.Query().Where(
		orderinfo.HasSampleWith(sample.IDEQ(int(event.SampleId))),
	).WithOrderFlags().First(ctx)
	if err != nil {
		t.Fatalf("failed to query order info: %v", err)
	}

	expectedFlags := []string{
		"report_delivered",
		"order_completed",
	}

	expectedFlagSet := make(map[string]bool, len(expectedFlags))
	for _, flagName := range expectedFlags {
		expectedFlagSet[flagName] = true
	}

	actualFlagSet := make(map[string]bool, len(order.Edges.OrderFlags))
	for _, flag := range order.Edges.OrderFlags {
		actualFlagSet[flag.OrderFlagName] = true
	}
	// Problem: flag not successful
	for flagName := range expectedFlagSet {
		if !actualFlagSet[flagName] {
			t.Fatalf("expected order to have flag %s, but it was not found. Order flag length %d.", flagName, len(order.Edges.OrderFlags))
		}
	}
}

func TestHandleSampleOrderGeneralEvent4(t *testing.T) {
	processor, server, asynqClient := setupSampleProcessorTest(t)
	defer cleanUpSampleProcessorTest(processor, server, asynqClient)

	testKafkaMessage := []byte(`{
		"event_time": "2025-01-22T01:43:30.206Z",
		"event_provider": "lis-report",
		"event_name": "personalized_report_ready",
		"event_action": "create",
		"event_comment": "This event is triggered when all the personalized reports under this barcode are already ready.",
		"addon_column": {},
		"tube_types": [],
		"test_ids": [],
		"sample_types": [],
		"products": [
			"AllReport"
		],
		"internal_user_id": null,
		"sample_id": null,
		"accession_id": "2412216035",
		"customer_id": null,
		"patient_id": null,
		"clinic_id": null,
		"order_id": null,
		"order_confirmation_number": null,
		"event_id": "85c1ddfc-a539-49af-a07d-6aba83a36a82"
	}`)

	event := &pb.GeneralEvent{}
	err := json.Unmarshal(testKafkaMessage, event)
	if err != nil {
		t.Fatalf("failed to unmarshal message: %v", err)
	}

	task, _ := tasks.NewSampleOrderGeneralEvent(&tasks.GeneralEventTask{
		Event: event,
	})

	ctx := context.Background()
	testSampleReportTime, _ := util.ParseEventTime("2025-01-01T01:43:30.206Z")
	createdSample, err := processor.dbClient.Sample.Create().
		SetAccessionID("2412216035").
		SetDelayedHours(0).
		SetSampleReportTime(testSampleReportTime).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = processor.dbClient.OrderInfo.Create().
		SetSampleID(createdSample.ID).
		SetOrderConfirmationNumber("test_order").
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, flagName := range existedOrderFlags {
		_, err = processor.dbClient.OrderFlag.Create().
			SetOrderFlagName(flagName).
			Save(ctx)

		if err != nil {
			t.Fatal(err)
		}
	}

	err = processor.HandleSampleOrderGeneralEvent(ctx, task)

	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	order, err := processor.dbClient.OrderInfo.Query().Where(
		orderinfo.HasSampleWith(sample.IDEQ(createdSample.ID)),
	).WithOrderFlags().First(ctx)
	if err != nil {
		t.Fatalf("failed to query order info: %v", err)
	}

	expectedFlags := []string{
		"order_completed",
		"report_ready",
		"kit_received",
		"kit_lab_received",
	}

	expectedFlagSet := make(map[string]bool, len(expectedFlags))
	for _, flagName := range expectedFlags {
		expectedFlagSet[flagName] = true
	}

	actualFlagSet := make(map[string]bool, len(order.Edges.OrderFlags))
	for _, flag := range order.Edges.OrderFlags {
		actualFlagSet[flag.OrderFlagName] = true
	}
	// Problem: flag not successful
	for flagName := range expectedFlagSet {
		if !actualFlagSet[flagName] {
			t.Fatalf("expected order to have flag %s, but it was not found. Order flag length %d.", flagName, len(order.Edges.OrderFlags))
		}
	}

	mockWriter := publisher.GetPublisher().GetWriter().(*publisher.MockKafkaWriter)
	if len(mockWriter.MessageQueue) < 1 {
		t.Fatalf("expected at least %d Kafka messages, but got %d", 1, len(mockWriter.MessageQueue))
	}

	// Validate the published message
	publishedMessage := mockWriter.MessageQueue[len(mockWriter.MessageQueue)-1]
	publishedEvent := &pb.GeneralEvent{}
	if err := json.Unmarshal(publishedMessage.Value, publishedEvent); err != nil {
		t.Fatalf("failed to unmarshal published Kafka message: %v", err)
	}
	// Validate the OrderMessage fields
	if publishedEvent.SampleId != int32(createdSample.ID) {
		t.Fatalf("expected SampleId %d in Kafka message, but got %d", int32(createdSample.ID), publishedEvent.SampleId)
	}

	if publishedEvent.EventName != "report_finished" {
		t.Fatalf("expected EventName %s, but got %s", "report_finished", publishedEvent.EventName)
	}

	updatedSample, err := processor.dbClient.Sample.Query().
		Where(sample.IDEQ(createdSample.ID)).First(ctx)
	if err != nil {
		t.Fatalf("failed to query sample: %v", err)
	}

	expectedSampleReportTime, _ := util.ParseEventTime(event.EventTime)
	if updatedSample.SampleReportTime != expectedSampleReportTime {
		t.Fatalf("sample report time not updated to %s", expectedSampleReportTime)
	}
}

func TestHandleSampleOrderGeneralEvent5(t *testing.T) {
	processor, server, asynqClient := setupSampleProcessorTest(t)
	defer cleanUpSampleProcessorTest(processor, server, asynqClient)

	testKafkaMessage := []byte(`{
		"event_time": "2025-01-16T02:30:02.377Z",
		"event_provider": "lis-report",
		"event_name": "redraw_personalized_report_ready",
		"event_action": "create",
		"event_comment": "This event is triggered when all the personalized reports under this redraw sample barcode are ready.",
		"addon_column": {},
		"tube_types": [],
		"test_ids": [],
		"sample_types": [],
		"products": [
			"AllReport"
		],
		"internal_user_id": null,
		"sample_id": null,
		"accession_id": "2412166232",
		"customer_id": null,
		"patient_id": null,
		"clinic_id": null,
		"order_id": null,
		"order_confirmation_number": null,
		"event_id": "9f2574f8-0b59-4fbc-8255-be7e8ed51d4a"
	}`)

	event := &pb.GeneralEvent{}
	err := json.Unmarshal(testKafkaMessage, event)
	if err != nil {
		t.Fatalf("failed to unmarshal message: %v", err)
	}

	task, _ := tasks.NewSampleOrderGeneralEvent(&tasks.GeneralEventTask{
		Event: event,
	})

	ctx := context.Background()
	testSampleReportTime, _ := util.ParseEventTime("2025-01-01T01:43:30.206Z")
	createdSample, err := processor.dbClient.Sample.Create().
		SetAccessionID("2412166232").
		SetDelayedHours(0).
		SetSampleReportTime(testSampleReportTime).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = processor.dbClient.OrderInfo.Create().
		SetSampleID(createdSample.ID).
		SetOrderConfirmationNumber("test_order").
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, flagName := range existedOrderFlags {
		_, err = processor.dbClient.OrderFlag.Create().
			SetOrderFlagName(flagName).
			Save(ctx)

		if err != nil {
			t.Fatal(err)
		}
	}

	err = processor.HandleSampleOrderGeneralEvent(ctx, task)

	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	order, err := processor.dbClient.OrderInfo.Query().Where(
		orderinfo.HasSampleWith(sample.IDEQ(createdSample.ID)),
	).WithOrderFlags().First(ctx)
	if err != nil {
		t.Fatalf("failed to query order info: %v", err)
	}

	expectedFlags := []string{
		"order_completed",
		"report_ready",
		"kit_received",
		"kit_lab_received",
	}

	expectedFlagSet := make(map[string]bool, len(expectedFlags))
	for _, flagName := range expectedFlags {
		expectedFlagSet[flagName] = true
	}

	actualFlagSet := make(map[string]bool, len(order.Edges.OrderFlags))
	for _, flag := range order.Edges.OrderFlags {
		actualFlagSet[flag.OrderFlagName] = true
	}
	// Problem: flag not successful
	for flagName := range expectedFlagSet {
		if !actualFlagSet[flagName] {
			t.Fatalf("expected order to have flag %s, but it was not found. Order flag length %d.", flagName, len(order.Edges.OrderFlags))
		}
	}

	mockWriter := publisher.GetPublisher().GetWriter().(*publisher.MockKafkaWriter)
	if len(mockWriter.MessageQueue) < 1 {
		t.Fatalf("expected at least %d Kafka messages, but got %d", 1, len(mockWriter.MessageQueue))
	}

	// Validate the published message
	publishedMessage := mockWriter.MessageQueue[len(mockWriter.MessageQueue)-1]
	publishedEvent := &pb.GeneralEvent{}
	if err := json.Unmarshal(publishedMessage.Value, publishedEvent); err != nil {
		t.Fatalf("failed to unmarshal published Kafka message: %v", err)
	}
	// Validate the OrderMessage fields
	if publishedEvent.SampleId != int32(createdSample.ID) {
		t.Fatalf("expected SampleId %d in Kafka message, but got %d", int32(createdSample.ID), publishedEvent.SampleId)
	}

	if publishedEvent.EventName != "redraw_report_finished" {
		t.Fatalf("expected EventName %s, but got %s", "redraw_report_finished", publishedEvent.EventName)
	}

	updatedSample, err := processor.dbClient.Sample.Query().
		Where(sample.IDEQ(createdSample.ID)).First(ctx)
	if err != nil {
		t.Fatalf("failed to query sample: %v", err)
	}

	expectedSampleReportTime, _ := util.ParseEventTime(event.EventTime)
	if updatedSample.SampleReportTime != expectedSampleReportTime {
		t.Fatalf("sample report time not updated to %s", expectedSampleReportTime)
	}
}

func TestHandleSampleOrderGeneralEvent6(t *testing.T) {
	processor, server, asynqClient := setupSampleProcessorTest(t)
	defer cleanUpSampleProcessorTest(processor, server, asynqClient)

	testKafkaMessage := []byte(`{
		"event_time": "2025-01-15T08:38:27.686Z",
		"event_provider": "lis-report",
		"event_name": "personalized_report_updated",
		"event_action": "update",
		"event_comment": "This event is triggered when any one of the personalized reports under this barcode just got updated.",
		"addon_column": {},
		"tube_types": [],
		"test_ids": [],
		"sample_types": [],
		"products": [
			"Micronutrients"
		],
		"internal_user_id": null,
		"sample_id": null,
		"accession_id": "2501026006",
		"customer_id": null,
		"patient_id": null,
		"clinic_id": null,
		"order_id": null,
		"order_confirmation_number": null,
		"event_id": "01796604-d503-4ace-85de-b1bfcb287da0"
	}`)

	event := &pb.GeneralEvent{}
	err := json.Unmarshal(testKafkaMessage, event)
	if err != nil {
		t.Fatalf("failed to unmarshal message: %v", err)
	}

	task, _ := tasks.NewSampleOrderGeneralEvent(&tasks.GeneralEventTask{
		Event: event,
	})

	ctx := context.Background()
	testSampleReportTime, _ := util.ParseEventTime("2025-01-01T01:43:30.206Z")
	createdSample, err := processor.dbClient.Sample.Create().
		SetAccessionID("2501026006").
		SetDelayedHours(0).
		SetSampleReportTime(testSampleReportTime).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = processor.dbClient.OrderInfo.Create().
		SetSampleID(createdSample.ID).
		SetOrderConfirmationNumber("test_order").
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, flagName := range existedOrderFlags {
		_, err = processor.dbClient.OrderFlag.Create().
			SetOrderFlagName(flagName).
			Save(ctx)

		if err != nil {
			t.Fatal(err)
		}
	}

	err = processor.HandleSampleOrderGeneralEvent(ctx, task)

	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	order, err := processor.dbClient.OrderInfo.Query().Where(
		orderinfo.HasSampleWith(sample.IDEQ(createdSample.ID)),
	).WithOrderFlags().First(ctx)
	if err != nil {
		t.Fatalf("failed to query order info: %v", err)
	}

	expectedFlags := []string{
		"order_processing",
		"report_not_ready",
	}

	expectedFlagSet := make(map[string]bool, len(expectedFlags))
	for _, flagName := range expectedFlags {
		expectedFlagSet[flagName] = true
	}

	actualFlagSet := make(map[string]bool, len(order.Edges.OrderFlags))
	for _, flag := range order.Edges.OrderFlags {
		actualFlagSet[flag.OrderFlagName] = true
	}
	// Problem: flag not successful
	for flagName := range expectedFlagSet {
		if !actualFlagSet[flagName] {
			t.Fatalf("expected order to have flag %s, but it was not found. Order flag length %d.", flagName, len(order.Edges.OrderFlags))
		}
	}

	updatedSample, err := processor.dbClient.Sample.Query().
		Where(sample.IDEQ(createdSample.ID)).First(ctx)
	if err != nil {
		t.Fatalf("failed to query sample: %v", err)
	}

	expectedSampleReportTime, _ := util.ParseEventTime(event.EventTime)
	if updatedSample.SampleReportTime != expectedSampleReportTime {
		t.Fatalf("sample report time not updated to %s", expectedSampleReportTime)
	}
}

func TestHandleCancelOrderEvent(t *testing.T) {
	processor, server, asynqClient := setupSampleProcessorTest(t)
	defer cleanUpSampleProcessorTest(processor, server, asynqClient)

	testKafkaMessage := []byte(`{
		"orderId": 30139923542650730,
		"sampleId": 675081,
		"patientId": 442134,
		"customerId": 8326
	}`)

	event := &pb.CancelOrderEvent{}
	err := json.Unmarshal(testKafkaMessage, event)
	if err != nil {
		t.Fatalf("failed to unmarshal message: %v", err)
	}
	task, err := tasks.NewCancelSampleOrderTask(&tasks.CancelOrderTask{
		Event: event,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	createdSample, err := processor.dbClient.Sample.Create().
		SetID(int(event.SampleId)).
		SetAccessionID("test_accession_id").
		SetDelayedHours(0).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = processor.dbClient.OrderInfo.Create().
		SetSampleID(createdSample.ID).
		SetOrderConfirmationNumber("test_order").
		SetOrderStatus("test_order_status").
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	for _, flagName := range existedOrderFlags {
		_, err = processor.dbClient.OrderFlag.Create().
			SetOrderFlagName(flagName).
			Save(ctx)

		if err != nil {
			t.Fatal(err)
		}
	}

	err = processor.HandleCancelOrderEvent(ctx, task)

	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	order, err := processor.dbClient.OrderInfo.Query().Where(
		orderinfo.HasSampleWith(sample.IDEQ(int(event.SampleId))),
	).WithOrderFlags().First(ctx)
	if err != nil {
		t.Fatalf("failed to query order info: %v", err)
	}

	hasCancelFlag := false
	for _, flag := range order.Edges.OrderFlags {
		if flag.OrderFlagName == "order_canceled_order" {
			hasCancelFlag = true
			break
		}
	}

	if !hasCancelFlag {
		t.Fatal("expected order to have the 'order_canceled_order' flag, but it was not found")
	}

	if order.OrderStatus != "order_canceled_order" ||
		order.OrderMajorStatus != "canceled_order" ||
		!order.OrderCanceled {
		t.Fatal("failed to update order status")
	}
}

func TestHandleClientTransactionShipping(t *testing.T) {
	processor, server, asynqClient := setupSampleProcessorTest(t)
	defer cleanUpSampleProcessorTest(processor, server, asynqClient)

	testKafkaMessage := []byte(`{
		"database": "inventory",
		"table": "client_transaction_shipping",
		"type": "update",
		"ts": 1739321116,
		"xid": 13925989445,
		"commit": true,
		"data": {
			"index": 1405700,
			"po_number": 482874,
			"tracking_id": "772022434470",
			"track_id_type": "OUT",
			"shipping_method": "FEDEX_2_DAY",
			"current_status": "Please check back later for shipment status or subscribe for e-mail notifications",
			"last_update_time": "2025-02-11 18:26:00",
			"kit_status": "READY_FOR_SHIPMENT",
			"estimated_delivery_date": "20250214",
			"display_est_delivery_time": "2000",
			"fedex_delivery_date": null,
			"box_receive_time": null,
			"box_receive_by": null,
			"track_id_delete_time": null,
			"track_id_delete_by": null,
			"track_id_delete_reason": null,
			"client_id": "26953",
			"customer_name": "QUINN SENN",
			"customer_practice_name": "",
			"customer_phone_number": "+1-9209739662",
			"customer_street": "1511 PHILIPPEN STREET",
			"customer_city": "MANITOWOC",
			"customer_state": "WI",
			"customer_zipcode": "54220",
			"customer_country": "US",
			"comment": null
		},
		"old": {
			"current_status": "LABEL CREATED",
			"last_update_time": null,
			"estimated_delivery_date": null,
			"display_est_delivery_time": null
		}
	}`)

	event := &pb.ClientTransactionShippingEvent{}
	err := json.Unmarshal(testKafkaMessage, event)
	if err != nil {
		t.Fatalf("failed to unmarshal message: %v", err)
	}
	task, err := tasks.NewClientTransactionShippingTask(&tasks.ClientTransactionShippingTask{
		Event: event,
	})
	if err != nil {
		t.Fatal(err)
	}

	// TODO: finish the test logic after HandleClientTransactionShipping logic is in

	ctx := context.Background()
	err = processor.HandleClientTransactionShipping(ctx, task)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}
}

func TestHandleClientTransactionShipping2(t *testing.T) {
	processor, server, asynqClient := setupSampleProcessorTest(t)
	defer cleanUpSampleProcessorTest(processor, server, asynqClient)

	testKafkaMessage := []byte(`{
	"database": "inventory",
	"table": "client_transaction_shipping",
	"type": "insert",
	"ts": 1739838713,
	"xid": 14086114576,
	"commit": true,
	"data": {
		"index": 1411006,
		"po_number": 484956,
		"tracking_id": "772137172011",
		"track_id_type": "OUT",
		"shipping_method": "FEDEX_2_DAY",
		"current_status": "LABEL CREATED",
		"last_update_time": null,
		"kit_status": "READY_FOR_SHIPMENT",
		"estimated_delivery_date": null,
		"display_est_delivery_time": null,
		"fedex_delivery_date": null,
		"box_receive_time": null,
		"box_receive_by": null,
		"track_id_delete_time": null,
		"track_id_delete_by": null,
		"track_id_delete_reason": null,
		"client_id": "34605",
		"customer_name": "Stephen McPherson",
		"customer_practice_name": "",
		"customer_phone_number": "8663640963",
		"customer_street": "3407 EAST HASTINGS AVENUE ",
		"customer_city": "MEAD",
		"customer_state": "WA",
		"customer_zipcode": "99021",
		"customer_country": "US",
		"comment": null
	}
}`)

	event := &pb.ClientTransactionShippingEvent{}
	err := json.Unmarshal(testKafkaMessage, event)
	if err != nil {
		t.Fatalf("failed to unmarshal message: %v", err)
	}
	task, err := tasks.NewClientTransactionShippingTask(&tasks.ClientTransactionShippingTask{
		Event: event,
	})
	if err != nil {
		t.Fatal(err)
	}

	// TODO: finish the test logic after HandleClientTransactionShipping logic is in

	ctx := context.Background()
	err = processor.HandleClientTransactionShipping(ctx, task)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}
}

func TestSampleProcessor_HandleRedrawOrderEvent(t *testing.T) {
	// Setup
	processor, server, asynqClient := setupSampleProcessorTest(t)
	defer cleanUpSampleProcessorTest(processor, server, asynqClient)

	ctx := context.Background()

	// Test variables
	sampleID := int32(123)
	testOrderID := 456

	// Create the test sample in the database
	_, err := processor.dbClient.Sample.Create().
		SetID(int(sampleID)).
		SetAccessionID("test_accession").
		SetDelayedHours(0).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Create the test order in the database
	_, err = processor.dbClient.OrderInfo.Create().
		SetID(testOrderID).
		SetSampleID(int(sampleID)).
		SetOrderConfirmationNumber("test_order").
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Create order flags
	for _, flagName := range existedOrderFlags {
		_, err = processor.dbClient.OrderFlag.Create().
			SetOrderFlagName(flagName).
			Save(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create the task payload directly as JSON
	taskPayload := []byte(fmt.Sprintf(`{
		"event": {
			"type": "insert",
			"data": {
				"sample_id": %d
			}
		}
	}`, sampleID))

	task := asynq.NewTask("redraw_order", taskPayload)

	// Execute the handler
	err = processor.HandleRedrawOrderEvent(ctx, task)
	assert.NoError(t, err, "Expected no error when handling valid redraw order event")

	// Verify the order was flagged correctly
	order, err := processor.dbClient.OrderInfo.Query().
		Where(orderinfo.IDEQ(testOrderID)).
		WithOrderFlags().
		First(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Check if the order has the "order_redraw_order" flag
	hasRedrawFlag := false
	for _, flag := range order.Edges.OrderFlags {
		if flag.OrderFlagName == "order_redraw_order" {
			hasRedrawFlag = true
			break
		}
	}
	assert.True(t, hasRedrawFlag, "Order should have the 'order_redraw_order' flag after processing")
}

func TestHandleEditOrderEvent1(t *testing.T) {
	processor, server, asynqClient := setupSampleProcessorTest(t)
	defer cleanUpSampleProcessorTest(processor, server, asynqClient)

	// 🔹 Step 1: Create a test Kafka message (mocked JSON payload)
	testKafkaMessage := []byte(`{
		"sample_id": 2262914,
		"julien_barcode": "2503076277",
		"add_on_test_list": [8, 10, 57],
		"delete_test_list": [],
		"order_id": "30139923543128369",
		"new_tube_info": {
			"actual_number_of_tubes": {"SST": 2, "EDTA": 1},
			"noOfDbsBloodTubes": {},
			"volumeRequired": {"SST": 2123, "EDTA": 1410},
			"noOfTubes": {"EDTA": 1, "SST": 2},
			"tube_order_map": {"EDTA": 1, "SST": 4},
			"actual_volume_required": {"EDTA": "410", "SST": "1123"}
		}
	}`)

	// 🔹 Step 2: Unmarshal JSON into `EditOrderEvent`
	event := &pb.EditOrderEvent{}
	err := json.Unmarshal(testKafkaMessage, event)
	if err != nil {
		t.Fatalf("failed to unmarshal message: %v", err)
	}

	// 🔹 Step 3: Create a test Asynq task
	task, err := tasks.NewEditOrderTask(&tasks.EditOrderTask{
		Event: event,
	})
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	// 🔹 Step 4: Insert test Sample data
	createdSample, err := processor.dbClient.Sample.Create().
		SetID(int(event.SampleId)).
		SetAccessionID("test_accession_id").
		SetDelayedHours(0).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// 🔹 Step 5: Insert test Order Info data
	_, err = processor.dbClient.OrderInfo.Create().
		SetSampleID(createdSample.ID).
		SetOrderConfirmationNumber("test_order").
		SetOrderStatus("test_order_status").
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	_, err = processor.dbClient.LabOrderSendHistory.Create().
		SetSampleID(createdSample.ID).
		SetTubeType("SST").
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	_, err = processor.dbClient.LabOrderSendHistory.Create().
		SetSampleID(createdSample.ID).
		SetTubeType("EDTA").
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	tubeType1, err := processor.dbClient.TubeType.Create().
		SetTubeTypeEnum("SST").
		SetTubeName("Tiger OR Red and yellow top SST").
		SetTubeTypeSymbol("S").
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	tubeType2, err := processor.dbClient.TubeType.Create().
		SetTubeTypeEnum("EDTA").
		SetTubeName("EDTA(Lavender)").
		SetTubeTypeSymbol("W").
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	_, err = processor.dbClient.SampleType.Create().
		SetSampleTypeEnum("SERUM").
		SetSampleTypeName("SERUM").
		SetSampleTypeCode("SERUM").
		SetSampleTypeDescription("Sample Type SERUM").
		SetPrimarySampleTypeGroup("Testing SERUM group").
		AddTubeTypes(tubeType1).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	_, err = processor.dbClient.SampleType.Create().
		SetSampleTypeEnum("EDTA").
		SetSampleTypeName("EDTA").
		SetSampleTypeCode("EDTA").
		SetSampleTypeDescription("Sample Type EDTA").
		SetPrimarySampleTypeGroup("Testing EDTA group").
		AddTubeTypes(tubeType2).
		Save(ctx)
	if err != nil {
		t.Fatal(err)
	}

	testIDs := []int32{8, 10}
	for _, testID := range testIDs {

		_, err = processor.dbClient.Test.Create().
			SetID(int(testID)).
			SetTestCode(fmt.Sprintf("Testing test code %d", testID)).
			SetTestName(fmt.Sprintf("Testing test name %d", testID)).
			SetDisplayName(fmt.Sprintf("Testing test displayName %d", testID)).
			SetTestDescription(fmt.Sprintf("Testing test description %d", testID)).
			SetAssayName(fmt.Sprintf("Testing test assayname %d", testID)).
			SetIsActive(true).
			Save(ctx)
		if err != nil {
			t.Fatalf("failed to create test with ID %d: %v", testID, err)
		}

		_, err = processor.dbClient.TestDetail.Create().
			SetTestID(int(testID)).
			SetTestDetailName("test_sample_type").
			SetTestDetailsValue("Serum").
			Save(ctx)
		if err != nil {
			t.Fatalf("failed to create test detail: %v", err)
		}
	}

	_, err = processor.dbClient.Test.Create().
		SetID(57).
		SetTestCode(fmt.Sprintf("Testing test code %d", 57)).
		SetTestName(fmt.Sprintf("Testing test name %d", 57)).
		SetDisplayName(fmt.Sprintf("Testing test displayName %d", 57)).
		SetTestDescription(fmt.Sprintf("Testing test description %d", 57)).
		SetAssayName(fmt.Sprintf("Testing test assayname %d", 57)).
		SetIsActive(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create test with ID %d: %v", 57, err)
	}

	_, err = processor.dbClient.TestDetail.Create().
		SetTestID(57).
		SetTestDetailName("test_sample_type").
		SetTestDetailsValue("N/A").
		Save(ctx)
	if err != nil {
		t.Fatalf("failed to create test detail: %v", err)
	}

	// 🔹 Step 6: Call `HandleEditOrderEvent`
	err = processor.HandleEditOrderEvent(ctx, task)
	if err != nil {
		t.Fatalf("expected no error, but got: %v", err)
	}

	// 🔹 Step 7: Validate Order Data After Processing
	order, err := processor.dbClient.OrderInfo.Query().Where(
		orderinfo.HasSampleWith(sample.IDEQ(int(event.SampleId))),
	).WithTests().First(ctx)
	if err != nil {
		t.Fatalf("failed to query order info: %v", err)
	}

	// 🔹 Step 8: Check if Add-On Tests were Added
	addedTestIDs := map[int]bool{8: false, 10: false, 57: false}
	for _, test := range order.Edges.Tests {
		if _, exists := addedTestIDs[test.ID]; exists {
			addedTestIDs[test.ID] = true
		}
	}

	for testID, added := range addedTestIDs {
		if !added {
			t.Fatalf("expected test %d to be added, but it was not found", testID)
		}
	}

	mockWriter := publisher.GetPublisher().GetWriter().(*publisher.MockKafkaWriter)
	if len(mockWriter.MessageQueue) < 1 {
		t.Fatalf("expected at least %d Kafka messages, but got %d", 1, len(mockWriter.MessageQueue))
	}

	// Validate the published message
	publishedMessage := mockWriter.MessageQueue[len(mockWriter.MessageQueue)-1]
	publishedEvent := &pb.OrderMessage{}
	if err := json.Unmarshal(publishedMessage.Value, publishedEvent); err != nil {
		t.Fatalf("failed to unmarshal published Kafka message: %v", err)
	}

	// Validate the OrderMessage fields
	if publishedEvent.SampleId != int32(createdSample.ID) {
		t.Fatalf("expected SampleId %d in Kafka message, but got %d", int32(createdSample.ID), publishedEvent.SampleId)
	}

	// Usage in test
	expectedTestIDs := []int32{8, 10}
	if !util.SliceEqual(publishedEvent.TestId, expectedTestIDs) {
		t.Fatalf("expected TestId %v, but got %v", expectedTestIDs, publishedEvent.TestId)
	}

	// // 🔹 Step 9: Check if Remove Tests were Removed
	// deletedTestIDs := map[int]bool{57: true, 107: true}
	// for _, test := range order.Edges.Tests {
	// 	if _, exists := deletedTestIDs[test.ID]; exists {
	// 		t.Fatalf("expected test %d to be removed, but it still exists", test.ID)
	// 	}
	// }

	// 🔹 Step 10: Check if Tube Requirements Were Updated
	tubeRequirements, err := processor.dbClient.TubeRequirement.Query().
		Where(tuberequirement.SampleIDEQ(int(event.SampleId))).
		All(ctx)
	if err != nil {
		t.Fatalf("failed to query tube requirements: %v", err)
	}

	expectedTubes := map[string]int{"EDTA": 1, "SST": 2}
	for _, tubeReq := range tubeRequirements {
		expectedCount, exists := expectedTubes[tubeReq.TubeType]
		if !exists || tubeReq.RequiredCount != expectedCount || tubeReq.RequiredBy != "Order Service Update" {
			t.Fatalf("tube type mismatch or incorrect count for %s: got %d, expected %d", tubeReq.TubeType, tubeReq.RequiredCount, expectedCount)
		}
	}
}
