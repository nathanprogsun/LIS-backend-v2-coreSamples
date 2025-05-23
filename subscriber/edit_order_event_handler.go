package subscriber

import (
	"context"
	"coresamples/common"
	"coresamples/ent"
	"coresamples/tasks"

	pb "coresamples/proto"

	"github.com/getsentry/sentry-go"
	"github.com/hibiken/asynq"
	"github.com/segmentio/kafka-go"
)

type EditOrderEventHandler struct {
	dbClient             *ent.Client
	ctx                  context.Context
	EditOrderEventReader *kafka.Reader
	asynqClient          tasks.AsynqClient
}

func (h *EditOrderEventHandler) HandleEditOrderEvent(event *pb.EditOrderEvent) {
	task, _ := tasks.NewEditOrderTask(&tasks.EditOrderTask{
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
