package common

import (
	"encoding/json"
	capi "github.com/hashicorp/consul/api"
	"log"
	"os"
)

const (
	DevConsulPrefix     = "micro/config/lis/"
	DevDockerComposeEnv = "dev_docker_compose"
)

var Env struct {
	RunEnv              string // Staging or Production
	ConsulConfigAddr    string
	ConsulPrefix        string
	ConsulToken         string
	ProdConsulToken     string
	LocalConsulToken    string
	PodIP               string
	ServiceName         string
	AsynqQueueName      string
	ServiceDeregistered bool
	LogLevel            string `json:"log_level,omitempty"`
	OrderServiceName    string `json:"order_service_name,omitempty"`
	DryRun              bool   `json:"dry_run,omitempty"`
	ConsulClient        *capi.Client
}

func InitEnv() {
	Env.RunEnv = os.Getenv("CORESAMPLES_ENV")
	Env.ConsulPrefix = os.Getenv("CONSUL_PREFIX")
	// Standard Consul address environment variable
	Env.ConsulConfigAddr = os.Getenv("CONSUL_ADDR")
	Env.ConsulToken = os.Getenv("CONSUL_TOKEN")
	Env.ProdConsulToken = os.Getenv("CONSUL_TOKEN_PROD")
	Env.LocalConsulToken = os.Getenv("CONSUL_TOKEN_LOCAL")
	Env.PodIP = os.Getenv("POD_IP")
	Env.ServiceDeregistered = false

	// Specific handling for Docker Compose environment
	if Env.RunEnv == DevDockerComposeEnv {
		// In Docker Compose, CONSUL_HTTP_ADDR is typically used to define the accessible address
		// for the Consul service container.
		consulHttpAddr := os.Getenv("CONSUL_HTTP_ADDR")
		if consulHttpAddr != "" {
			Env.ConsulConfigAddr = consulHttpAddr
		}
		// Set default log level from ENV if provided, otherwise InitEnvFromConsul might override it
		// or it might rely on Consul. For docker-compose, direct ENV is better.
		logLevel := os.Getenv("LOG_LEVEL")
		if logLevel != "" {
			Env.LogLevel = logLevel
		} else {
			Env.LogLevel = "debug" // Default for dev_docker_compose if not set
		}
	} else if Env.RunEnv == "" {
		Env.RunEnv = DevEnv // Default to DevEnv if CORESAMPLES_ENV is not set at all
	}

	if Env.ConsulPrefix == "" {
		Env.ConsulPrefix = DevConsulPrefix
	}

	// Service Name Logic (ensure DevDockerComposeEnv results in a .dev-like name or is handled)
	if Env.RunEnv == StagingEnv || Env.RunEnv == AksStagingEnv {
		Env.ServiceName = "go.micro.lis.service.coresamples.v2.staging"
	} else if Env.RunEnv == AksProductionEnv {
		Env.ServiceName = "go.micro.lis.service.coresamples.v2"
	} else if Env.RunEnv == DevDockerComposeEnv {
		Env.ServiceName = "go.micro.lis.service.coresamples.v2.devcompose" // Or use the .dev default
	} else { // Defaults to DevEnv or any other non-specified CORESAMPLES_ENV value
		Env.ServiceName = "go.micro.lis.service.coresamples.v2.dev"
	}
	Env.AsynqQueueName = "coresamplesv2_" + Env.RunEnv

	// If LogLevel wasn't set by DevDockerComposeEnv specific logic or by InitEnvFromConsul later,
	// ensure a default if it's still empty.
	// However, InitZapLogger in main.go will use Env.LogLevel, so it should be set before that.
	// The InitEnvFromConsul might overwrite Env.LogLevel.
	// For dev_docker_compose, we want the ENV LOG_LEVEL to be authoritative.
	// This might mean reading LOG_LEVEL env var *after* InitEnvFromConsul in main.go for this specific mode.
	// For now, the above setting in DevDockerComposeEnv block is the primary attempt.
}

func InitEnvFromConsul(client *capi.Client, prefix string, key string) {
	val, _, err := client.KV().Get(prefix+"/"+key, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(val.Value, &Env)
	if err != nil {
		log.Fatal(err)
	}
	Env.ConsulClient = client
}
