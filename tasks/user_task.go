package tasks

import (
	"coresamples/common"
	pb "coresamples/proto"
	"encoding/json"

	"github.com/hibiken/asynq"
)

const (
	TypeUserHubspot = "user:hubspot"
)

type HubspotEventTask struct {
	Event *pb.HubspotEvent
}

func NewHubspotEventTask(task *HubspotEventTask) (*asynq.Task, error) {
	payload, err := json.Marshal(task)
	common.Debugf("hubspot task payload %s, event %v", payload, task.Event)
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeUserHubspot, payload), nil
}
