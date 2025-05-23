package subscriber

import (
	"context"
	"coresamples/common"
	"coresamples/ent"
	pb "coresamples/proto"
	"coresamples/tasks"

	"github.com/getsentry/sentry-go"
	"github.com/hibiken/asynq"
	"github.com/segmentio/kafka-go"
)

type ClientTransactionShippingEventHandler struct {
	dbClient                             *ent.Client
	ctx                                  context.Context
	ClientTransactionShippingEventReader *kafka.Reader
	asynqClient                          tasks.AsynqClient
}

func (h *ClientTransactionShippingEventHandler) HandleClientTransactionShippingEvent(event *pb.ClientTransactionShippingEvent) {
	task, _ := tasks.NewClientTransactionShippingTask(&tasks.ClientTransactionShippingTask{
		Event: event,
	})

	// Enqueue the task with Asynq
	taskInfo, err := h.asynqClient.Enqueue(task,
		asynq.Queue(common.Env.AsynqQueueName),
		asynq.MaxRetry(100))
	if err != nil {
		common.Error(err)
		sentry.CaptureException(err)
		return
	}
	common.LogTaskInfo(taskInfo)
}
