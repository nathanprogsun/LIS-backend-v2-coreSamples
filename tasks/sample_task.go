package tasks

import (
	"coresamples/model"
	pb "coresamples/proto"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

const (
	TypePostSampleOrder               = "sample:create_sample_order"
	TypeFlagOrderOnReceiving          = "sample:flag_order_on_receiving"
	TypeSendOrderOnReceiving          = "sample:send_order_on_receiving"
	TypeSendSampleReceiveGeneralEvent = "sample:send_sample_receive_general_event"
	TypeCancelSampleOrder             = "sample:cancel_order"
	TypeClientTransactionShipping     = "sample:client_transaction_shipping"
	TypeRedrawOrder                   = "sample:redraw_order"
	TypeEditOrder                     = "sample:edit_order"
)

type PostSampleOrderTask struct {
	SampleId                 int32
	PatientId                int32
	OrderConfirmationNumber  string
	CustomerId               int32
	Tests                    []int
	ClinicId                 int32
	BillingOrderId           string
	AccessionId              string
	BloodKitDeliverMethod    string
	NonBloodKitDeliverMethod string
	RequiredNumberOfTubes    map[string]int32
}

type SampleTubeReceiveTask struct {
	SampleId     int32
	TubeDetails  []*model.SampleTubeDetails
	ReceivedTime time.Time
	IsRedraw     bool
}

type CancelOrderTask struct {
	Event *pb.CancelOrderEvent
}

type ClientTransactionShippingTask struct {
	Event *pb.ClientTransactionShippingEvent
}

type RedrawOrderInfoTask struct {
	Event *pb.RedrawOrderInfoEvent
}

type EditOrderTask struct {
	Event *pb.EditOrderEvent
}

func NewPostSampleOrderTask(task *PostSampleOrderTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypePostSampleOrder, payload), nil
}

func NewFlagOrderOnReceivingTask(task *SampleTubeReceiveTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeFlagOrderOnReceiving, payload), nil
}

func NewSendOrderOnReceivingTask(task *SampleTubeReceiveTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSendOrderOnReceiving, payload), nil
}

func NewSendSampleReceiveGeneralEventTask(task *pb.GeneralEvent) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSendSampleReceiveGeneralEvent, payload), nil
}

func NewCancelSampleOrderTask(task *CancelOrderTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeCancelSampleOrder, payload), nil
}

func NewClientTransactionShippingTask(task *ClientTransactionShippingTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeClientTransactionShipping, payload), nil
}

func NewRedrawOrderInfoTask(event *RedrawOrderInfoTask) (*asynq.Task, error) {
	payload, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeRedrawOrder, payload), nil
}

func NewEditOrderTask(task *EditOrderTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeEditOrder, payload), nil
}
