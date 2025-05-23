package common

import (
	"errors"
	capi "github.com/hashicorp/consul/api"
	"math/rand"
	"strconv"
	"time"
)

const (
	TestGroupServiceName = "lis-order-dev-default"
)

var extSvc *ExternalServices

type ExternalServices struct {
	client    *capi.Client
	consulURL string
	token     string
}

func GetExternalServices() *ExternalServices {
	return extSvc
}

func InitExternalServices(client *capi.Client, consulURL string, token string) {
	if extSvc != nil {
		return
	}
	extSvc = &ExternalServices{
		consulURL: consulURL,
		client:    client,
		token:     token,
	}
}

func DiscoverServiceEndpoint(client *capi.Client, name string) (string, error) {
	svc, _, err := client.Catalog().Service(name, "", nil)
	if err != nil {
		return "", err
	}

	if svc == nil || len(svc) == 0 {
		return "", errors.New("Empty service: " + name)
	}
	rand.Seed(time.Now().Unix())
	pod := svc[rand.Intn(len(svc))]
	return pod.ServiceAddress + ":" + strconv.Itoa(pod.ServicePort), nil
}

func DiscoverServices(client *capi.Client, name string) ([]*capi.CatalogService, error) {
	svc, _, err := client.Catalog().Service(name, "", nil)
	if err != nil {
		return nil, err
	}

	if svc == nil || len(svc) == 0 {
		return nil, errors.New("Empty service: " + name)
	}
	return svc, nil
}
