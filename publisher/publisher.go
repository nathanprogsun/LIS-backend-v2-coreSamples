package publisher

import (
	"context"
	"coresamples/common"
	pb "coresamples/proto"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

const (
	EmailFrom = "notification@vibrant-america.com"
)

var publisher *Publisher

type Publisher struct {
	ctx    context.Context
	writer KafkaWriter
}

func InitPublisher(ctx context.Context, transport *kafka.Transport, addrs []string) {
	if publisher == nil {
		publisher = &Publisher{
			ctx:    ctx,
			writer: initKafkaWriter(transport, addrs),
		}
	}
}

func InitMockPublisher() {
	if publisher == nil {
		publisher = &Publisher{
			ctx:    context.Background(),
			writer: &MockKafkaWriter{},
		}
	}
}

func initKafkaWriter(transport *kafka.Transport, addrs []string) *kafka.Writer {
	var writer *kafka.Writer
	if transport != nil {
		writer = &kafka.Writer{
			Addr:                   kafka.TCP(addrs...),
			Balancer:               &kafka.Hash{},
			Transport:              transport,
			AllowAutoTopicCreation: true,
			Logger:                 kafka.LoggerFunc(common.Infof),
			ErrorLogger:            kafka.LoggerFunc(common.ErrorLogger),
		}
	} else {
		writer = &kafka.Writer{
			Addr:                   kafka.TCP(addrs...),
			Balancer:               &kafka.Hash{},
			AllowAutoTopicCreation: true,
			Logger:                 kafka.LoggerFunc(common.Infof),
			ErrorLogger:            kafka.LoggerFunc(common.ErrorLogger),
		}
	}
	return writer
}

func GetPublisher() *Publisher {
	if publisher == nil {
		common.Fatal(fmt.Errorf("publisher is not initialized"))
	}
	return publisher
}

func (p *Publisher) GetWriter() KafkaWriter {
	return p.writer
}

func (p *Publisher) SendTestMessage() error {
	email := &pb.SubscriptionCancellationEmail{
		MessageID:     time.Now().Format("2006-01-02 15:04:05"),
		Tag:           "Subscription cancellation",
		Subject:       "Longevity Club Membership Cancellation",
		From:          EmailFrom,
		Cc:            "",
		Bcc:           "",
		TemplateId:    8340335,
		Delay:         0,
		MessageStream: "outbound",
		Type:          "Email",
	}
	emailb, err := json.Marshal(email)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Topic: common.LocalKafkaConfigs.TopicEmail,
		Key:   []byte(email.MessageID),
		Value: emailb,
	}
	common.InfoFields("Sending message", zap.String("topic", msg.Topic), zap.String("message", fmt.Sprintf("%v", email)))
	if common.Env.DryRun {
		return nil
	}
	err = p.writer.WriteMessages(p.ctx, msg)
	if err != nil {
		return err
	}
	return nil
}

func (p *Publisher) SendSubscriptionSuccessEmail(template *pb.SubscriptionConfirmationEmailTemplate, to string, subscriptionID string) error {
	email := &pb.SubscriptionConfirmationEmail{
		MessageID:     "subscription " + subscriptionID + "_" + time.Now().Format("2006-01-02 15:04:05"),
		Tag:           "Subscription confirmation",
		Subject:       "Welcome to the Longevity Club!",
		From:          EmailFrom,
		To:            to,
		Cc:            "",
		Bcc:           "",
		TemplateId:    35604717,
		TemplateModel: template,
		Delay:         0,
		MessageStream: "outbound",
		Type:          "Email",
	}
	emailb, err := json.Marshal(email)
	if err != nil {
		common.Error(err)
		return err
	}
	msg := kafka.Message{
		Topic: common.LocalKafkaConfigs.TopicEmail,
		Key:   []byte(email.MessageID),
		Value: emailb,
	}
	common.InfoFields("Sending message", zap.String("topic", msg.Topic), zap.String("message", fmt.Sprintf("%v", email)))
	if common.Env.DryRun {
		return nil
	}
	return p.writer.WriteMessages(p.ctx, msg)
}

