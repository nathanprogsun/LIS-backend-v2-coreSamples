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

type HubspotEventHandler struct {
	dbClient           *ent.Client
	ctx                context.Context
	HubspotEventReader *kafka.Reader
	asynqClient        tasks.AsynqClient
}

func (h *HubspotEventHandler) HandleHubspotEvent(event *pb.HubspotEvent) {
	task, err := tasks.NewHubspotEventTask(&tasks.HubspotEventTask{
		Event: event,
	})
	if err != nil {
		common.Error(err)
		return
	}
	taskInfo, err := h.asynqClient.Enqueue(task,
		asynq.Queue(common.Env.AsynqQueueName),
		asynq.MaxRetry(10))

	if err != nil {
		common.Error(err)
		sentry.CaptureException(err)
		return
	}
	common.LogTaskInfo(taskInfo)
}
