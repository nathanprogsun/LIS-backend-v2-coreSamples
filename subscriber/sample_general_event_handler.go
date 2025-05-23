package subscriber

import (
	"context"
	"coresamples/common"
	"coresamples/ent"
	pb "coresamples/proto"
	"coresamples/tasks"

	"github.com/getsentry/sentry-go"
	"github.com/hibiken/asynq"
)

type SampleOrderGeneralEventHandler struct {
	dbClient    *ent.Client
	ctx         context.Context
	asynqClient tasks.AsynqClient
}

// HandleSampleOrderGeneralEvent processes a GeneralEvent and enqueues it as an Asynq task.
func (oh *SampleOrderGeneralEventHandler) HandleSampleOrderGeneralEvent(event *pb.GeneralEvent) {
	switch event.EventProvider {
	case "lis-order", "lis-accessioning", "lis-report", "lis-issue-system":
		// Create a new task for the General Event
		task, _ := tasks.NewSampleOrderGeneralEvent(&tasks.GeneralEventTask{
			Event: event,
		})

		// Enqueue the task with Asynq
		taskInfo, err := oh.asynqClient.Enqueue(task,
			asynq.Queue(common.Env.AsynqQueueName),
			asynq.MaxRetry(100))
		if err != nil {
			common.Error(err)
			sentry.CaptureException(err)
			return
		}
		common.LogTaskInfo(taskInfo)

	}
}