func (p *Publisher) SendPaymentUpdateEmail(template *pb.PaymentUpdateEmailTemplate, to string, subscriptionID string) error {
	email := &pb.PaymentUpdateEmail{
		MessageID:     "subscription " + subscriptionID + "_" + time.Now().Format("2006-01-02 15:04:05"),
		Tag:           "Require payment update",
		Subject:       "Action Required: Update Your Payment Information",
		From:          EmailFrom,
		To:            to,
		Cc:            "",
		Bcc:           "",
		TemplateId:    35605030,
		TemplateModel: template,
		Delay:         0,
		MessageStream: "outbound",
		Type:          "Email",
	}
	emailb, err := json.Marshal(email)
	if err != nil {
		common.Error(err)
		return err
	}
	msg := kafka.Message{
		Topic: common.LocalKafkaConfigs.TopicEmail,
		Key:   []byte(email.MessageID),
		Value: emailb,
	}
	common.InfoFields("Sending message", zap.String("topic", msg.Topic), zap.String("message", fmt.Sprintf("%v", email)))
	if common.Env.DryRun {
		return nil
	}
	return p.writer.WriteMessages(p.ctx, msg)
}

func (p *Publisher) SendSubscriptionCancellationEmail(template *pb.SubscriptionCancellationEmailTemplate, to string, subscriptionID string) error {
	email := &pb.SubscriptionCancellationEmail{
		MessageID:     "subscription " + subscriptionID + "_" + time.Now().Format("2006-01-02 15:04:05"),
		Tag:           "Subscription cancellation",
		Subject:       "Longevity Club Membership Cancellation",
		From:          EmailFrom,
		To:            to,
		Cc:            "",
		Bcc:           "",
		TemplateId:    8340335,
		TemplateModel: template,
		Delay:         0,
		MessageStream: "outbound",
		Type:          "Email",
	}
	emailb, err := json.Marshal(email)
	if err != nil {
		common.Error(err)
		return err
	}
	msg := kafka.Message{
		Topic: common.LocalKafkaConfigs.TopicEmail,
		Key:   []byte(email.MessageID),
		Value: emailb,
	}
	common.InfoFields("Sending message", zap.String("topic", msg.Topic), zap.String("message", fmt.Sprintf("%v", email)))
	if common.Env.DryRun {
		return nil
	}
	return p.writer.WriteMessages(p.ctx, msg)
}

func (p *Publisher) SendGeneralEvent(event *pb.GeneralEvent) error {
	// only send message in dev env for test purpose
	if common.Env.RunEnv != "" && common.Env.RunEnv != common.DevEnv {
		return nil
	}
	val, err := json.Marshal(event)
	if err != nil {
		common.Error(err)
		return err
	}
	msg := kafka.Message{
		Topic: common.LocalKafkaConfigs.TopicGeneralEvent,
		Key:   []byte(uuid.NewString()),
		Value: val,
	}
	common.InfoFields("Sending message", zap.String("topic", msg.Topic), zap.String("message", fmt.Sprintf("%v", event)))
	if common.Env.DryRun {
		return nil
	}
	return p.writer.WriteMessages(p.ctx, msg)
}

func (p *Publisher) SendOrderMessage(message *pb.OrderMessage) error {
	// only send message in dev env for test purpose
	if common.Env.RunEnv != "" && common.Env.RunEnv != common.DevEnv {
		return nil
	}
	val, err := json.Marshal(message)
	if err != nil {
		common.Error(err)
		return err
	}
	msg := kafka.Message{
		Topic: common.LocalKafkaConfigs.TopicSampleTest,
		Key:   []byte(strconv.Itoa(int(message.SampleId))),
		Value: val,
	}
	common.InfoFields("Sending message", zap.String("topic", msg.Topic), zap.String("message", fmt.Sprintf("%v", message)))
	if common.Env.DryRun {
		return nil
	}
	return p.writer.WriteMessages(p.ctx, msg)
}
