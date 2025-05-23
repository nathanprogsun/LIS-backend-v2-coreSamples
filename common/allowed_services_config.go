package common

import (
	"encoding/json"
	capi "github.com/hashicorp/consul/api"
)

func GetAllowedServicesFromConsul(client *capi.Client, prefix string, key string) []string {
	var allowedServices []string
	val, _, err := client.KV().Get(prefix+"/"+key, nil)
	if err != nil {
		Fatal(err)
	}
	err = json.Unmarshal(val.Value, &allowedServices)
	if err != nil {
		Fatal(err)
	}
	return allowedServices
}
