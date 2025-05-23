package subscriber

import (
	"context"
	"coresamples/common"
	pb "coresamples/proto"
	"coresamples/tasks"
	"encoding/json"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/hibiken/asynq"
	"github.com/segmentio/kafka-go"
)

type CoreCDCUpdatesHandler struct {
	AddressCDCUpdateReader                  *kafka.Reader
	ClinicCDCUpdateReader                   *kafka.Reader
	ContactCDCUpdateReader                  *kafka.Reader
	CustomerCDCUpdateReader                 *kafka.Reader
	InternalUserCDCUpdateReader             *kafka.Reader
	PatientCDCUpdateReader                  *kafka.Reader
	SettingCDCUpdateReader                  *kafka.Reader
	UserCDCUpdateReader                     *kafka.Reader
	CustomerToPatientCDCUpdateReader        *kafka.Reader
	CustomerSettingOnClinicsCDCUpdateReader *kafka.Reader
	ClinicToCustomerCDCUpdateReader         *kafka.Reader
	ClinicToPatientCDCUpdateReader          *kafka.Reader
	ClinicToSettingCDCUpdateReader          *kafka.Reader
	asynqClient                             tasks.AsynqClient
	evSubRecoverCount                       int
}

func NewCoreCDCUpdatesHandler(addrs []string, dialer *kafka.Dialer, asynqClient tasks.AsynqClient) *CoreCDCUpdatesHandler {
	addressReader := newLISCoreCDCReader(common.LocalKafkaConfigs.TopicLISCoreCDCAddress, "address", addrs, dialer)
	clinicReader := newLISCoreCDCReader(common.LocalKafkaConfigs.TopicLISCoreCDCClinic, "clinic", addrs, dialer)
	contactReader := newLISCoreCDCReader(common.LocalKafkaConfigs.TopicLISCoreCDCContact, "contact", addrs, dialer)
	customerReader := newLISCoreCDCReader(common.LocalKafkaConfigs.TopicLISCoreCDCCustomer, "customer", addrs, dialer)
	internalUserReader := newLISCoreCDCReader(common.LocalKafkaConfigs.TopicLISCoreCDCInternalUser, "internal_user", addrs, dialer)
	patientReader := newLISCoreCDCReader(common.LocalKafkaConfigs.TopicLISCoreCDCPatient, "patient", addrs, dialer)
	settingReader := newLISCoreCDCReader(common.LocalKafkaConfigs.TopicLISCoreCDCSetting, "setting", addrs, dialer)
	userReader := newLISCoreCDCReader(common.LocalKafkaConfigs.TopicLISCoreCDCUser, "user", addrs, dialer)
	customerToPatientReader := newLISCoreCDCReader(common.LocalKafkaConfigs.TopicLISCoreCDCCustomerToPatient, "_customertopatient", addrs, dialer)
	customerSettingReader := newLISCoreCDCReader(common.LocalKafkaConfigs.TopicLISCoreCDCCustomerSettingOnClinics, "customersettingonclinics", addrs, dialer)
	clinicToCustomerReader := newLISCoreCDCReader(common.LocalKafkaConfigs.TopicLISCoreCDCClinicToCustomer, "_clinictocustomer", addrs, dialer)
	clinicToPatientReader := newLISCoreCDCReader(common.LocalKafkaConfigs.TopicLISCoreCDCClinicToPatient, "_clinictopatient", addrs, dialer)
	clinicToSettingReader := newLISCoreCDCReader(common.LocalKafkaConfigs.TopicLISCoreCDCClinicToSetting, "_clinictosetting", addrs, dialer)

	return &CoreCDCUpdatesHandler{
		AddressCDCUpdateReader:                  addressReader,
		ClinicCDCUpdateReader:                   clinicReader,
		ContactCDCUpdateReader:                  contactReader,
		CustomerCDCUpdateReader:                 customerReader,
		InternalUserCDCUpdateReader:             internalUserReader,
		PatientCDCUpdateReader:                  patientReader,
		SettingCDCUpdateReader:                  settingReader,
		UserCDCUpdateReader:                     userReader,
		CustomerToPatientCDCUpdateReader:        customerToPatientReader,
		CustomerSettingOnClinicsCDCUpdateReader: customerSettingReader,
		ClinicToCustomerCDCUpdateReader:         clinicToCustomerReader,
		ClinicToPatientCDCUpdateReader:          clinicToPatientReader,
		ClinicToSettingCDCUpdateReader:          clinicToSettingReader,
		asynqClient:                             asynqClient,
	}
}

