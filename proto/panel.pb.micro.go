// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/panel.proto

package coresamples_service

import (
	fmt "fmt"
	proto "google.golang.org/protobuf/proto"
	_ "google.golang.org/protobuf/types/known/durationpb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
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

// Api Endpoints for PanelService service

func NewPanelServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for PanelService service

type PanelService interface {
	CreatePanel(ctx context.Context, in *CreatePanelRequest, opts ...client.CallOption) (*Panel, error)
}

type panelService struct {
	c    client.Client
	name string
}

func NewPanelService(name string, c client.Client) PanelService {
	return &panelService{
		c:    c,
		name: name,
	}
}

func (c *panelService) CreatePanel(ctx context.Context, in *CreatePanelRequest, opts ...client.CallOption) (*Panel, error) {
	req := c.c.NewRequest(c.name, "PanelService.CreatePanel", in)
	out := new(Panel)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for PanelService service

type PanelServiceHandler interface {
	CreatePanel(context.Context, *CreatePanelRequest, *Panel) error
}

func RegisterPanelServiceHandler(s server.Server, hdlr PanelServiceHandler, opts ...server.HandlerOption) error {
	type panelService interface {
		CreatePanel(ctx context.Context, in *CreatePanelRequest, out *Panel) error
	}
	type PanelService struct {
		panelService
	}
	h := &panelServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&PanelService{h}, opts...))
}

type panelServiceHandler struct {
	PanelServiceHandler
}

func (h *panelServiceHandler) CreatePanel(ctx context.Context, in *CreatePanelRequest, out *Panel) error {
	return h.PanelServiceHandler.CreatePanel(ctx, in, out)
}
