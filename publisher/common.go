package publisher

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type KafkaWriter interface {
	Close() error
	WriteMessages(ctx context.Context, msgs ...kafka.Message) error
}

type MockKafkaWriter struct {
	MessageQueue []kafka.Message
}

func (w *MockKafkaWriter) Close() error {
	w.MessageQueue = nil
	return nil
}

func (w *MockKafkaWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	w.MessageQueue = append(w.MessageQueue, msgs...)
	return nil
}
