package main

import (
	"context"
	"coresamples/common"
	"coresamples/external"
	"coresamples/handler"
	"coresamples/middleware"
	"coresamples/processor"
	"coresamples/publisher"
	"coresamples/service"
	"coresamples/subscriber"
	"coresamples/util"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/hibiken/asynq"

	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"github.com/asim/go-micro/plugins/server/grpc/v4"
	httpServer "github.com/asim/go-micro/plugins/server/http/v4"
	limiter "github.com/asim/go-micro/plugins/wrapper/ratelimiter/uber/v4"
	"github.com/getsentry/sentry-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	capi "github.com/hashicorp/consul/api"
	"github.com/opentracing/opentracing-go"
	"github.com/segmentio/kafka-go"
	"go-micro.dev/v4"
	"go-micro.dev/v4/cache"
	"go-micro.dev/v4/registry"
	"go-micro.dev/v4/server"
)

const (
	version     = "latest"
	address     = "0.0.0.0:8084"
	httpAddress = "0.0.0.0:8083"
	QPS         = 100 // QPS for rate limiter
)

func main() {
	var err error
	var secret string
	common.Infof("init env")
	common.InitEnv()

	common.Infof("init consul")
	// Config center
	var consulClient *capi.Client
	var localConsulClient *capi.Client
	if common.Env.RunEnv == common.DevEnv {
		consulClient = common.GetConsulApiClient("http", "localhost", 8500, common.Env.ConsulToken, "")
	} else if common.Env.RunEnv == common.AksProductionEnv || common.Env.RunEnv == common.AksStagingEnv {
		consulClient = common.GetConsulApiClient("http", "consul.vibrant-wellness.com", 0, common.Env.ConsulToken, "")
		if common.Env.RunEnv == common.AksProductionEnv {
			localConsulClient = common.GetConsulApiClient("http", "192.168.60.9", 8500, common.Env.LocalConsulToken, "")
		}
	} else {
		consulClient = common.GetConsulApiClient("http", common.Env.ConsulConfigAddr, 8500, common.Env.ConsulToken, "")
	}
	// additional env variables stored in consul
	common.InitEnvFromConsul(consulClient, common.Env.ConsulPrefix, "env")
	// Init logger, this has to be done after Env is initialized completed from env and consul so that we know the log level
	common.InitZapLogger(common.Env.LogLevel)

	// fetch info from config center
	// DB config
	var mysqlInfo *common.MySQLConfig
	if common.Env.RunEnv != common.AksStagingEnv {
		mysqlInfo = common.GetMySqlConfigFromConsul(consulClient, common.Env.ConsulPrefix, "mysql")
	} else {
		mysqlInfo = common.GetMySqlConfigFromConsul(consulClient, common.Env.ConsulPrefix, "mysqlStaging")
	}

	// endpoints
	common.Infof("init external services and endpoints")
	if common.Env.RunEnv == common.AksStagingEnv || common.Env.RunEnv == common.AksProductionEnv {
		common.InitExternalServices(consulClient, common.Env.ConsulConfigAddr, common.Env.ConsulToken)
	}

	if common.Env.RunEnv == common.AksProductionEnv {
		common.InitEndpointsFromConsul(consulClient, common.Env.ConsulPrefix, "endpoints")
	} else {
		common.InitEndpointsFromConsul(consulClient, common.Env.ConsulPrefix, "endpointsStaging")
	}

	// secrets
	common.InitSecretsFromConsul(consulClient, common.Env.ConsulPrefix, "secrets")

	// external services
	if common.Env.RunEnv == common.AksProductionEnv {
		secret = common.Secrets.Secret
	} else {
		secret = common.Secrets.SecretStaging
	}

	external.InitAccountingService(secret)
	external.InitChargingService(secret)
	external.InitCustomerService(common.Env.ConsulToken, common.Env.ConsulConfigAddr)

	external.InitOrderService(secret)
	allowedServices := common.GetAllowedServicesFromConsul(consulClient, common.Env.ConsulPrefix, "allowedServices")
	allowedClinics := common.GetAllowedClinicsFromConsul(consulClient, common.Env.ConsulPrefix, "allowedClinics")
	// Init Redis Cache
	redisConfig := common.GetRedisConfigFromConsul(consulClient, common.Env.ConsulPrefix, "redisSentinel")
	redisClient := common.GetRedisReadWriteClient(redisConfig)
	asynqClient := common.GetAsynqClient(redisConfig)
	asynqServer := common.GetAsyncServer(redisConfig)
	rs := common.GetRedisSync(redisClient.RedisWriteClient)

	defer func() {
		if err := redisClient.Close(); err != nil {
			common.Error(err)
		}
	}()

	if common.Env.RunEnv == common.DevEnv ||
		common.Env.RunEnv == common.AksProductionEnv {
		common.GetLocalKafkaConfigsFromConsul(consulClient, common.Env.ConsulPrefix, "kafkaLocal")
	} else {
		common.GetLocalKafkaConfigsFromConsul(consulClient, common.Env.ConsulPrefix, "kafkaStaging")
	}

	//Registry center
	var consulRegistry registry.Registry
	if common.Env.RunEnv == common.DevEnv {
		consulRegistry = consul.NewRegistry(
			consul.Config(&capi.Config{
				Address: "127.0.0.1:8500",
				Token:   common.Env.ConsulToken,
			}))
	} else if common.Env.RunEnv == common.AksProductionEnv || common.Env.RunEnv == common.AksStagingEnv {
		consulRegistry = consul.NewRegistry(
			consul.Config(&capi.Config{
				Scheme:  "http",
				Address: "consul.vibrant-wellness.com",
				Token:   common.Env.ConsulToken,
			}),
		)
	} else {
		consulRegistry = consul.NewRegistry(
			consul.Config(&capi.Config{
				Address: common.Env.ConsulConfigAddr + ":8500",
				Token:   common.Env.ConsulToken,
			}))
	}

	//Tracer
	var tracer opentracing.Tracer
	var ioCloser io.Closer
	if common.Env.RunEnv == common.DevEnv {
		tracer, ioCloser, err = common.NewTracer("coresamplesv2", "127.0.0.1:6831")
	} else {
		tracer, ioCloser, err = common.NewTracer("coresamplesv2", "192.168.60.9:6831")
	}
	if err != nil {
		common.Fatal(err)
	}
	defer func() {
		if err := ioCloser.Close(); err != nil {
			common.Error(err)
		}
	}()
	opentracing.SetGlobalTracer(tracer)

	// Init DB
	common.Infof("init database")
	if common.Env.RunEnv == common.AksProductionEnv {
		rootCertPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile("/etc/ssl/certs/DBDigiCertGlobalRootCA.crt.pem")
		if err != nil {
			common.Fatal(err)
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			common.Fatal(fmt.Errorf("failed to append PEM"))
		}
		err = mysql.RegisterTLSConfig("custom", &tls.Config{
			RootCAs:    rootCertPool,
			ServerName: "lisportalprod2.mysql.database.azure.com",
		})
		if err != nil {
			common.Fatal(err)
		}
	}
	var dataSource string
	if common.Env.RunEnv == common.AksProductionEnv {
		dataSource = mysqlInfo.User + ":" + mysqlInfo.Pwd + "@tcp(" + mysqlInfo.Host + ":" + strconv.Itoa(mysqlInfo.Port) + ")/" + mysqlInfo.Database + "?parseTime=true&tls=custom"
	} else {
		dataSource = mysqlInfo.User + ":" + mysqlInfo.Pwd + "@tcp(" + mysqlInfo.Host + ":" + strconv.Itoa(mysqlInfo.Port) + ")/" + mysqlInfo.Database + "?parseTime=true"
	}

	dbClient, err := util.EntOpen("mysql", dataSource)
	if err != nil {
		common.Fatalf("failed opening connection to MySQL", err)
	}
	defer func() {
		if err := dbClient.Close(); err != nil {
			common.Error(err)
		}
	}()

	externDbClient, err := util.EntOpen("mysql", mysqlInfo.ExternDataSource)
	if err != nil {
		common.Fatalf("failed opening external connection to MySQL", err)
	}
	defer func() {
		if err := externDbClient.Close(); err != nil {
			common.Error(err)
		}
	}()

	if common.Env.RunEnv == common.DevEnv || common.Env.RunEnv == common.StagingEnv {
		err = dbClient.Schema.Create(context.Background())
		if err != nil {
			common.Fatal(err)
		}
	}

	grpcServer := grpc.NewServer(
		server.Name(common.Env.ServiceName),
		server.Version(version),
		server.Address(address),
		server.Registry(consulRegistry),
	)

	//var kbroker broker.Broker
	common.Infof("setup kafka: %v", common.LocalKafkaConfigs.Address)

	var localDialer *kafka.Dialer
	localDialer = &kafka.Dialer{
		DualStack: true,
	}

	publisher.InitPublisher(context.Background(), nil, common.LocalKafkaConfigs.Address)
	defer publisher.GetPublisher().GetWriter().Close()

	common.Infof("create http server")
	srv := httpServer.NewServer(
		server.Version(version),
		server.Address(httpAddress),
	)
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(middleware.AuthMiddleware(true))
	router.Use(middleware.Ginzap("", true))

	var v1 *gin.RouterGroup
	//if common.Env.RunEnv == common.AksProductionEnv {
	//	v1 = router.Group("/v1/lis/coresamples-v2-service-http")
	//} else if common.Env.RunEnv == common.AksStagingEnv || common.Env.RunEnv == common.StagingEnv {
	//	v1 = router.Group("/v1/lis/coresamples-v2-service-http-staging")
	//} else {
	//	v1 = router.Group("/")
	//}
	v1 = router.Group("/")
	err = sentry.Init(sentry.ClientOptions{
		// Either set your DSN here or set the SENTRY_DSN environment variable.
		Dsn: "https://fde61d937fcf4562bc7de1e439be1f21@sentry1.vibrant-america.com/37",
		// Either set environment and release here or set the SENTRY_ENVIRONMENT
		// and SENTRY_RELEASE environment variables.
		// Environment: "Production",
		// Release:     "my-project-name@1.0.0",
		// Enable printing of SDK debug messages.
		// Useful when getting started or trying to figure something out.
		Debug:       true,
		Environment: common.Env.RunEnv,
		// Set TracesSampleRate to 1.0 to capture 100%
		// of transactions for performance monitoring.
		// We recommend adjusting this value in production,
		TracesSampleRate: 1,
		AttachStacktrace: true,
	})

	if err != nil {
		common.Fatalf("sentry.Init", err)
	}
	defer sentry.Flush(2 * time.Second)

	//explicitly init all services
	service.GetTubeService(dbClient, redisClient)
	service.GetRBACService(dbClient, redisClient, "mysql", dataSource)
	service.GetServiceshipSvc(dbClient, redisClient, allowedClinics)
	service.GetTestService(dbClient, redisClient)
	service.GetPatientSvc(externDbClient, redisClient, secret)
	service.GetOrderService(dbClient, externDbClient, redisClient)
	service.GetSampleService(dbClient, redisClient, asynqClient)
	service.GetSalesSvc(dbClient, redisClient)
	service.GetUserService(dbClient, redisClient)
	service.GetCustomerService(dbClient, redisClient)

	// Register handler
	common.Infof("register handlers")
	handler.RegisterTubeHandler(grpcServer)
	handler.RegisterRBACHandler(grpcServer, allowedServices)
	handler.RegisterHealthHandler(grpcServer)
	handler.RegisterServiceshipHandler(v1.Group("/serviceship"))
	handler.RegisterTestHandler(v1.Group("/test"), grpcServer)
	handler.RegisterOrderHandler(v1.Group("/order"), grpcServer)
	//TODO: replace this with local ent client once we import the table
	handler.RegisterPatientHandler(v1.Group("/patient"))
	handler.RegisterSampleHandler(grpcServer)
	handler.RegisterSalesHandler(grpcServer)
	handler.RegisterInternalUserHandler(grpcServer)
	handler.RegisterCustomerHandler(grpcServer)

	testAuth := router.Group("/healthcheck")
	testAuth.GET("/", func(c *gin.Context) {
		if common.Env.ServiceDeregistered {
			c.JSON(500, gin.H{
				"message": "service restart",
			})
		} else {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		}
	})
	hd := srv.NewHandler(router)
	if err := srv.Handle(hd); err != nil {
		common.Fatal(err)
	}

	common.Infof("create micro service")
	httpSvc := micro.NewService(
		micro.Server(srv),
		micro.WrapHandler(middleware.NewTracingHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapHandler(limiter.NewHandlerWrapper(QPS)),
		micro.WrapHandler(middleware.LogHandler),
		micro.WrapHandler(middleware.RecoveryWithZap),
		micro.Cache(cache.NewCache()),
	)
	httpSvc.Init()
	go func() {
		if err := httpSvc.Run(); err != nil {
			common.Fatal(err)
		}
	}()
	// Create service
	svc := micro.NewService(
		micro.Server(grpcServer),
		micro.WrapHandler(middleware.NewTracingHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapHandler(limiter.NewHandlerWrapper(QPS)),
		micro.WrapHandler(middleware.LogHandler),
		micro.WrapHandler(middleware.RecoveryWithZap),
		micro.Cache(cache.NewCache()),
	)

	// Asynq handler
	common.Infof("registering Asynq handler")
	mux := asynq.NewServeMux()
	processor.RegisterSampleProcessor(mux, dbClient, redisClient, rs)
	processor.RegisterUserProcessor(mux, dbClient, redisClient)
	processor.RegisterCDCUpdateProcessor(mux, dbClient, redisClient, rs)

	go func() {
		if err = asynqServer.Run(mux); err != nil {
			common.Fatalf("Could not run Asynq server", err)
		}
	}()
	defer func(Client *asynq.Client) {
		err = Client.Close()
		if err != nil {
			common.Error(err)
		}
	}(asynqClient)

	common.Infof("run subscriber")
	eventHandler := subscriber.NewEventHandler(dbClient, asynqClient, context.Background(), common.LocalKafkaConfigs.Address, localDialer)
	eventHandler.Run()

	if common.Env.RunEnv == common.AksProductionEnv || common.Env.RunEnv == common.DevEnv {
		cdcUpdateHandler := subscriber.NewCoreCDCUpdatesHandler(common.LocalKafkaConfigs.Address, localDialer, asynqClient)
		cdcUpdateHandler.Run()
	}

	// start health check
	//handler.HandleHTTPHealthCheck()

	deregDone := make(chan bool, 1)
	if common.Env.RunEnv == common.AksProductionEnv {
		ip := common.Env.PodIP
		name := "go.micro.lis.service.coresamples.v2.prod"
		common.Infof("prod address: %s:%d", ip, 8084)
		id := "go.micro.lis.service.coresamples.v2.prod" + "-" + uuid.NewString()
		registration := &capi.AgentServiceRegistration{
			ID:      id,
			Name:    name,
			Port:    8084,
			Address: ip,
			Check: &capi.AgentServiceCheck{
				TTL:                            "15s",
				DeregisterCriticalServiceAfter: "10m",
				CheckID:                        id,
			},
		}

		go func() {
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL)
			common.Infof("sigs registered")
			err = localConsulClient.Agent().ServiceRegister(registration)
			if err != nil {
				common.Fatal(err)
			}
			common.Infof("registered service " + id)
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					err = localConsulClient.Agent().UpdateTTL(id, "periodic pass", "passing")
					if err != nil {
						common.Error(err)
						if strings.Contains(err.Error(), "Unknown check ID") {
							common.Infof("Service deregistered, restart now")
							common.Env.ServiceDeregistered = true
							deregDone <- true
							return
						}
					}
				case <-sigs:
					common.Infof("Stop updating TTL, deregistering service " + id)
					err := localConsulClient.Agent().ServiceDeregister(id)
					if err != nil {
						common.Error(err)
					} else {
						common.Infof("Deregistered service " + id)
					}
					// mark dregistration done
					deregDone <- true
					return
				}
			}
		}()
	} else {
		// For other environment, no need to wait for explicit deregistration
		deregDone <- true
	}

	svc.Init()

	common.Infof("init grpc server")
	err = grpcServer.Init(server.Wait(nil))

	if err != nil {
		common.Fatal(err)
	}

	common.Infof("run service")
	// Run service
	if err := svc.Run(); err != nil {
		common.Fatal(err)
	}
	// wait for service deregistration, if there is a need, in case program terminates before deregistration done
	common.Infof("wait for dereg done")
	<-deregDone
	common.Infof("dereg done")
}
