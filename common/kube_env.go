package common

import (
	"encoding/json"
	capi "github.com/hashicorp/consul/api"
	"log"
	"os"
)

const (
	DevConsulPrefix = "micro/config/lis/"
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
	Env.ConsulConfigAddr = os.Getenv("CONSUL_ADDR")
	Env.ConsulToken = os.Getenv("CONSUL_TOKEN")
	Env.ProdConsulToken = os.Getenv("CONSUL_TOKEN_PROD")
	Env.LocalConsulToken = os.Getenv("CONSUL_TOKEN_LOCAL")
	Env.PodIP = os.Getenv("POD_IP")
	Env.ServiceDeregistered = false
	if Env.RunEnv == "" {
		Env.RunEnv = DevEnv
	}
	if Env.ConsulPrefix == "" {
		Env.ConsulPrefix = DevConsulPrefix
	}
	if Env.RunEnv == StagingEnv || Env.RunEnv == AksStagingEnv {
		Env.ServiceName = "go.micro.lis.service.coresamples.v2.staging"
	} else if Env.RunEnv == AksProductionEnv {
		Env.ServiceName = "go.micro.lis.service.coresamples.v2"
	} else {
		Env.ServiceName = "go.micro.lis.service.coresamples.v2.dev"
	}
	Env.AsynqQueueName = "coresamplesv2_" + Env.RunEnv
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
