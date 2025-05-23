package processor

import (
	"coresamples/common"
	"coresamples/ent"
	"coresamples/tasks"
	"github.com/go-redsync/redsync/v4"
	"github.com/hibiken/asynq"
)

func RegisterSampleProcessor(mux *asynq.ServeMux, dbClient *ent.Client, redisClient *common.RedisClient, rs *redsync.Redsync) {
	processor := NewSampleProcessor(dbClient, redisClient, rs)
	mux.HandleFunc(tasks.TypePostSampleOrder, processor.HandlePostSampleOrder)
	mux.HandleFunc(tasks.TypeFlagOrderOnReceiving, processor.HandleFlagOrderOnReceiving)
	mux.HandleFunc(tasks.TypeSendOrderOnReceiving, processor.HandleSendOrderOnReceiving)
	mux.HandleFunc(tasks.TypeSendSampleReceiveGeneralEvent, processor.HandleSendSampleReceiveGeneralEvent)
	mux.HandleFunc(tasks.TypeSampleOrderGeneralEvent, processor.HandleSampleOrderGeneralEvent)
	mux.HandleFunc(tasks.TypeCancelSampleOrder, processor.HandleCancelOrderEvent)
	mux.HandleFunc(tasks.TypeClientTransactionShipping, processor.HandleClientTransactionShipping)
	mux.HandleFunc(tasks.TypeRedrawOrder, processor.HandleRedrawOrderEvent)
	mux.HandleFunc(tasks.TypeEditOrder, processor.HandleEditOrderEvent)
}

func RegisterUserProcessor(mux *asynq.ServeMux, dbClient *ent.Client, redisClient *common.RedisClient) {
	processor := NewUserProcessor(dbClient, redisClient)
	mux.HandleFunc(tasks.TypeUserHubspot, processor.HandleHubspotEvent)
}

func RegisterCDCUpdateProcessor(mux *asynq.ServeMux, dbClient *ent.Client, redisClient *common.RedisClient, rs *redsync.Redsync) {
	processor := NewCDCProcessor(dbClient, redisClient, rs)
	mux.HandleFunc(tasks.TypeAddressCDCUpdate, processor.HandleAddressUpdates)
	mux.HandleFunc(tasks.TypeClinicCDCUpdate, processor.HandleClinicUpdates)
	mux.HandleFunc(tasks.TypeContactCDCUpdate, processor.HandleContactUpdates)
	mux.HandleFunc(tasks.TypeCustomerCDCUpdate, processor.HandleCustomerUpdates)
	mux.HandleFunc(tasks.TypeInternalUserCDCUpdate, processor.HandleInternalUserUpdates)
	mux.HandleFunc(tasks.TypePatientCDCUpdate, processor.HandlePatientUpdates)
	mux.HandleFunc(tasks.TypeSettingCDCUpdate, processor.HandleSettingUpdates)
	mux.HandleFunc(tasks.TypeUserCDCUpdate, processor.HandleUserUpdates)
	mux.HandleFunc(tasks.TypeCustomerToPatientCDCUpdate, processor.HandleCustomerToPatientUpdates)
	mux.HandleFunc(tasks.TypeCustomerSettingOnClinicsCDCUpdate, processor.HandleCustomerSettingOnClinicsUpdates)
	mux.HandleFunc(tasks.TypeClinicToCustomerCDCUpdate, processor.HandleClinicToCustomerUpdates)
	mux.HandleFunc(tasks.TypeClinicToPatientCDCUpdate, processor.HandleClinicToPatientUpdates)
	mux.HandleFunc(tasks.TypeClinicToSettingCDCUpdate, processor.HandleClinicToSettingUpdates)
}
