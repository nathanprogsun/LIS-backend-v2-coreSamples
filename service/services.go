package service

import (
	"coresamples/common"
	"coresamples/ent"
	"coresamples/tasks"
	"fmt"
)

var OrderSvc IOrderService

// TODO: this service uses v1 database,
// remove this after complete migration
var OrderSvcV1 IOrderService
var TestSvc ITestService
var TubeSvc ITubeService
var RBACSvc IRBACService
var SampleSvc ISampleService
var ServiceshipSvc IServiceshipService
var PatientSvc IPatientService
var SalesSvc ISalesService
var UserSvc IUserService
var CustomerSvc ICustomerService

var ErrServiceNotInitialized = fmt.Errorf("service has not been initialized")

type Service struct {
	dbClient    *ent.Client
	redisClient *common.RedisClient
}

func InitService(dbClient *ent.Client, redisClient *common.RedisClient) Service {
	return Service{
		dbClient:    dbClient,
		redisClient: redisClient,
	}
}

func GetOrderService(dbClient *ent.Client, externDbClient *ent.Client, redisClient *common.RedisClient) IOrderService {
	if OrderSvc == nil {
		OrderSvc = NewOrderService(dbClient, redisClient)
	}
	if OrderSvcV1 == nil {
		OrderSvcV1 = NewOrderService(externDbClient, redisClient)
	}
	return OrderSvc
}

func GetCurrentOrderService() IOrderService {
	if OrderSvc == nil {
		common.Fatal(ErrServiceNotInitialized)
	}
	return OrderSvc
}

func GetCurrentOrderServiceV1() IOrderService {
	if OrderSvcV1 == nil {
		common.Fatal(ErrServiceNotInitialized)
	}
	return OrderSvcV1
}

func GetTestService(dbClient *ent.Client, redisClient *common.RedisClient) ITestService {
	if TestSvc == nil {
		TestSvc = NewTestService(dbClient, redisClient)
	}
	return TestSvc
}

func GetCurrentTestService() ITestService {
	if TestSvc == nil {
		common.Fatal(ErrServiceNotInitialized)
	}
	return TestSvc
}

func GetTubeService(dbClient *ent.Client, redisClient *common.RedisClient) ITubeService {
	if TubeSvc == nil {
		TubeSvc = newTubeService(dbClient, redisClient)
	}
	return TubeSvc
}

func GetCurrentTubeService() ITubeService {
	if TubeSvc == nil {
		common.Fatal(ErrServiceNotInitialized)
	}
	return TubeSvc
}

func GetRBACService(dbClient *ent.Client, redisClient *common.RedisClient, driverName string, dataSource string) IRBACService {
	if RBACSvc == nil {
		RBACSvc = newRBACService(dbClient, redisClient, driverName, dataSource)
	}
	return RBACSvc
}

func GetCurrentRBACService() IRBACService {
	if RBACSvc == nil {
		common.Fatal(ErrServiceNotInitialized)
	}
	return RBACSvc
}

func GetSampleService(dbClient *ent.Client, redisClient *common.RedisClient, asynqClient tasks.AsynqClient) ISampleService {
	if SampleSvc == nil {
		SampleSvc = NewSampleService(dbClient, redisClient, asynqClient)
	}
	return SampleSvc
}

func GetCurrentSampleService() ISampleService {
	if SampleSvc == nil {
		common.Fatal(ErrServiceNotInitialized)
	}
	return SampleSvc
}

func GetServiceshipSvc(dbClient *ent.Client, redisClient *common.RedisClient, allowedClinics []string) IServiceshipService {
	if ServiceshipSvc == nil {
		ServiceshipSvc = newServiceshipService(dbClient, redisClient, allowedClinics)
	}
	return ServiceshipSvc
}

func GetCurrentServiceshipSvc() IServiceshipService {
	if ServiceshipSvc == nil {
		common.Fatal(ErrServiceNotInitialized)
	}
	return ServiceshipSvc
}

func GetPatientSvc(dbClient *ent.Client, redisClient *common.RedisClient, secret string) IPatientService {
	if PatientSvc == nil {
		PatientSvc = newPatientService(dbClient, redisClient, secret)
	}
	return PatientSvc
}

func GetCurrentPatientSvc() IPatientService {
	if PatientSvc == nil {
		common.Fatal(ErrServiceNotInitialized)
	}
	return PatientSvc
}

func GetSalesSvc(dbClient *ent.Client, redisClient *common.RedisClient) ISalesService {
	if SalesSvc == nil {
		SalesSvc = newSalesService(dbClient, redisClient)
	}
	return SalesSvc
}

func GetCurrentSalesSvc() ISalesService {
	if SalesSvc == nil {
		common.Fatal(ErrServiceNotInitialized)
	}
	return SalesSvc
}

func GetCustomerService(dbClient *ent.Client, redisClient *common.RedisClient) ICustomerService {
	if CustomerSvc == nil {
		CustomerSvc = NewCustomerService(dbClient, redisClient)
	}

	return CustomerSvc
}

func GetCurrentCustomerService() ICustomerService {
	if CustomerSvc == nil {
		common.Fatal(ErrServiceNotInitialized)
	}
	return CustomerSvc
}
