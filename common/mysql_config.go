package common

import (
	"encoding/json"
	capi "github.com/hashicorp/consul/api"
)

type MySQLConfig struct {
	Host             string `json:"host,omitempty"`
	User             string `json:"user,omitempty"`
	Pwd              string `json:"password,omitempty"`
	Database         string `json:"database,omitempty"`
	Port             int    `json:"port,omitempty"`
	ExternDataSource string `json:"extern_data_source,omitempty"`
}

func GetMySqlConfigFromConsul(client *capi.Client, prefix string, key string) *MySQLConfig {
	mysqlConfig := &MySQLConfig{}
	val, _, err := client.KV().Get(prefix+"/"+key, nil)
	if err != nil {
		Fatal(err)
	}
	err = json.Unmarshal(val.Value, mysqlConfig)
	if err != nil {
		Fatal(err)
	}
	return mysqlConfig
}
