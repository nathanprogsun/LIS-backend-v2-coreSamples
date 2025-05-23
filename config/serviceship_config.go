package config

import (
	"coresamples/common"
	_ "embed"
	"gopkg.in/yaml.v3"
)

//go:embed serviceship.yaml
var configb []byte

type ServiceshipConfig struct {
	Membership map[string]MembershipConfig `yaml:"membership,omitempty"`
}

type MembershipConfig struct {
	Name            string                 `yaml:"name,omitempty"`
	Description     string                 `yaml:"description,omitempty"`
	AccountType     string                 `yaml:"account_type,omitempty"`
	AllowedServices []string               `yaml:"allowed_services,omitempty"`
	Bonus           map[string]interface{} `yaml:"bonus"`
}

var svcConfig *ServiceshipConfig

func GetServiceshipConfig() *ServiceshipConfig {
	if svcConfig == nil {
		svcConfig = &ServiceshipConfig{}
		if err := yaml.Unmarshal(configb, svcConfig); err != nil {
			common.Fatal(err)
		}
	}
	return svcConfig
}
