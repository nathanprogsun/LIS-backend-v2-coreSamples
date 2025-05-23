package subscriber

import (
	"context"
	"coresamples/common"
	"coresamples/ent"
	pb "coresamples/proto"
	"coresamples/tasks"
	"coresamples/util"
	"fmt"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/hibiken/asynq"
	"github.com/segmentio/kafka-go"
)

type PostOrderEventHandler struct {
	dbClient             *ent.Client
	ctx                  context.Context
	PostOrderEventReader *kafka.Reader
	asynqClient          tasks.AsynqClient
}

type TubeInfo struct {
	TubeNumber        map[string]int32   `json:"noOfTubes,omitempty"`
	VolumeRequired    map[string]float32 `json:"volumeRequired,omitempty"`
	NoOfDbsBloodTubes map[string]int32   `json:"noOfDbsBloodTubes,omitempty"`
	IsDbsPossible     bool               `json:"isDbsPossible,omitempty"`
}

func (h *PostOrderEventHandler) HandlePostOrderEvent(eventKey string, event *pb.PostOrderEvent) {
	sampleId, err := strconv.Atoi(event.SampleId)
	if err != nil {
		common.Error(err)
		return
	}
	customerId, err := strconv.Atoi(event.CustomerId)
	if err != nil {
		common.Error(err)
		return
	}
	//TODO: confirm this
	clinicId := event.ClinicId
	if event.TubeInfo == nil {
		common.Error(fmt.Errorf("Post order: nil tube info %v", event))
		return
	}
	if clinicId != 0 {
		//TODO: try getting clinic info
	}

	//create sample(with processor queue?)
	task, err := tasks.NewPostSampleOrderTask(&tasks.PostSampleOrderTask{
		SampleId:                 int32(sampleId),
		Tests:                    util.Int32ArrayToIntArray(event.OrderContents.Tests),
		PatientId:                event.OrderInfo.PatientId,
		OrderConfirmationNumber:  event.OrderConfirmationNumber,
		CustomerId:               int32(customerId),
		ClinicId:                 clinicId,
		BillingOrderId:           eventKey,
		AccessionId:              strconv.FormatUint(event.OrderInfo.JulienBarcode, 10),
		BloodKitDeliverMethod:    event.OrderInfo.BloodKitDeliveryMethod,
		NonBloodKitDeliverMethod: event.OrderInfo.NonBloodKitDeliveryMethod,
		RequiredNumberOfTubes:    event.TubeInfo.NoOfTubes,
	})
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
