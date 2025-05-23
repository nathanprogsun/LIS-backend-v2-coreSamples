package subscriber

import (
	"context"
	"coresamples/common"
	"coresamples/ent"
	pb "coresamples/proto"
	"encoding/json"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/hibiken/asynq"
	"github.com/segmentio/kafka-go"
)

var EvSubRecoverCount = 0

const (
	SubscriberGroupIDPrefix = "lis_coresamplesv2_"
)

type EventHandler struct {
	GeneralEventHandler
	PostOrderEventHandler
	CancelOrderEventHandler
	ClientTransactionShippingEventHandler
	RedrawOrderInfoInNativeKafkaHandler
	EditOrderEventHandler
	HubspotEventHandler
	evSubRecoverCount int
}

func (eh *EventHandler) Run() {
	go eh.handleGeneralEvents()
	go eh.handlePostOrderEvents()
	go eh.handleCancelOrderEvent()
	go eh.handleClientTransactionShippingEvent()
	go eh.handleRedrawOrderEvent()
	go eh.handleEditOrderEvent()
	go eh.handleHubspotEvents()
}

func (h *EventHandler) recoverEventSubscriber(eventHandler func(), reader *kafka.Reader) {
	if r := recover(); r != nil {
		common.Error(fmt.Errorf("panic in event subscriber %v\n", r))
		sentry.CaptureMessage("event subscriber panicked")
		h.evSubRecoverCount += 1
		if h.evSubRecoverCount == 10 {
			h.evSubRecoverCount = 0
			sentry.CaptureMessage("watcher recovers too many times")
			reader.Close()
		} else {
			go eventHandler()
		}
	} else {
		reader.Close()
	}
}

func (eh *EventHandler) handleGeneralEvents() {
	defer eh.recoverEventSubscriber(eh.handleGeneralEvents, eh.GeneralEventReader)
	//err := eh.GeneralEventReader.SetOffsetAt(eh.ctx, time.Now())
	//if err != nil {
	//	common.Fatal(err)
	//}
	for {
		m, err := eh.GeneralEventReader.ReadMessage(context.Background())
		if err != nil {
			common.Error(err)
			sentry.CaptureMessage(err.Error())
			continue
		}
		event := &pb.GeneralEvent{}
		err = json.Unmarshal(m.Value, event)
		if err != nil {
			common.Debugf("General events handler: unable to unmarshal event %s, %v", string(m.Value), err)
			continue
		}
		eh.HandleMembershipGeneralEvent(event)
		eh.HandleSampleOrderGeneralEvent(event)
	}
}

func (eh *EventHandler) handlePostOrderEvents() {
	defer eh.recoverEventSubscriber(eh.handlePostOrderEvents, eh.PostOrderEventReader)
	//err := eh.GeneralEventReader.SetOffsetAt(eh.ctx, time.Now())
	//if err != nil {
	//	common.Fatal(err)
	//}
	for {
		m, err := eh.PostOrderEventReader.ReadMessage(context.Background())
		if err != nil {
			common.Error(err)
			sentry.CaptureMessage(err.Error())
			continue
		}
		event := &pb.PostOrderEvent{}
		err = json.Unmarshal(m.Value, event)
		if err != nil {
			common.Error(err)
			continue
		}
		eh.HandlePostOrderEvent(string(m.Key), event)
	}
}

func (eh *EventHandler) handleCancelOrderEvent() {
	defer eh.recoverEventSubscriber(eh.handleCancelOrderEvent, eh.CancelOrderEventReader)
	//err := eh.GeneralEventReader.SetOffsetAt(eh.ctx, time.Now())
	//if err != nil {
	//	common.Fatal(err)
	//}
	for {
		m, err := eh.CancelOrderEventReader.ReadMessage(context.Background())
		if err != nil {
			common.Error(err)
			sentry.CaptureMessage(err.Error())
			continue
		}
		event := &pb.CancelOrderEvent{}
		err = json.Unmarshal(m.Value, event)
		if err != nil {
			common.Error(err)
			continue
		}
		eh.HandleCancelOrderEvent(event)
	}
}

