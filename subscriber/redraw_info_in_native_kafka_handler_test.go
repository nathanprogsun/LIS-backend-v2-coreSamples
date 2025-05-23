package subscriber

import (
	"context"
	"testing"

	"coresamples/common"
	pb "coresamples/proto"
	"coresamples/tasks"
)

// setupRedrawHandlerTest creates the test environment
func setupRedrawHandlerTest(t *testing.T) (*RedrawOrderInfoInNativeKafkaHandler, *tasks.MockAsynqClient, context.Context) {
	// Use the existing mock implementation but keep a reference to the concrete type
	mockClient := tasks.NewMockAsynqClient().(*tasks.MockAsynqClient)
	ctx := context.Background()

	common.InitZapLogger("debug")

	handler := &RedrawOrderInfoInNativeKafkaHandler{
		dbClient:    nil, // Not needed for this test
		ctx:         ctx,
		asynqClient: mockClient,
	}

	return handler, mockClient, ctx
}

func cleanUpRedrawHandlerTest(client *tasks.MockAsynqClient) {
	_ = client.Close()
}

func TestHandleRedrawOrderInfoFromKafka(t *testing.T) {
	handler, mockClient, _ := setupRedrawHandlerTest(t)
	defer cleanUpRedrawHandlerTest(mockClient)

	// Setup test data
	event := &pb.RedrawOrderInfoEvent{
		Database: "test database",
		Data: &pb.RedrawOrderInfoEvent_Data{
			SampleId: 111, // Using int32 instead of string
		},
	}

	// Call the method under test
	handler.HandleRedrawOrderInfoFromKafka(event)

	// Verify the task was enqueued
	if len(mockClient.TaskQueue) != 1 {
		t.Fatalf("Expected one task to be enqueued, got %d", len(mockClient.TaskQueue))
	}

	// Verify the task type - update to match the actual task type
	expectedTaskType := "sample:redraw_order"
	if mockClient.TaskQueue[0].Type() != expectedTaskType {
		t.Errorf("Expected task type '%s', got '%s'", expectedTaskType, mockClient.TaskQueue[0].Type())
	}
}
