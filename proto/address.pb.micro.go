// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/address.proto

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

// Api Endpoints for AddressService service

func NewAddressServiceEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for AddressService service

type AddressService interface {
	GetAddress(ctx context.Context, in *GetAddressRequest, opts ...client.CallOption) (*Address, error)
	UpdateAddress(ctx context.Context, in *UpdateAddressRequest, opts ...client.CallOption) (*Address, error)
	CreateAddress(ctx context.Context, in *CreateAddressRequest, opts ...client.CallOption) (*Address, error)
	UpdateGroupAddress(ctx context.Context, in *UpdateGroupAddressRequest, opts ...client.CallOption) (*CreateOrUpdateGroupAddressResponse, error)
	CreateOrUpdateGroupAddress(ctx context.Context, in *CreateOrUpdateGroupAddressRequest, opts ...client.CallOption) (*CreateOrUpdateGroupAddressResponse, error)
	ShowCustomerAddress(ctx context.Context, in *ShowCustomerAddressRequest, opts ...client.CallOption) (*ShowCustomerAddressResponse, error)
	ShowClinicAddress(ctx context.Context, in *ShowClinicAddressRequest, opts ...client.CallOption) (*ShowCustomerAddressResponse, error)
}

type addressService struct {
	c    client.Client
	name string
}

func NewAddressService(name string, c client.Client) AddressService {
	return &addressService{
		c:    c,
		name: name,
	}
}

func (c *addressService) GetAddress(ctx context.Context, in *GetAddressRequest, opts ...client.CallOption) (*Address, error) {
	req := c.c.NewRequest(c.name, "AddressService.GetAddress", in)
	out := new(Address)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addressService) UpdateAddress(ctx context.Context, in *UpdateAddressRequest, opts ...client.CallOption) (*Address, error) {
	req := c.c.NewRequest(c.name, "AddressService.UpdateAddress", in)
	out := new(Address)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addressService) CreateAddress(ctx context.Context, in *CreateAddressRequest, opts ...client.CallOption) (*Address, error) {
	req := c.c.NewRequest(c.name, "AddressService.CreateAddress", in)
	out := new(Address)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addressService) UpdateGroupAddress(ctx context.Context, in *UpdateGroupAddressRequest, opts ...client.CallOption) (*CreateOrUpdateGroupAddressResponse, error) {
	req := c.c.NewRequest(c.name, "AddressService.UpdateGroupAddress", in)
	out := new(CreateOrUpdateGroupAddressResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addressService) CreateOrUpdateGroupAddress(ctx context.Context, in *CreateOrUpdateGroupAddressRequest, opts ...client.CallOption) (*CreateOrUpdateGroupAddressResponse, error) {
	req := c.c.NewRequest(c.name, "AddressService.CreateOrUpdateGroupAddress", in)
	out := new(CreateOrUpdateGroupAddressResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addressService) ShowCustomerAddress(ctx context.Context, in *ShowCustomerAddressRequest, opts ...client.CallOption) (*ShowCustomerAddressResponse, error) {
	req := c.c.NewRequest(c.name, "AddressService.ShowCustomerAddress", in)
	out := new(ShowCustomerAddressResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addressService) ShowClinicAddress(ctx context.Context, in *ShowClinicAddressRequest, opts ...client.CallOption) (*ShowCustomerAddressResponse, error) {
	req := c.c.NewRequest(c.name, "AddressService.ShowClinicAddress", in)
	out := new(ShowCustomerAddressResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AddressService service

type AddressServiceHandler interface {
	GetAddress(context.Context, *GetAddressRequest, *Address) error
	UpdateAddress(context.Context, *UpdateAddressRequest, *Address) error
	CreateAddress(context.Context, *CreateAddressRequest, *Address) error
	UpdateGroupAddress(context.Context, *UpdateGroupAddressRequest, *CreateOrUpdateGroupAddressResponse) error
	CreateOrUpdateGroupAddress(context.Context, *CreateOrUpdateGroupAddressRequest, *CreateOrUpdateGroupAddressResponse) error
	ShowCustomerAddress(context.Context, *ShowCustomerAddressRequest, *ShowCustomerAddressResponse) error
	ShowClinicAddress(context.Context, *ShowClinicAddressRequest, *ShowCustomerAddressResponse) error
}

func RegisterAddressServiceHandler(s server.Server, hdlr AddressServiceHandler, opts ...server.HandlerOption) error {
	type addressService interface {
		GetAddress(ctx context.Context, in *GetAddressRequest, out *Address) error
		UpdateAddress(ctx context.Context, in *UpdateAddressRequest, out *Address) error
		CreateAddress(ctx context.Context, in *CreateAddressRequest, out *Address) error
		UpdateGroupAddress(ctx context.Context, in *UpdateGroupAddressRequest, out *CreateOrUpdateGroupAddressResponse) error
		CreateOrUpdateGroupAddress(ctx context.Context, in *CreateOrUpdateGroupAddressRequest, out *CreateOrUpdateGroupAddressResponse) error
		ShowCustomerAddress(ctx context.Context, in *ShowCustomerAddressRequest, out *ShowCustomerAddressResponse) error
		ShowClinicAddress(ctx context.Context, in *ShowClinicAddressRequest, out *ShowCustomerAddressResponse) error
	}
	type AddressService struct {
		addressService
	}
	h := &addressServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&AddressService{h}, opts...))
}

type addressServiceHandler struct {
	AddressServiceHandler
}

func (h *addressServiceHandler) GetAddress(ctx context.Context, in *GetAddressRequest, out *Address) error {
	return h.AddressServiceHandler.GetAddress(ctx, in, out)
}

func (h *addressServiceHandler) UpdateAddress(ctx context.Context, in *UpdateAddressRequest, out *Address) error {
	return h.AddressServiceHandler.UpdateAddress(ctx, in, out)
}

func (h *addressServiceHandler) CreateAddress(ctx context.Context, in *CreateAddressRequest, out *Address) error {
	return h.AddressServiceHandler.CreateAddress(ctx, in, out)
}

func (h *addressServiceHandler) UpdateGroupAddress(ctx context.Context, in *UpdateGroupAddressRequest, out *CreateOrUpdateGroupAddressResponse) error {
	return h.AddressServiceHandler.UpdateGroupAddress(ctx, in, out)
}

func (h *addressServiceHandler) CreateOrUpdateGroupAddress(ctx context.Context, in *CreateOrUpdateGroupAddressRequest, out *CreateOrUpdateGroupAddressResponse) error {
	return h.AddressServiceHandler.CreateOrUpdateGroupAddress(ctx, in, out)
}

func (h *addressServiceHandler) ShowCustomerAddress(ctx context.Context, in *ShowCustomerAddressRequest, out *ShowCustomerAddressResponse) error {
	return h.AddressServiceHandler.ShowCustomerAddress(ctx, in, out)
}

func (h *addressServiceHandler) ShowClinicAddress(ctx context.Context, in *ShowClinicAddressRequest, out *ShowCustomerAddressResponse) error {
	return h.AddressServiceHandler.ShowClinicAddress(ctx, in, out)
}
