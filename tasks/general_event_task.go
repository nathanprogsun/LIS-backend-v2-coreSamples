package tasks

import (
	pb "coresamples/proto"
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeSampleOrderGeneralEvent = "general_event:sample_order"
)

// GeneralEventTask defines the structure of a task for handling general events.
type GeneralEventTask struct {
	Event *pb.GeneralEvent
}

// NewGeneralEventTask creates a new Asynq task for handling general events.
func NewSampleOrderGeneralEvent(task *GeneralEventTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSampleOrderGeneralEvent, payload), nil
}
