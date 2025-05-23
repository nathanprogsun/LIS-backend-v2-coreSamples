package common

import (
	"encoding/json"
	capi "github.com/hashicorp/consul/api"
)

type EndpointsConfig struct {
	GetTestGroup        string `json:"getTestGroup,omitempty"`
	GetTestGroupMapping string `json:"getTestGroupMapping,omitempty"`
	Accounting          string `json:"accounting,omitempty"`
	Charging            string `json:"charging,omitempty"`
	Order               string `json:"order,omitempty"`
}

var EndpointsInfo = &EndpointsConfig{}

func InitEndpointsFromConsul(client *capi.Client, prefix string, key string) {
	val, _, err := client.KV().Get(prefix+"/"+key, nil)
	if err != nil {
		Error(err)
	}
	err = json.Unmarshal(val.Value, EndpointsInfo)
	if err != nil {
		Error(err)
	}
}
