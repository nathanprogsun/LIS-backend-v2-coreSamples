package handler

import (
	"context"
	pb "coresamples/proto"
	"coresamples/service"
)

type InternalUserHandler struct {
	UserService service.IUserService
}

func (h *InternalUserHandler) CreateInternalUser(ctx context.Context, req *pb.CreateInternalUserRequest, resp *pb.InternalUser) error {
	return nil
}

func (h *InternalUserHandler) GetSalesClinics(ctx context.Context, req *pb.ListSalesNameRequest, resp *pb.ListCustomerPracticeResponse) error {
	return nil
}

func (h *InternalUserHandler) GetInternalUser(ctx context.Context, req *pb.GetInternalUserRequest, resp *pb.GetInternalUserResponse) error {
	// Extract parameters from the request
	role := req.Role
	roleIDs := req.RoleIds
	usernames := req.Usernames

	// Call the service method
	result, err := h.UserService.GetInternalUser(ctx, role, roleIDs, usernames)
	if err != nil {
		return err
	}

	// Copy the response
	resp.Response = result.Response
	return nil
}

func (h *InternalUserHandler) TransferSalesCustomer(ctx context.Context, req *pb.TransferSalesCustomerRequest, resp *pb.TransferSalesCustomerResponse) error {
	// Extract parameters from the request
	fromSalesID := req.FromSalesId
	toSalesID := req.ToSalesId
	customerID := req.CustomerId

	// Call the service method
	result, err := h.UserService.TransferSalesCustomer(ctx, fromSalesID, toSalesID, customerID)
	if err != nil {
		return err
	}

	resp.Status = result.Status

	return nil
}

func (h *InternalUserHandler) GetLowerLevelInternalUsers(ctx context.Context, req *pb.GetLowerLevelInternalUsersRequest, resp *pb.GetLowerLevelInternalUsersResponse) error {
	return nil
}

func (h *InternalUserHandler) SetLowerLevelInternalUsers(ctx context.Context, req *pb.SetLowerLevelInternalUsersRequest, resp *pb.SetLowerLevelInternalUsersResponse) error {
	return nil
}

func (h *InternalUserHandler) CreateSampleNavigatorNote(ctx context.Context, req *pb.CreateSampleNavigatorNoteRequest, resp *pb.CreateSampleNavigatorNoteResponse) error {
	return nil
}

func (h *InternalUserHandler) ModifySampleNavigatorNote(ctx context.Context, req *pb.ModifySampleNavigatorNoteRequest, resp *pb.ModifySampleNavigatorNoteResponse) error {
	return nil
}

func (h *InternalUserHandler) DeleteSampleNavigatorNote(ctx context.Context, req *pb.DeleteSampleNavigatorNoteRequest, resp *pb.DeleteSampleNavigatorNoteResponse) error {
	return nil
}

func (h *InternalUserHandler) GetInternalUserByid(ctx context.Context, req *pb.GetInternalUserByidRequest, resp *pb.GetInternalUserByidResponse) error {
	return nil
}

func (h *InternalUserHandler) CheckCustomerNavigator(ctx context.Context, req *pb.CheckNavigatorCustomerRequest, resp *pb.CheckNavigatorCustomerResponse) error {
	return nil
}