func (h *CoreCDCUpdatesHandler) recoverEventSubscriber(eventHandler func(), reader *kafka.Reader) {
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

func (h *CoreCDCUpdatesHandler) Run() {
	go h.handleAddressUpdates()
	go h.handleClinicUpdates()
	go h.handleContactUpdates()
	go h.handleCustomerUpdates()
	go h.handleInternalUserUpdates()
	go h.handlePatientUpdates()
	go h.handleSettingUpdates()
	go h.handleUserUpdates()

	go h.handleCustomerSettingOnClinicsUpdates()
	go h.handleCustomerToPatientUpdates()
	go h.handleClinicToCustomerUpdates()
	go h.handleClinicToPatientUpdates()
	go h.handleClinicToSettingUpdates()
}

func (h *CoreCDCUpdatesHandler) handleAddressUpdates() {
	defer h.recoverEventSubscriber(h.handleAddressUpdates, h.AddressCDCUpdateReader)
	for {
		m, err := h.AddressCDCUpdateReader.ReadMessage(context.Background())
		if err != nil {
			common.Error(err)
			continue
		}
		event := &pb.AddressCDCUpdate{}
		err = json.Unmarshal(m.Value, event)
		if err != nil {
			common.Error(err)
			continue
		}
		task, _ := tasks.NewAddressCDCUpdateTask(&tasks.AddressCDCUpdateTask{
			Event: event,
		})
		taskInfo, err := h.asynqClient.Enqueue(task, asynq.Queue(common.Env.AsynqQueueName), asynq.MaxRetry(10))
		if err != nil {
			common.Error(err)
		}
		common.LogTaskInfo(taskInfo)
	}
}

func (h *CoreCDCUpdatesHandler) handleClinicUpdates() {
	defer h.recoverEventSubscriber(h.handleClinicUpdates, h.ClinicCDCUpdateReader)
	for {
		m, err := h.ClinicCDCUpdateReader.ReadMessage(context.Background())
		if err != nil {
			common.Error(err)
			continue
		}
		event := &pb.ClinicCDCUpdate{}
		err = json.Unmarshal(m.Value, event)
		if err != nil {
			common.Error(err)
			continue
		}
		task, _ := tasks.NewClinicCDCUpdateTask(&tasks.ClinicCDCUpdateTask{
			Event: event,
		})
		taskInfo, err := h.asynqClient.Enqueue(task, asynq.Queue(common.Env.AsynqQueueName), asynq.MaxRetry(10))
		if err != nil {
			common.Error(err)
		}
		common.LogTaskInfo(taskInfo)
	}
}

func (h *CoreCDCUpdatesHandler) handleContactUpdates() {
	defer h.recoverEventSubscriber(h.handleContactUpdates, h.ContactCDCUpdateReader)
	for {
		m, err := h.ContactCDCUpdateReader.ReadMessage(context.Background())
		if err != nil {
			common.Error(err)
			continue
		}
		event := &pb.ContactCDCUpdate{}
		err = json.Unmarshal(m.Value, event)
		if err != nil {
			common.Error(err)
			continue
		}
		task, _ := tasks.NewContactCDCUpdateTask(&tasks.ContactCDCUpdateTask{
			Event: event,
		})
		taskInfo, err := h.asynqClient.Enqueue(task, asynq.Queue(common.Env.AsynqQueueName), asynq.MaxRetry(10))
		if err != nil {
			common.Error(err)
		}
		common.LogTaskInfo(taskInfo)
	}
}

func (h *CoreCDCUpdatesHandler) handleCustomerUpdates() {
	defer h.recoverEventSubscriber(h.handleCustomerUpdates, h.CustomerCDCUpdateReader)
	for {
		m, err := h.CustomerCDCUpdateReader.ReadMessage(context.Background())
		if err != nil {
			common.Error(err)
			sentry.CaptureMessage(err.Error())
			continue
		}
		event := &pb.CustomerCDCUpdate{}
		err = json.Unmarshal(m.Value, event)
		if err != nil {
			common.Error(err)
			continue
		}
		task, _ := tasks.NewCustomerCDCUpdateTask(&tasks.CustomerCDCUpdateTask{
			Event: event,
		})

		// Enqueue the task with Asynq
		taskInfo, err := h.asynqClient.Enqueue(task,
			asynq.Queue(common.Env.AsynqQueueName),
			asynq.MaxRetry(100))
		if err != nil {
			common.Error(err)
			sentry.CaptureException(err)
			continue
		}
		common.LogTaskInfo(taskInfo)
	}
}

func (h *CoreCDCUpdatesHandler) handleInternalUserUpdates() {
	defer h.recoverEventSubscriber(h.handleInternalUserUpdates, h.InternalUserCDCUpdateReader)
	for {
		m, err := h.InternalUserCDCUpdateReader.ReadMessage(context.Background())
		if err != nil {
			common.Error(err)
			sentry.CaptureMessage(err.Error())
			continue
		}
		event := &pb.InternalUserCDCUpdate{}
		err = json.Unmarshal(m.Value, event)
		if err != nil {
			common.Error(err)
			continue
		}
		task, _ := tasks.NewInternalUserCDCUpdateTask(&tasks.InternalUserCDCUpdateTask{
			Event: event,
		})

		// Enqueue the task with Asynq
		taskInfo, err := h.asynqClient.Enqueue(task,
			asynq.Queue(common.Env.AsynqQueueName),
			asynq.MaxRetry(100))
		if err != nil {
			common.Error(err)
			sentry.CaptureException(err)
			continue
		}
		common.LogTaskInfo(taskInfo)
	}
}

func (h *CoreCDCUpdatesHandler) handlePatientUpdates() {
	defer h.recoverEventSubscriber(h.handlePatientUpdates, h.PatientCDCUpdateReader)
	for {
		m, err := h.PatientCDCUpdateReader.ReadMessage(context.Background())
		if err != nil {
			common.Error(err)
			continue
		}
		event := &pb.PatientCDCUpdate{}
		err = json.Unmarshal(m.Value, event)
		if err != nil {
			common.Error(err)
			continue
		}
		task, _ := tasks.NewPatientCDCUpdateTask(&tasks.PatientCDCUpdateTask{
			Event: event,
		})
		taskInfo, err := h.asynqClient.Enqueue(task, asynq.Queue(common.Env.AsynqQueueName), asynq.MaxRetry(10))
		if err != nil {
			common.Error(err)
		}
		common.LogTaskInfo(taskInfo)
	}
}

func (h *CoreCDCUpdatesHandler) handleSettingUpdates() {
	defer h.recoverEventSubscriber(h.handleSettingUpdates, h.SettingCDCUpdateReader)

	for {
		m, err := h.SettingCDCUpdateReader.ReadMessage(context.Background())
		if err != nil {
			common.Error(err)
			sentry.CaptureMessage(err.Error())
			continue
		}
		event := &pb.SettingCDCUpdate{}
		err = json.Unmarshal(m.Value, event)
		if err != nil {
			common.Error(err)
			continue
		}
		task, _ := tasks.NewSettingCDCUpdateTask(&tasks.SettingCDCUpdateTask{
			Event: event,
		})

		// Enqueue the task with Asynq
		taskInfo, err := h.asynqClient.Enqueue(task,
			asynq.Queue(common.Env.AsynqQueueName),
			asynq.MaxRetry(100))
		if err != nil {
			common.Error(err)
			sentry.CaptureException(err)
		}
		common.LogTaskInfo(taskInfo)
	}
}

func (h *CoreCDCUpdatesHandler) handleUserUpdates() {
	defer h.recoverEventSubscriber(h.handleUserUpdates, h.UserCDCUpdateReader)
	for {
		m, err := h.UserCDCUpdateReader.ReadMessage(context.Background())
		if err != nil {
			common.Error(err)
			continue
		}
		event := &pb.UserCDCUpdate{}
		err = json.Unmarshal(m.Value, event)
		if err != nil {
			common.Error(err)
			continue
		}
		task, _ := tasks.NewUserCDCUpdateTask(&tasks.UserCDCUpdateTask{
			Event: event,
		})
		taskInfo, err := h.asynqClient.Enqueue(task, asynq.Queue(common.Env.AsynqQueueName), asynq.MaxRetry(10))
		if err != nil {
			common.Error(err)
		}
		common.LogTaskInfo(taskInfo)
	}
}

func (h *CoreCDCUpdatesHandler) handleCustomerToPatientUpdates() {
	defer h.recoverEventSubscriber(h.handleCustomerToPatientUpdates, h.CustomerToPatientCDCUpdateReader)
	for {
		m, err := h.CustomerToPatientCDCUpdateReader.ReadMessage(context.Background())
		if err != nil {
			common.Error(err)
			continue
		}
		event := &pb.CustomerToPatientCDCUpdate{}
		err = json.Unmarshal(m.Value, event)
		if err != nil {
			common.Error(err)
			continue
		}
		task, _ := tasks.NewCustomerToPatientCDCUpdateTask(&tasks.CustomerToPatientCDCUpdateTask{
			Event: event,
		})
		taskInfo, err := h.asynqClient.Enqueue(task, asynq.Queue(common.Env.AsynqQueueName), asynq.MaxRetry(10))
		if err != nil {
			common.Error(err)
		}
		common.LogTaskInfo(taskInfo)
	}
}

func (h *CoreCDCUpdatesHandler) handleCustomerSettingOnClinicsUpdates() {
	defer h.recoverEventSubscriber(h.handleCustomerSettingOnClinicsUpdates, h.CustomerSettingOnClinicsCDCUpdateReader)

	for {
		m, err := h.CustomerSettingOnClinicsCDCUpdateReader.ReadMessage(context.Background())
		if err != nil {
			common.Error(err)
			sentry.CaptureMessage(err.Error())
			continue
		}
		event := &pb.CustomerSettingOnClinicsCDCUpdate{}
		err = json.Unmarshal(m.Value, event)
		if err != nil {
			common.Error(err)
			continue
		}
		task, _ := tasks.NewCustomerSettingOnClinicsCDCUpdateTask(&tasks.CustomerSettingOnClinicsCDCUpdateTask{
			Event: event,
		})

		// Enqueue the task with Asynq
		taskInfo, err := h.asynqClient.Enqueue(task,
			asynq.Queue(common.Env.AsynqQueueName),
			asynq.MaxRetry(100))
		if err != nil {
			common.Error(err)
			sentry.CaptureException(err)
		}
		common.LogTaskInfo(taskInfo)
	}
}

func (h *CoreCDCUpdatesHandler) handleClinicToCustomerUpdates() {
	defer h.recoverEventSubscriber(h.handleClinicToCustomerUpdates, h.ClinicToCustomerCDCUpdateReader)
	for {
		m, err := h.ClinicToCustomerCDCUpdateReader.ReadMessage(context.Background())
		if err != nil {
			common.Error(err)
			continue
		}
		event := &pb.ClinicToCustomerCDCUpdate{}
		err = json.Unmarshal(m.Value, event)
		if err != nil {
			common.Error(err)
			continue
		}
		task, _ := tasks.NewClinicToCustomerCDCUpdateTask(&tasks.ClinicToCustomerCDCUpdateTask{
			Event: event,
		})
		taskInfo, err := h.asynqClient.Enqueue(task, asynq.Queue(common.Env.AsynqQueueName), asynq.MaxRetry(10))
		if err != nil {
			common.Error(err)
		}
		common.LogTaskInfo(taskInfo)
	}
}

func (h *CoreCDCUpdatesHandler) handleClinicToPatientUpdates() {
	defer h.recoverEventSubscriber(h.handleClinicToPatientUpdates, h.ClinicToPatientCDCUpdateReader)
	for {
		m, err := h.ClinicToPatientCDCUpdateReader.ReadMessage(context.Background())
		if err != nil {
			common.Error(err)
			continue
		}
		event := &pb.ClinicToPatientCDCUpdate{}
		err = json.Unmarshal(m.Value, event)
		if err != nil {
			common.Error(err)
			continue
		}
		task, _ := tasks.NewClinicToPatientCDCUpdateTask(&tasks.ClinicToPatientCDCUpdateTask{
			Event: event,
		})
		taskInfo, err := h.asynqClient.Enqueue(task, asynq.Queue(common.Env.AsynqQueueName), asynq.MaxRetry(10))
		if err != nil {
			common.Error(err)
		}
		common.LogTaskInfo(taskInfo)
	}
}

func (h *CoreCDCUpdatesHandler) handleClinicToSettingUpdates() {
	defer h.recoverEventSubscriber(h.handleClinicToSettingUpdates, h.ClinicToSettingCDCUpdateReader)
	for {
		m, err := h.ClinicToSettingCDCUpdateReader.ReadMessage(context.Background())
		if err != nil {
			common.Error(err)
			sentry.CaptureMessage(err.Error())
			continue
		}
		event := &pb.ClinicToSettingCDCUpdate{}
		err = json.Unmarshal(m.Value, event)
		if err != nil {
			common.Error(err)
			continue
		}
		task, _ := tasks.NewClinicToSettingCDCUpdateTask(&tasks.ClinicToSettingCDCUpdateTask{
			Event: event,
		})

		// Enqueue the task with Asynq
		taskInfo, err := h.asynqClient.Enqueue(task,
			asynq.Queue(common.Env.AsynqQueueName),
			asynq.MaxRetry(100))
		if err != nil {
			common.Error(err)
			sentry.CaptureException(err)
		}
		common.LogTaskInfo(taskInfo)
	}

}

func newLISCoreCDCReader(topic string, suffix string, addrs []string, dialer *kafka.Dialer) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        addrs,
		Topic:          topic,
		SessionTimeout: 100 * time.Second,
		Dialer:         dialer,
		GroupID:        SubscriberGroupIDPrefix + common.LocalKafkaConfigs.GroupIDLISCoreCDC + "_" + suffix,
		ErrorLogger:    kafka.LoggerFunc(common.ErrorLogger),
		StartOffset:    kafka.LastOffset,
	})
}
