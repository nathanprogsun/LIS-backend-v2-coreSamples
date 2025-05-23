package tasks

import (
	"github.com/hibiken/asynq"
)

type AsynqClient interface {
	Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error)
	Close() error
}

type MockAsynqClient struct {
	TaskQueue []*asynq.Task
}

func NewMockAsynqClient() AsynqClient {
	return &MockAsynqClient{}
}

func (c *MockAsynqClient) Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	c.TaskQueue = append(c.TaskQueue, task)
	return &asynq.TaskInfo{
			Queue:   "mock_queue",
			Type:    task.Type(),
			Payload: task.Payload(),
		},
		nil
}

func (c *MockAsynqClient) Close() error {
	c.TaskQueue = nil
	return nil
}
