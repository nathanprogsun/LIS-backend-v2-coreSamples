package handler

import (
	"coresamples/common"
	pb "coresamples/proto"
	coresamplesservice "coresamples/service"

	"github.com/gin-gonic/gin"
	"go-micro.dev/v4/server"
)

var membershipHTTPHandler *ServiceshipHTTPHandler

func RegisterTubeHandler(grpcServer server.Server) {
	err := pb.RegisterTubeServiceHandler(grpcServer, &TubeHandler{
		TubeService: coresamplesservice.GetCurrentTubeService(),
	})
	if err != nil {
		common.Fatal(err)
	}
}

func RegisterRBACHandler(grpcServer server.Server, allowedServices []string) {
	err := pb.RegisterRBACServiceHandler(grpcServer, &RBACHandler{
		RBACService:     coresamplesservice.GetCurrentRBACService(),
		allowedServices: allowedServices,
	})
	if err != nil {
		common.Fatal(err)
	}
}

func RegisterHealthHandler(grpcServer server.Server) {
	err := pb.RegisterHealthHandler(grpcServer, &HealthcheckHandler{})
	if err != nil {
		common.Fatal(err)
	}
}

func RegisterServiceshipHandler(router *gin.RouterGroup) {
	service := coresamplesservice.GetCurrentServiceshipSvc()
	//registerMembershipRPCHandler(grpcServer, service)
	registerMembershipHTTPHandler(router, service)
}

func RegisterTestHandler(router *gin.RouterGroup, grpcServer server.Server) {
	service := coresamplesservice.GetCurrentTestService()
	err := pb.RegisterTestServiceHandler(grpcServer, &TestHandler{
		TestService: service,
	})
	if err != nil {
		common.Fatal(err)
	}

	registerTestHTTPHandler(router, service)
}

func RegisterOrderHandler(router *gin.RouterGroup, grpcServer server.Server) {
	service := coresamplesservice.GetCurrentOrderService()
	err := pb.RegisterOrderServiceHandler(grpcServer, &OrderHandler{
		OrderService: service,
	})
	if err != nil {
		common.Fatal(err)
	}
	registerOrderHTTPHandler(router, coresamplesservice.GetCurrentOrderServiceV1())
}

func RegisterPatientHandler(router *gin.RouterGroup) {
	service := coresamplesservice.GetCurrentPatientSvc()
	registerPatientHTTPHandler(router, service)
}

func RegisterSampleHandler(grpcServer server.Server) {
	service := coresamplesservice.GetCurrentSampleService()
	err := pb.RegisterSampleServiceHandler(grpcServer, &SampleHandler{
		SampleService: service,
	})
	if err != nil {
		common.Fatal(err)
	}
}

func RegisterSalesHandler(grpcServer server.Server) {
	service := coresamplesservice.GetCurrentSalesSvc()
	err := pb.RegisterSalesServiceHandler(grpcServer, &SalesHandler{
		SalesService: service,
	})
	if err != nil {
		common.Fatal(err)
	}
}

func RegisterUserHandler(grpcServer server.Server) {
	service := coresamplesservice.GetCurrentUserService()
	err := pb.RegisterUserServiceHandler(grpcServer, &UserHandler{
		UserService: service,
	})
	if err != nil {
		common.Fatal(err)
	}
}

func RegisterInternalUserHandler(grpcServer server.Server) {
	service := coresamplesservice.GetCurrentUserService()
	err := pb.RegisterInternalUserServiceHandler(grpcServer, &InternalUserHandler{
		UserService: service,
	})
	if err != nil {
		common.Fatal(err)
	}
}

func RegisterCustomerHandler(grpcServer server.Server) {
	service := coresamplesservice.GetCurrentCustomerService()
	err := pb.RegisterCustomerServiceHandler(grpcServer, &CustomerHandler{
		CustomerService: service,
	})
	if err != nil {
		common.Fatal(err)
	}
}

//func registerMembershipRPCHandler(grpcServer server.Server, serviceshipService coresamplesservice.IServiceshipService) {
//	err := pb.RegisterServiceshipServiceHandler(grpcServer, &ServiceshipHTTPHandler{
//		service: serviceshipService,
//	})
//	if err != nil {
//		common.Fatal(err)
//	}
//}

func registerMembershipHTTPHandler(router *gin.RouterGroup, membershipService coresamplesservice.IServiceshipService) {
	membershipHTTPHandler = NewServiceshipHTTPHandler(membershipService)
	router.GET("/allow/:clinic_name", membershipHTTPHandler.SubscriptionAllowed)
	router.POST("/", membershipHTTPHandler.CreateServiceship)
	router.GET("/", membershipHTTPHandler.GetServiceshipsByType)
	router.GET("/:id", membershipHTTPHandler.GetServiceshipByID)
	router.POST("/subscription", membershipHTTPHandler.Subscribe)
	router.GET("/subscription", membershipHTTPHandler.GetAccountSubscriptions)
	router.POST("/subscription/update", membershipHTTPHandler.UpdateSubscription)
	router.POST("/subscription/pause", membershipHTTPHandler.PauseAutoRenew)
	router.POST("/subscription/resume", membershipHTTPHandler.ResumeAutoRenew)
	router.GET("/subscription/charging", membershipHTTPHandler.GetChargingSubscription)
	router.GET("/subscription/transaction", membershipHTTPHandler.GetSubscriptionTransactions)
	router.POST("/billingplan/", membershipHTTPHandler.CreateBillingPlanSet)
	router.POST("/billingplan/update", membershipHTTPHandler.AddBillingPlan)
	router.GET("/billingplan", membershipHTTPHandler.GetLatestBillingPlanSet)
	router.GET("/paymentmethod", membershipHTTPHandler.GetPaymentMethods)
	router.POST("/paymentmethod", membershipHTTPHandler.CreatePaymentMethod)
	router.DELETE("/paymentmethod", membershipHTTPHandler.DeletePaymentMethod)
}

func registerTestHTTPHandler(router *gin.RouterGroup, testService coresamplesservice.ITestService) {
	testHTTPHandler := NewTestHTTPHandler(testService)

	// Define HTTP routes for the Test service
	router.GET("/test", testHTTPHandler.GetTest)
	router.GET("/test/fields", testHTTPHandler.GetTestField)
	router.POST("/test", testHTTPHandler.CreateTest)
	router.GET("/test/codes", testHTTPHandler.GetTestIDsFromTestCodes)
	router.GET("/test/duplicates", testHTTPHandler.GetDuplicateAssayGroupTest)
}

func registerPatientHTTPHandler(router *gin.RouterGroup, patientService coresamplesservice.IPatientService) {
	patientHTTPHandler := NewPatientHTTPHandler(patientService)

	router.POST("/guest_login", patientHTTPHandler.PatientGuestLogIn)
}

func registerOrderHTTPHandler(router *gin.RouterGroup, orderService coresamplesservice.IOrderService) {
	orderHTTPHandler := NewOrderHTTPHandler(orderService)
	router.POST("/set_kit_status", orderHTTPHandler.UpdateOrderKitStatus)
}
