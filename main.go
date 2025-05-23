package main

import (
	"context"
	"coresamples/common"
	"coresamples/external"
	"coresamples/handler"
	"coresamples/middleware"
	"coresamples/processor"
	"coresamples/service"
	"coresamples/util"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"go.uber.org/zap"

	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"github.com/asim/go-micro/plugins/server/grpc/v4"
	httpServer "github.com/asim/go-micro/plugins/server/http/v4"
	limiter "github.com/asim/go-micro/plugins/wrapper/ratelimiter/uber/v4"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	capi "github.com/hashicorp/consul/api"
	"github.com/opentracing/opentracing-go"
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
	} else if common.Env.RunEnv == common.DevDockerComposeEnv {
		// Use localhost for Consul in dev_docker_compose mode
		consulClient = common.GetConsulApiClient("http", "localhost", 8500, common.Env.ConsulToken, "")
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
	// var err error // err is already declared in main
	mysqlInfo, err = common.GetCurrentMySQLConfig(consulClient, common.Env.ConsulPrefix)
	if err != nil {
		common.Fatal(fmt.Errorf("critical: failed to initialize MySQL config: %w", err))
	}
	if mysqlInfo == nil { // Should not happen if GetCurrentMySQLConfig guarantees non-nil or error
		common.Fatal(fmt.Errorf("critical: MySQL config is nil after attempting to load"))
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

	// Override JWTSecret if in dev_docker_compose mode and ENV VAR is set
	if common.Env.RunEnv == common.DevDockerComposeEnv {
		jwtSecretEnv := os.Getenv("JWT_SECRET_VAL") // Corrected ENV var name to match docker-compose.yaml
		if jwtSecretEnv != "" {
			if common.Secrets == nil { // Should be initialized by InitSecretsFromConsul, even if to an empty struct on error
				common.Warn("common.Secrets was nil before JWT_SECRET_VAL ENV override; initializing.")
				common.Secrets = &common.SecretsConfig{}
			}
			common.Secrets.JWTSecret = jwtSecretEnv
			common.InfoFields("Loaded JWT_SECRET_VAL from environment variable for dev_docker_compose mode, overriding any Consul value.",
				zap.String("source", "environment_variable"))
		} else {
			common.Warn("dev_docker_compose mode: JWT_SECRET_VAL environment variable is not set. Using value from Consul or default if any.")
		}
	}

	// Critical validation for JWTSecret after all loading attempts
	if common.Secrets == nil || common.Secrets.JWTSecret == "" {
		// This check is important. If JWTSecret is empty, auth will fail or be insecure.
		common.Fatal(fmt.Errorf("CRITICAL: common.Secrets.JWTSecret is empty after all configuration attempts. " +
			"Ensure it's set in Consul (as jwt_secret) or via JWT_SECRET_VAL environment variable for dev_docker_compose mode."))
	}

	// external services
	// The rest of the code that uses common.Secrets.Secret / SecretStaging for 'secret' variable:
	// var secret string // Assuming 'secret' was declared earlier as per original main.go
	if common.Env.RunEnv == common.AksProductionEnv {
		if common.Secrets != nil && common.Secrets.Secret != "" {
			secret = common.Secrets.Secret
		} else {
			common.Fatal(fmt.Errorf("common.Secrets is nil or common.Secrets.Secret is empty, cannot retrieve 'Secret' for AksProductionEnv"))
		}
	} else {
		if common.Secrets != nil && common.Secrets.SecretStaging != "" {
			secret = common.Secrets.SecretStaging
		} else {
			common.Fatal(fmt.Errorf("common.Secrets is nil or common.Secrets.SecretStaging is empty, cannot retrieve 'SecretStaging' for non-AksProductionEnv"))
		}
	}
	// TODO for future: Consider if 'secret' and 'secretStaging' also need ENV var overrides for dev_docker_compose
	// if external.InitXService(secret) calls are essential for basic startup in dev_docker_compose.
	// For now, this task focuses on JWT_SECRET.

	// EXTERNAL SERVICE INITIALIZATION MODIFICATION
	if common.Env.RunEnv == common.DevDockerComposeEnv {
		common.Info("Initializing external services for dev_docker_compose (errors will be logged as warnings).")
		func() {
			defer func() {
				if r := recover(); r != nil {
					common.Warn("Panic recovered during InitAccountingService/InitChargingService/InitOrderService for dev_docker_compose")
				}
			}()
			common.Info("Attempting to initialize Accounting, Charging, Order services for dev_docker_compose...")
			external.InitAccountingService(secret)
			external.InitChargingService(secret)
			external.InitOrderService(secret)
			common.Info("Finished attempting to initialize Accounting, Charging, Order services for dev_docker_compose.")
		}()
		func() {
			defer func() {
				if r := recover(); r != nil {
					common.Warn("Panic recovered during InitCustomerService for dev_docker_compose")
				}
			}()
			common.Info("Attempting to initialize CustomerService for dev_docker_compose. Note: it may try to use Consul for discovery unless internally adapted.")
			external.InitCustomerService(common.Env.ConsulToken, common.Env.ConsulConfigAddr)
			common.Info("Finished attempting to initialize CustomerService for dev_docker_compose.")
		}()
	} else {
		common.Info("Initializing external services for standard mode.")
		external.InitAccountingService(secret)
		external.InitChargingService(secret)
		external.InitCustomerService(common.Env.ConsulToken, common.Env.ConsulConfigAddr)
		external.InitOrderService(secret)
	}

	allowedServices := common.GetAllowedServicesFromConsul(consulClient, common.Env.ConsulPrefix, "allowedServices")
	allowedClinics := common.GetAllowedClinicsFromConsul(consulClient, common.Env.ConsulPrefix, "allowedClinics")
	// Init Redis Cache
	var redisConfig *common.RedisConfig
	// var err error; // Assuming err is already declared
	redisConfig, err = common.GetCurrentRedisConfig(consulClient, common.Env.ConsulPrefix)
	if err != nil {
		common.Fatal(fmt.Errorf("critical: failed to initialize Redis config: %w", err))
	}
	if redisConfig == nil { // Should not happen if GetCurrentRedisConfig guarantees non-nil or error
		common.Fatal(fmt.Errorf("critical: Redis config is nil after attempting to load"))
	}
	redisClient := common.GetRedisReadWriteClient(redisConfig)
	asynqClient := common.GetAsynqClient(redisConfig)
	asynqServer := common.GetAsyncServer(redisConfig)
	rs := common.GetRedisSync(redisClient.RedisWriteClient)

	defer func() {
		if err := redisClient.Close(); err != nil {
			common.Error(err)
		}
	}()

	// Kafka configuration
	/*
	var kafkaAvailable bool = false
	if common.Env.RunEnv == common.DevDockerComposeEnv {
		kafkaBrokerEnv := os.Getenv("KAFKA_BROKERS")
		if kafkaBrokerEnv != "" {
			if common.LocalKafkaConfigs == nil {
				common.LocalKafkaConfigs = &common.KafkaConfiguration{}
			}
			common.LocalKafkaConfigs.Address = []string{kafkaBrokerEnv}
			common.InfoFields("Configuring Kafka from KAFKA_BROKERS ENV VAR for dev_docker_compose.", zap.String("brokers", kafkaBrokerEnv))
		} else {
			common.Warn("dev_docker_compose mode: KAFKA_BROKERS environment variable is not set. Kafka functionality will be limited.")
			if common.LocalKafkaConfigs == nil {
				common.LocalKafkaConfigs = &common.KafkaConfiguration{}
			}
			common.LocalKafkaConfigs.Address = []string{} // Ensure it's empty
		}
	} else {
		// Existing logic for GetLocalKafkaConfigsFromConsul
		if common.Env.RunEnv == common.AksProductionEnv {
			common.GetLocalKafkaConfigsFromConsul(consulClient, common.Env.ConsulPrefix, "kafkaLocal")
		} else {
			common.GetLocalKafkaConfigsFromConsul(consulClient, common.Env.ConsulPrefix, "kafkaStaging")
		}
	}

	// After this, check common.LocalKafkaConfigs and set kafkaAvailable
	if common.LocalKafkaConfigs != nil && len(common.LocalKafkaConfigs.Address) > 0 && common.LocalKafkaConfigs.Address[0] != "" {
		kafkaAvailable = true
		common.InfoFields("Kafka configuration loaded from Consul.", zap.Strings("brokers", common.LocalKafkaConfigs.Address))
	}
	*/

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
	} else if common.Env.RunEnv == common.DevDockerComposeEnv {
		consulRegistry = consul.NewRegistry(
			consul.Config(&capi.Config{
				Address: "consul:8500",
				Token:   common.Env.ConsulToken,
			}))
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
	// var err error // Assuming err is already declared in main()

	var jaegerAddr string
	// common.Env.ServiceName should be correctly set by InitEnv based on Env.RunEnv
	// If common.Env.ServiceName is not suitable, use a fixed string like "coresamplesv2" as before.
	serviceNameForTracer := common.Env.ServiceName 
	if serviceNameForTracer == "" { // Fallback if service name isn't set for some reason
		serviceNameForTracer = "coresamplesv2"
		common.Warn("common.Env.ServiceName is empty, using default 'coresamplesv2' for Jaeger.")
	}


	if common.Env.RunEnv == common.DevDockerComposeEnv {
		jaegerHost := os.Getenv("JAEGER_AGENT_HOST")
		jaegerPort := os.Getenv("JAEGER_AGENT_PORT")

		if jaegerHost == "" {
			jaegerHost = "localhost" // Default if not set, expecting 'jaeger' service from compose
			common.Warn("JAEGER_AGENT_HOST environment variable not set for dev_docker_compose, defaulting to 'localhost'. This should ideally be 'jaeger' (the service name in docker-compose.yaml).")
		}
		if jaegerPort == "" {
			jaegerPort = "6831" // Default Jaeger agent UDP port
			common.Warn("JAEGER_AGENT_PORT environment variable not set for dev_docker_compose, defaulting to '6831'.")
		}
		jaegerAddr = fmt.Sprintf("%s:%s", jaegerHost, jaegerPort)
		common.InfoFields("Initializing Jaeger tracer for dev_docker_compose mode.",
			zap.String("service_name", serviceNameForTracer),
			zap.String("resolved_jaeger_address", jaegerAddr))
	} else if common.Env.RunEnv == common.DevEnv { // Assuming common.DevEnv is a defined constant
		jaegerAddr = "127.0.0.1:6831" // Original DevEnv logic
		common.InfoFields("Initializing Jaeger tracer for dev mode.",
			zap.String("service_name", serviceNameForTracer),
			zap.String("jaeger_address", jaegerAddr))
	} else { // Production, Staging, etc.
		// This was the original "else" logic for non-DevEnv
		jaegerAddr = "192.168.60.9:6831" 
		common.InfoFields("Initializing Jaeger tracer for production/staging mode.",
			zap.String("service_name", serviceNameForTracer),
			zap.String("jaeger_address", jaegerAddr))
	}

	tracer, ioCloser, err = common.NewTracer(serviceNameForTracer, jaegerAddr)
	if err != nil {
		// Use common.Fatalf if it's the project's standard for fatal errors with formatted messages
		common.Fatal(fmt.Errorf("failed to create Jaeger tracer: %w", err))
	}
	defer func() {
		if err := ioCloser.Close(); err != nil {
			common.Error(err) // Assuming common.Error logs the error
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
	// Check mysqlInfo for nil before accessing its fields, although previous checks should prevent this.
	if mysqlInfo == nil {
		common.Fatal(fmt.Errorf("MySQL configuration is nil before creating dataSource")) // Should have been caught earlier
	}

	connectionParams := "?parseTime=true"
	if common.Env.RunEnv == common.AksProductionEnv && !mysqlInfo.DisableTLS {
		// Ensure TLS setup for DB is still attempted for production if not explicitly disabled
		// (The original code for RegisterTLSConfig for "custom" should remain for this path)
		// This example assumes RegisterTLSConfig has been called if needed.
		connectionParams += "&tls=custom"
		common.Info("Attempting to use custom TLS for MySQL connection in Production.")
	} else if mysqlInfo.DisableTLS {
		// For dev_docker_compose with MYSQL_DISABLE_TLS=true, or other envs where DisableTLS is true in config
		common.Info("MySQL TLS explicitly disabled via configuration.")
		// No specific TLS parameter, or ensure 'tls=false' or similar if required by driver for explicit disable
	}
	// For other cases (e.g., dev/staging not using TLS from original logic), no explicit TLS param was added.

	dataSource = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s%s",
		mysqlInfo.User,
		mysqlInfo.Pwd,
		mysqlInfo.Host,
		mysqlInfo.Port,
		mysqlInfo.Database,
		connectionParams,
	)

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
	// Kafka Publisher Initialization
	/*
	var localDialer *kafka.Dialer // Declare localDialer here so it's in scope for subscribers if Kafka is available
	if kafkaAvailable {
		common.InfoFields("Kafka is configured. Initializing publisher and subscribers.", zap.Strings("Brokers", common.LocalKafkaConfigs.Address))
		localDialer = &kafka.Dialer{DualStack: true} // Assuming this is generally safe

		publisher.InitPublisher(context.Background(), nil, common.LocalKafkaConfigs.Address)
		defer func() {
			if pub := publisher.GetPublisher(); pub != nil {
				if writer := pub.GetWriter(); writer != nil {
					if err := writer.Close(); err != nil {
						common.Error(fmt.Errorf("error closing kafka writer: %w", err))
					}
				} else {
					common.Debugf("Kafka writer was nil, nothing to close.")
				}
			} else {
				common.Debugf("Kafka publisher was nil, nothing to close.")
			}
		}()
	} else {
		common.Warn("Kafka is not available or not configured. Publisher and subscribers will not be initialized.")
	}
	*/

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
	// err = sentry.Init(sentry.ClientOptions{
	// 	// Either set your DSN here or set the SENTRY_DSN environment variable.
	// 	Dsn: "https://fde61d937fcf4562bc7de1e439be1f21@sentry1.vibrant-america.com/37",
	// 	// Either set environment and release here or set the SENTRY_ENVIRONMENT
	// 	// and SENTRY_RELEASE environment variables.
	// 	// Environment: "Production",
	// 	// Release:     "my-project-name@1.0.0",
	// 	// Enable printing of SDK debug messages.
	// 	// Useful when getting started or trying to figure something out.
	// 	Debug:       true,
	// 	Environment: common.Env.RunEnv,
	// 	// Set TracesSampleRate to 1.0 to capture 100%
	// 	// of transactions for performance monitoring.
	// 	// We recommend adjusting this value in production,
	// 	TracesSampleRate: 1,
	// 	AttachStacktrace: true,
	// })

	// if err != nil {
	// 	common.Fatalf("sentry.Init", err)
	// }
	// defer sentry.Flush(2 * time.Second)

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

	// Kafka Subscriber Initialization
	/*
	if kafkaAvailable {
		common.Infof("Starting Kafka subscribers...")
		eventHandler := subscriber.NewEventHandler(dbClient, asynqClient, context.Background(), common.LocalKafkaConfigs.Address, localDialer)
		if common.Env.RunEnv == common.DevDockerComposeEnv {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						common.Warn("Panic recovered in eventHandler.Run() for dev_docker_compose")
					}
				}()
				common.Info("eventHandler.Run() starting in goroutine for dev_docker_compose")
				eventHandler.Run()
			}()
		} else {
			eventHandler.Run()
		}

		// Handle cdcUpdateHandler similarly
		// Original condition: common.Env.RunEnv == common.AksProductionEnv || common.Env.RunEnv == common.DevEnv
		// For dev_docker_compose, it should also run if kafka is available.
		originalCdcEnvMatch := common.Env.RunEnv == common.AksProductionEnv || common.Env.RunEnv == common.DevEnv
		if originalCdcEnvMatch || (common.Env.RunEnv == common.DevDockerComposeEnv && kafkaAvailable) {
			cdcUpdateHandler := subscriber.NewCoreCDCUpdatesHandler(common.LocalKafkaConfigs.Address, localDialer, asynqClient)
			if common.Env.RunEnv == common.DevDockerComposeEnv {
				go func() {
					defer func() {
						if r := recover(); r != nil {
							common.Warn("Panic recovered in cdcUpdateHandler.Run() for dev_docker_compose")
						}
					}()
					common.Info("cdcUpdateHandler.Run() starting in goroutine for dev_docker_compose")
					cdcUpdateHandler.Run()
				}()
			} else {
				cdcUpdateHandler.Run()
			}
		}
	}
	*/

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
