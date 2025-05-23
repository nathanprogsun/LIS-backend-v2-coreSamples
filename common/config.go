package common

import (
	"encoding/json"
	capi "github.com/hashicorp/consul/api"
	"strconv"
)

func GetConsulApiClient(scheme string, host string, port int64, token string, CAFile string) *capi.Client {
	if port != 0 {
		host = host + ":" + strconv.FormatInt(port, 10)
	}
	cfg := capi.Config{
		Scheme:  scheme,
		Address: host,
		Token:   token,
	}
	if CAFile != "" {
		cfg.TLSConfig = capi.TLSConfig{
			CAFile:             CAFile,
			InsecureSkipVerify: true,
		}
	}

	capiClient, err := capi.NewClient(&cfg)
	if err != nil {
		Fatal(err)
	}

	return capiClient
}

//func GetConsulConfig(host string, port int64, prefix string) (config.Config, error) {
//	consulSource := consul.NewSource(
//		consul.WithAddress(host+":"+strconv.FormatInt(port, 10)),
//		consul.WithPrefix(prefix),
//		consul.StripPrefix(true),
//		consul.WithToken(ConsulToken),
//	)
//
//	consulConfig, err := config.NewConfig()
//	if err != nil {
//		return consulConfig, err
//	}
//
//	err = consulConfig.Load(consulSource)
//	return consulConfig, err
//}

func getConsulArrayConfig(client *capi.Client, prefix string, key string) []string {
	var ArrayConfig []string
	val, _, err := client.KV().Get(prefix+"/"+key, nil)
	if err != nil {
		Fatal(err)
	}
	err = json.Unmarshal(val.Value, &ArrayConfig)
	if err != nil {
		Fatal(err)
	}
	return ArrayConfig
}
