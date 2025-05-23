package handler

import (
	"context"
	pb "coresamples/proto"
	"coresamples/service"
)

type AddressHandler struct {
	AddressService service.IAddressService
}

func (h *AddressHandler) GetAddress(ctx context.Context, req *pb.GetAddressRequest, resp *pb.Address) error {
	return nil
}

func (h *AddressHandler) UpdateAddress(ctx context.Context, req *pb.UpdateAddressRequest, resp *pb.Address) error {
	return nil
}

func (h *AddressHandler) CreateAddress(ctx context.Context, req *pb.CreateAddressRequest, resp *pb.Address) error {
	return nil
}

func (h *AddressHandler) UpdateGroupAddress(ctx context.Context, req *pb.UpdateGroupAddressRequest, resp *pb.CreateOrUpdateGroupAddressResponse) error {
	return nil
}

func (h *AddressHandler) CreateOrUpdateGroupAddress(ctx context.Context, req *pb.CreateOrUpdateGroupAddressRequest, resp *pb.CreateOrUpdateGroupAddressResponse) error {
	return nil
}

func (h *AddressHandler) ShowCustomerAddress(ctx context.Context, req *pb.ShowCustomerAddressRequest, resp *pb.ShowCustomerAddressResponse) error {
	return nil
}

func (h *AddressHandler) ShowClinicAddress(ctx context.Context, req *pb.ShowClinicAddressRequest, resp *pb.ShowCustomerAddressResponse) error {
	return nil
}
