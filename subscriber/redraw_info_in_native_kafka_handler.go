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

const (
	asynqMaxRery = 100
)

type RedrawOrderInfoInNativeKafkaHandler struct {
	dbClient                           *ent.Client
	ctx                                context.Context
	RedrawOrderInfoInNativeKafkaReader *kafka.Reader
	asynqClient                        tasks.AsynqClient
}

func (h *RedrawOrderInfoInNativeKafkaHandler) HandleRedrawOrderInfoFromKafka(event *pb.RedrawOrderInfoEvent) {
	task, _ := tasks.NewRedrawOrderInfoTask(&tasks.RedrawOrderInfoTask{
		Event: event,
	})

	// Enqueue the task with Asynq
	taskInfo, err := h.asynqClient.Enqueue(task,
		asynq.Queue(common.Env.AsynqQueueName),
		asynq.MaxRetry(asynqMaxRery))
	if err != nil {
		common.Error(err)
		sentry.CaptureException(err)
		return
	}
	common.LogTaskInfo(taskInfo)
}
