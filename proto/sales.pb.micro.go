// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/sales.proto

package coresamples_service

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	math "math"
)

import (
	context "context"
	api "go-micro.dev/v4/api"
	client "go-micro.dev/v4/client"
	server "go-micro.dev/v4/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for SalesService service

func NewSalesServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for SalesService service

type SalesService interface {
	GetSalesByTerritory(ctx context.Context, in *Territory, opts ...client.CallOption) (*SalesName, error)
}

type salesService struct {
	c    client.Client
	name string
}

func NewSalesService(name string, c client.Client) SalesService {
	return &salesService{
		c:    c,
		name: name,
	}
}

func (c *salesService) GetSalesByTerritory(ctx context.Context, in *Territory, opts ...client.CallOption) (*SalesName, error) {
	req := c.c.NewRequest(c.name, "SalesService.GetSalesByTerritory", in)
	out := new(SalesName)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for SalesService service

type SalesServiceHandler interface {
	GetSalesByTerritory(context.Context, *Territory, *SalesName) error
}

func RegisterSalesServiceHandler(s server.Server, hdlr SalesServiceHandler, opts ...server.HandlerOption) error {
	type salesService interface {
		GetSalesByTerritory(ctx context.Context, in *Territory, out *SalesName) error
	}
	type SalesService struct {
		salesService
	}
	h := &salesServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&SalesService{h}, opts...))
}

type salesServiceHandler struct {
	SalesServiceHandler
}

func (h *salesServiceHandler) GetSalesByTerritory(ctx context.Context, in *Territory, out *SalesName) error {
	return h.SalesServiceHandler.GetSalesByTerritory(ctx, in, out)
}
