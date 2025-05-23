package common

import (
	"encoding/json"
	capi "github.com/hashicorp/consul/api"
)

type SecretsConfig struct {
	Secret        string `json:"secret,omitempty"`
	SecretStaging string `json:"secret_staging"`
	OrderToken    string `json:"orderToken,omitempty"`
	JWTSecret     string
}

var Secrets = &SecretsConfig{}

func InitSecretsFromConsul(client *capi.Client, prefix string, key string) {
	val, _, err := client.KV().Get(prefix+"/"+key, nil)
	if err != nil {
		Fatal(err)
	}
	err = json.Unmarshal(val.Value, Secrets)
	if err != nil {
		Fatal(err)
	}
	if Env.RunEnv == AksProductionEnv {
		Secrets.JWTSecret = Secrets.Secret
	} else {
		Secrets.JWTSecret = Secrets.SecretStaging
	}
}
