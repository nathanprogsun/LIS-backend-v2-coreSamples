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

type CancelOrderEventHandler struct {
	dbClient               *ent.Client
	ctx                    context.Context
	CancelOrderEventReader *kafka.Reader
	asynqClient            tasks.AsynqClient
}

func (h *CancelOrderEventHandler) HandleCancelOrderEvent(event *pb.CancelOrderEvent) {
	task, _ := tasks.NewCancelSampleOrderTask(&tasks.CancelOrderTask{
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