func (eh *EventHandler) handleClientTransactionShippingEvent() {
	defer eh.recoverEventSubscriber(eh.handleClientTransactionShippingEvent, eh.ClientTransactionShippingEventReader)
	//err := eh.GeneralEventReader.SetOffsetAt(eh.ctx, time.Now())
	//if err != nil {
	//	common.Fatal(err)
	//}
	for {
		m, err := eh.ClientTransactionShippingEventReader.ReadMessage(context.Background())
		if err != nil {
			common.Error(err)
			sentry.CaptureMessage(err.Error())
			continue
		}
		event := &pb.ClientTransactionShippingEvent{}
		err = json.Unmarshal(m.Value, event)
		if err != nil {
			common.Error(err)
			continue
		}
		eh.HandleClientTransactionShippingEvent(event)
	}
}

func (eh *EventHandler) handleRedrawOrderEvent() {
	// release kafka reader resource when event handle exit
	defer eh.recoverEventSubscriber(eh.handleRedrawOrderEvent, eh.RedrawOrderInfoInNativeKafkaReader)

	for {
		// Read event from kafka
		m, err := eh.RedrawOrderInfoInNativeKafkaReader.ReadMessage(context.Background())
		if err != nil {
			common.Error(err)
			sentry.CaptureMessage(err.Error())
			continue
		}

		// Deserialize the pb data to the event struct
		event := &pb.RedrawOrderInfoEvent{}
		err = json.Unmarshal(m.Value, event)
		if err != nil {
			common.Error(err)
			continue
		}

		eh.HandleRedrawOrderInfoFromKafka(event)
	}
}

func (eh *EventHandler) handleEditOrderEvent() {
	defer eh.recoverEventSubscriber(eh.handleEditOrderEvent, eh.EditOrderEventReader)

	for {
		m, err := eh.EditOrderEventReader.ReadMessage(context.Background())
		if err != nil {
			common.Error(err)
			sentry.CaptureMessage(err.Error())
			continue
		}
		event := &pb.EditOrderEvent{}
		err = json.Unmarshal(m.Value, event)
		if err != nil {
			common.Error(err)
			continue
		}
		eh.HandleEditOrderEvent(event)
	}
}

func (eh *EventHandler) handleHubspotEvents() {
	defer eh.recoverEventSubscriber(eh.handleHubspotEvents, eh.HubspotEventReader)
	for {
		m, err := eh.HubspotEventReader.ReadMessage(context.Background())
		if err != nil {
			common.Error(err)
			sentry.CaptureMessage(err.Error())
			continue
		}
		event := &pb.HubspotEvent{}
		err = json.Unmarshal(m.Value, event)
		if err != nil {
			common.Error(err)
			continue
		}
		eh.HandleHubspotEvent(event)
	}
}

