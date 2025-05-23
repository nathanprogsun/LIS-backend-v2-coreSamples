package handler

import (
	"context"
	pb "coresamples/proto"
	"coresamples/service"
)

type ContactHandler struct {
	ContactService service.IContactService
}

func (h *ContactHandler) UpdateContact(ctx context.Context, req *pb.UpdateContactRequest, resp *pb.Contact) error {
	return nil
}

func (h *ContactHandler) CreateContact(ctx context.Context, req *pb.CreateContactRequest, resp *pb.Contact) error {
	return nil
}

func (h *ContactHandler) UpdateGroupContact(ctx context.Context, req *pb.UpdateGroupContactRequest, resp *pb.CreateOrUpdateGroupContactResponse) error {
	return nil
}

func (h *ContactHandler) CreateOrUpdateGroupContact(ctx context.Context, req *pb.CreateOrUpdateGroupContactRequest, resp *pb.CreateOrUpdateGroupContactResponse) error {
	return nil
}

// TODO: change to customer contact on clinics instead
// func (h *ContactHandler) ShowCustomerContact(ctx context.Context, req *pb.ShowCustomerContactRequest, resp *pb.ShowCustomerContactResponse) error {
// 	return nil
// }

func (h *ContactHandler) ShowClinicContact(ctx context.Context, req *pb.ShowClinicContactRequest, resp *pb.ShowCustomerContactResponse) error {
	return nil
}

func (h *ContactHandler) DeleteContact(ctx context.Context, req *pb.DeleteContactRequest, resp *pb.DeleteContactResponse) error {
	return nil
}