func NewEventHandler(dbClient *ent.Client, asynqClient *asynq.Client, ctx context.Context, addrs []string, dialer *kafka.Dialer) *EventHandler {
	generalEventReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        addrs,
		Topic:          common.LocalKafkaConfigs.TopicGeneralEvent,
		SessionTimeout: 100 * time.Second,
		Dialer:         dialer,
		GroupID:        SubscriberGroupIDPrefix + common.LocalKafkaConfigs.GroupIDGeneralEvent,
		ErrorLogger:    kafka.LoggerFunc(common.ErrorLogger),
		StartOffset:    kafka.LastOffset,
	})
	postOrderEventReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        addrs,
		Topic:          common.LocalKafkaConfigs.TopicPostOrder,
		SessionTimeout: 100 * time.Second,
		Dialer:         dialer,
		GroupID:        SubscriberGroupIDPrefix + common.LocalKafkaConfigs.GroupIDPostOrder,
		ErrorLogger:    kafka.LoggerFunc(common.ErrorLogger),
		StartOffset:    kafka.LastOffset,
	})
	cancelOrderEventReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        addrs,
		Topic:          common.LocalKafkaConfigs.TopicCancelOrder,
		SessionTimeout: 100 * time.Second,
		Dialer:         dialer,
		GroupID:        SubscriberGroupIDPrefix + common.LocalKafkaConfigs.GroupIDCancelOrder,
		ErrorLogger:    kafka.LoggerFunc(common.ErrorLogger),
		StartOffset:    kafka.LastOffset,
	})
	ClientTransactionShippingEventReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        addrs,
		Topic:          common.LocalKafkaConfigs.TopicTransactionShipping,
		SessionTimeout: 100 * time.Second,
		Dialer:         dialer,
		GroupID:        SubscriberGroupIDPrefix + common.LocalKafkaConfigs.GroupIDTransactionShipping,
		ErrorLogger:    kafka.LoggerFunc(common.ErrorLogger),
		StartOffset:    kafka.LastOffset,
	})
	redrawOrderInfoEventReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        addrs,
		Topic:          common.LocalKafkaConfigs.TopicRedrawOrder,
		SessionTimeout: 100 * time.Second, //the source code value is 95s
		Dialer:         dialer,
		GroupID:        SubscriberGroupIDPrefix + common.LocalKafkaConfigs.GroupIDRedrawOrder,
		ErrorLogger:    kafka.LoggerFunc(common.ErrorLogger),
		StartOffset:    kafka.LastOffset,
	})
	editOrderEventReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        addrs,
		Topic:          common.LocalKafkaConfigs.TopicEditOrder,
		SessionTimeout: 100 * time.Second, //the source code value is 95s
		Dialer:         dialer,
		GroupID:        SubscriberGroupIDPrefix + common.LocalKafkaConfigs.GroupIDEditOrder,
		ErrorLogger:    kafka.LoggerFunc(common.ErrorLogger),
		StartOffset:    kafka.LastOffset,
	})
	HubspotEventReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        addrs,
		Topic:          common.LocalKafkaConfigs.TopicHubspot,
		SessionTimeout: 100 * time.Second,
		Dialer:         dialer,
		GroupID:        SubscriberGroupIDPrefix + common.LocalKafkaConfigs.GroupIDHubspotEvent,
		ErrorLogger:    kafka.LoggerFunc(common.ErrorLogger),
		StartOffset:    kafka.LastOffset,
	})
	return &EventHandler{
		GeneralEventHandler: GeneralEventHandler{
			MembershipEventHandler: MembershipEventHandler{
				dbClient: dbClient,
				ctx:      ctx,
			},
			SampleOrderGeneralEventHandler: SampleOrderGeneralEventHandler{
				dbClient:    dbClient,
				ctx:         ctx,
				asynqClient: asynqClient,
			},
			GeneralEventReader: generalEventReader,
		},
		PostOrderEventHandler: PostOrderEventHandler{
			dbClient:             dbClient,
			ctx:                  ctx,
			PostOrderEventReader: postOrderEventReader,
			asynqClient:          asynqClient,
		},
		CancelOrderEventHandler: CancelOrderEventHandler{
			dbClient:               dbClient,
			ctx:                    ctx,
			CancelOrderEventReader: cancelOrderEventReader,
			asynqClient:            asynqClient,
		},
		ClientTransactionShippingEventHandler: ClientTransactionShippingEventHandler{
			dbClient:                             dbClient,
			ctx:                                  ctx,
			ClientTransactionShippingEventReader: ClientTransactionShippingEventReader,
			asynqClient:                          asynqClient,
		},
		RedrawOrderInfoInNativeKafkaHandler: RedrawOrderInfoInNativeKafkaHandler{
			dbClient:                           dbClient,
			ctx:                                ctx,
			RedrawOrderInfoInNativeKafkaReader: redrawOrderInfoEventReader,
		},
		EditOrderEventHandler: EditOrderEventHandler{
			dbClient:             dbClient,
			ctx:                  ctx,
			EditOrderEventReader: editOrderEventReader,
			asynqClient:          asynqClient,
		},
		HubspotEventHandler: HubspotEventHandler{
			dbClient:           dbClient,
			ctx:                ctx,
			HubspotEventReader: HubspotEventReader,
			asynqClient:        asynqClient,
		},
	}
}
