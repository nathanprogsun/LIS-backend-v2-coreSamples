// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: proto/clinic.proto

package coresamples_service

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	ClinicService_CreateClinic_FullMethodName                   = "/coresamples_service.ClinicService/CreateClinic"
	ClinicService_GetClinic_FullMethodName                      = "/coresamples_service.ClinicService/GetClinic"
	ClinicService_ListClinic_FullMethodName                     = "/coresamples_service.ClinicService/ListClinic"
	ClinicService_ListClinicCustomers_FullMethodName            = "/coresamples_service.ClinicService/ListClinicCustomers"
	ClinicService_GetCustomerClinicNames_FullMethodName         = "/coresamples_service.ClinicService/GetCustomerClinicNames"
	ClinicService_GetClinicsClinicAccountIDs_FullMethodName     = "/coresamples_service.ClinicService/GetClinicsClinicAccountIDs"
	ClinicService_ListClinicCustomersByClinicID_FullMethodName  = "/coresamples_service.ClinicService/ListClinicCustomersByClinicID"
	ClinicService_GetClinicByID_FullMethodName                  = "/coresamples_service.ClinicService/GetClinicByID"
	ClinicService_UpdateClinicNPINumber_FullMethodName          = "/coresamples_service.ClinicService/UpdateClinicNPINumber"
	ClinicService_EditClinicName_FullMethodName                 = "/coresamples_service.ClinicService/EditClinicName"
	ClinicService_ModifyClinicNames_FullMethodName              = "/coresamples_service.ClinicService/modifyClinicNames"
	ClinicService_SearchClinicsByName_FullMethodName            = "/coresamples_service.ClinicService/SearchClinicsByName"
	ClinicService_SignUpClinicForExistingAccount_FullMethodName = "/coresamples_service.ClinicService/SignUpClinicForExistingAccount"
	ClinicService_CheckClinicAttributes_FullMethodName          = "/coresamples_service.ClinicService/CheckClinicAttributes"
	ClinicService_GetFirstCustomerOfClinic_FullMethodName       = "/coresamples_service.ClinicService/GetFirstCustomerOfClinic"
	ClinicService_FuzzySearchClinics_FullMethodName             = "/coresamples_service.ClinicService/FuzzySearchClinics"
	ClinicService_GetClinicAddress_FullMethodName               = "/coresamples_service.ClinicService/GetClinicAddress"
)

// ClinicServiceClient is the client API for ClinicService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClinicServiceClient interface {
	CreateClinic(ctx context.Context, in *CreateClinicRequest, opts ...grpc.CallOption) (*Clinic, error)
	GetClinic(ctx context.Context, in *GetClinicRequest, opts ...grpc.CallOption) (*GetClinicResponse, error)
	ListClinic(ctx context.Context, in *ListClinicRequest, opts ...grpc.CallOption) (*GetClinicResponse, error)
	ListClinicCustomers(ctx context.Context, in *ListClinicCustomerRequest, opts ...grpc.CallOption) (*ListClinicCustomerResponse, error)
	GetCustomerClinicNames(ctx context.Context, in *GetCustomerClinicNamesRequest, opts ...grpc.CallOption) (*GetCustomerClinicNamesResponse, error)
	GetClinicsClinicAccountIDs(ctx context.Context, in *GetClinicsClinicAccountIDsRequest, opts ...grpc.CallOption) (*GetClinicResponse, error)
	ListClinicCustomersByClinicID(ctx context.Context, in *ListClinicCustomersByClinicIDRequest, opts ...grpc.CallOption) (*ListClinicCustomerByIDResponse, error)
	// Version 0.7.3.7
	GetClinicByID(ctx context.Context, in *ClinicID, opts ...grpc.CallOption) (*FullClinic, error)
	// Version 0.7.3.9
	UpdateClinicNPINumber(ctx context.Context, in *UpdateClinicNPINumberRequest, opts ...grpc.CallOption) (*GetClinicResponse, error)
	EditClinicName(ctx context.Context, in *EditClinicNameRequest, opts ...grpc.CallOption) (*GetClinicResponse, error)
	ModifyClinicNames(ctx context.Context, in *ModifyClinicNamesRequest, opts ...grpc.CallOption) (*GetClinicResponse, error)
	SearchClinicsByName(ctx context.Context, in *SearchClinicsNameRequest, opts ...grpc.CallOption) (*SearchClinicsInfoResponse, error)
	SignUpClinicForExistingAccount(ctx context.Context, in *SignUpClinicRequest, opts ...grpc.CallOption) (*SignUpClinicResponse, error)
	CheckClinicAttributes(ctx context.Context, in *CheckClinicAttributesRequest, opts ...grpc.CallOption) (*CheckClinicAttributesResponse, error)
	// VP-4965
	GetFirstCustomerOfClinic(ctx context.Context, in *GetFirstCustomerOfClinicRequest, opts ...grpc.CallOption) (*GetFirstCustomerOfClinicResponse, error)
	FuzzySearchClinics(ctx context.Context, in *FuzzySearchClinicsRequest, opts ...grpc.CallOption) (*SearchClinicsInfoResponse, error)
	GetClinicAddress(ctx context.Context, in *ClinicID, opts ...grpc.CallOption) (*GetClinicAddressResponse, error)
}

type clinicServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewClinicServiceClient(cc grpc.ClientConnInterface) ClinicServiceClient {
	return &clinicServiceClient{cc}
}

func (c *clinicServiceClient) CreateClinic(ctx context.Context, in *CreateClinicRequest, opts ...grpc.CallOption) (*Clinic, error) {
	out := new(Clinic)
	err := c.cc.Invoke(ctx, ClinicService_CreateClinic_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clinicServiceClient) GetClinic(ctx context.Context, in *GetClinicRequest, opts ...grpc.CallOption) (*GetClinicResponse, error) {
	out := new(GetClinicResponse)
	err := c.cc.Invoke(ctx, ClinicService_GetClinic_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clinicServiceClient) ListClinic(ctx context.Context, in *ListClinicRequest, opts ...grpc.CallOption) (*GetClinicResponse, error) {
	out := new(GetClinicResponse)
	err := c.cc.Invoke(ctx, ClinicService_ListClinic_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clinicServiceClient) ListClinicCustomers(ctx context.Context, in *ListClinicCustomerRequest, opts ...grpc.CallOption) (*ListClinicCustomerResponse, error) {
	out := new(ListClinicCustomerResponse)
	err := c.cc.Invoke(ctx, ClinicService_ListClinicCustomers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clinicServiceClient) GetCustomerClinicNames(ctx context.Context, in *GetCustomerClinicNamesRequest, opts ...grpc.CallOption) (*GetCustomerClinicNamesResponse, error) {
	out := new(GetCustomerClinicNamesResponse)
	err := c.cc.Invoke(ctx, ClinicService_GetCustomerClinicNames_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clinicServiceClient) GetClinicsClinicAccountIDs(ctx context.Context, in *GetClinicsClinicAccountIDsRequest, opts ...grpc.CallOption) (*GetClinicResponse, error) {
	out := new(GetClinicResponse)
	err := c.cc.Invoke(ctx, ClinicService_GetClinicsClinicAccountIDs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clinicServiceClient) ListClinicCustomersByClinicID(ctx context.Context, in *ListClinicCustomersByClinicIDRequest, opts ...grpc.CallOption) (*ListClinicCustomerByIDResponse, error) {
	out := new(ListClinicCustomerByIDResponse)
	err := c.cc.Invoke(ctx, ClinicService_ListClinicCustomersByClinicID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clinicServiceClient) GetClinicByID(ctx context.Context, in *ClinicID, opts ...grpc.CallOption) (*FullClinic, error) {
	out := new(FullClinic)
	err := c.cc.Invoke(ctx, ClinicService_GetClinicByID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clinicServiceClient) UpdateClinicNPINumber(ctx context.Context, in *UpdateClinicNPINumberRequest, opts ...grpc.CallOption) (*GetClinicResponse, error) {
	out := new(GetClinicResponse)
	err := c.cc.Invoke(ctx, ClinicService_UpdateClinicNPINumber_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clinicServiceClient) EditClinicName(ctx context.Context, in *EditClinicNameRequest, opts ...grpc.CallOption) (*GetClinicResponse, error) {
	out := new(GetClinicResponse)
	err := c.cc.Invoke(ctx, ClinicService_EditClinicName_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clinicServiceClient) ModifyClinicNames(ctx context.Context, in *ModifyClinicNamesRequest, opts ...grpc.CallOption) (*GetClinicResponse, error) {
	out := new(GetClinicResponse)
	err := c.cc.Invoke(ctx, ClinicService_ModifyClinicNames_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clinicServiceClient) SearchClinicsByName(ctx context.Context, in *SearchClinicsNameRequest, opts ...grpc.CallOption) (*SearchClinicsInfoResponse, error) {
	out := new(SearchClinicsInfoResponse)
	err := c.cc.Invoke(ctx, ClinicService_SearchClinicsByName_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clinicServiceClient) SignUpClinicForExistingAccount(ctx context.Context, in *SignUpClinicRequest, opts ...grpc.CallOption) (*SignUpClinicResponse, error) {
	out := new(SignUpClinicResponse)
	err := c.cc.Invoke(ctx, ClinicService_SignUpClinicForExistingAccount_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clinicServiceClient) CheckClinicAttributes(ctx context.Context, in *CheckClinicAttributesRequest, opts ...grpc.CallOption) (*CheckClinicAttributesResponse, error) {
	out := new(CheckClinicAttributesResponse)
	err := c.cc.Invoke(ctx, ClinicService_CheckClinicAttributes_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clinicServiceClient) GetFirstCustomerOfClinic(ctx context.Context, in *GetFirstCustomerOfClinicRequest, opts ...grpc.CallOption) (*GetFirstCustomerOfClinicResponse, error) {
	out := new(GetFirstCustomerOfClinicResponse)
	err := c.cc.Invoke(ctx, ClinicService_GetFirstCustomerOfClinic_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clinicServiceClient) FuzzySearchClinics(ctx context.Context, in *FuzzySearchClinicsRequest, opts ...grpc.CallOption) (*SearchClinicsInfoResponse, error) {
	out := new(SearchClinicsInfoResponse)
	err := c.cc.Invoke(ctx, ClinicService_FuzzySearchClinics_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *clinicServiceClient) GetClinicAddress(ctx context.Context, in *ClinicID, opts ...grpc.CallOption) (*GetClinicAddressResponse, error) {
	out := new(GetClinicAddressResponse)
	err := c.cc.Invoke(ctx, ClinicService_GetClinicAddress_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClinicServiceServer is the server API for ClinicService service.
// All implementations must embed UnimplementedClinicServiceServer
// for forward compatibility
type ClinicServiceServer interface {
	CreateClinic(context.Context, *CreateClinicRequest) (*Clinic, error)
	GetClinic(context.Context, *GetClinicRequest) (*GetClinicResponse, error)
	ListClinic(context.Context, *ListClinicRequest) (*GetClinicResponse, error)
	ListClinicCustomers(context.Context, *ListClinicCustomerRequest) (*ListClinicCustomerResponse, error)
	GetCustomerClinicNames(context.Context, *GetCustomerClinicNamesRequest) (*GetCustomerClinicNamesResponse, error)
	GetClinicsClinicAccountIDs(context.Context, *GetClinicsClinicAccountIDsRequest) (*GetClinicResponse, error)
	ListClinicCustomersByClinicID(context.Context, *ListClinicCustomersByClinicIDRequest) (*ListClinicCustomerByIDResponse, error)
	// Version 0.7.3.7
	GetClinicByID(context.Context, *ClinicID) (*FullClinic, error)
	// Version 0.7.3.9
	UpdateClinicNPINumber(context.Context, *UpdateClinicNPINumberRequest) (*GetClinicResponse, error)
	EditClinicName(context.Context, *EditClinicNameRequest) (*GetClinicResponse, error)
	ModifyClinicNames(context.Context, *ModifyClinicNamesRequest) (*GetClinicResponse, error)
	SearchClinicsByName(context.Context, *SearchClinicsNameRequest) (*SearchClinicsInfoResponse, error)
	SignUpClinicForExistingAccount(context.Context, *SignUpClinicRequest) (*SignUpClinicResponse, error)
	CheckClinicAttributes(context.Context, *CheckClinicAttributesRequest) (*CheckClinicAttributesResponse, error)
	// VP-4965
	GetFirstCustomerOfClinic(context.Context, *GetFirstCustomerOfClinicRequest) (*GetFirstCustomerOfClinicResponse, error)
	FuzzySearchClinics(context.Context, *FuzzySearchClinicsRequest) (*SearchClinicsInfoResponse, error)
	GetClinicAddress(context.Context, *ClinicID) (*GetClinicAddressResponse, error)
	mustEmbedUnimplementedClinicServiceServer()
}

// UnimplementedClinicServiceServer must be embedded to have forward compatible implementations.
type UnimplementedClinicServiceServer struct {
}

func (UnimplementedClinicServiceServer) CreateClinic(context.Context, *CreateClinicRequest) (*Clinic, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateClinic not implemented")
}
func (UnimplementedClinicServiceServer) GetClinic(context.Context, *GetClinicRequest) (*GetClinicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClinic not implemented")
}
func (UnimplementedClinicServiceServer) ListClinic(context.Context, *ListClinicRequest) (*GetClinicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListClinic not implemented")
}
func (UnimplementedClinicServiceServer) ListClinicCustomers(context.Context, *ListClinicCustomerRequest) (*ListClinicCustomerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListClinicCustomers not implemented")
}
func (UnimplementedClinicServiceServer) GetCustomerClinicNames(context.Context, *GetCustomerClinicNamesRequest) (*GetCustomerClinicNamesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCustomerClinicNames not implemented")
}
func (UnimplementedClinicServiceServer) GetClinicsClinicAccountIDs(context.Context, *GetClinicsClinicAccountIDsRequest) (*GetClinicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClinicsClinicAccountIDs not implemented")
}
func (UnimplementedClinicServiceServer) ListClinicCustomersByClinicID(context.Context, *ListClinicCustomersByClinicIDRequest) (*ListClinicCustomerByIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListClinicCustomersByClinicID not implemented")
}
func (UnimplementedClinicServiceServer) GetClinicByID(context.Context, *ClinicID) (*FullClinic, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClinicByID not implemented")
}
func (UnimplementedClinicServiceServer) UpdateClinicNPINumber(context.Context, *UpdateClinicNPINumberRequest) (*GetClinicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateClinicNPINumber not implemented")
}
func (UnimplementedClinicServiceServer) EditClinicName(context.Context, *EditClinicNameRequest) (*GetClinicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditClinicName not implemented")
}
func (UnimplementedClinicServiceServer) ModifyClinicNames(context.Context, *ModifyClinicNamesRequest) (*GetClinicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ModifyClinicNames not implemented")
}
func (UnimplementedClinicServiceServer) SearchClinicsByName(context.Context, *SearchClinicsNameRequest) (*SearchClinicsInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchClinicsByName not implemented")
}
func (UnimplementedClinicServiceServer) SignUpClinicForExistingAccount(context.Context, *SignUpClinicRequest) (*SignUpClinicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SignUpClinicForExistingAccount not implemented")
}
func (UnimplementedClinicServiceServer) CheckClinicAttributes(context.Context, *CheckClinicAttributesRequest) (*CheckClinicAttributesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckClinicAttributes not implemented")
}
func (UnimplementedClinicServiceServer) GetFirstCustomerOfClinic(context.Context, *GetFirstCustomerOfClinicRequest) (*GetFirstCustomerOfClinicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFirstCustomerOfClinic not implemented")
}
func (UnimplementedClinicServiceServer) FuzzySearchClinics(context.Context, *FuzzySearchClinicsRequest) (*SearchClinicsInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FuzzySearchClinics not implemented")
}
func (UnimplementedClinicServiceServer) GetClinicAddress(context.Context, *ClinicID) (*GetClinicAddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetClinicAddress not implemented")
}
func (UnimplementedClinicServiceServer) mustEmbedUnimplementedClinicServiceServer() {}

// UnsafeClinicServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClinicServiceServer will
// result in compilation errors.
type UnsafeClinicServiceServer interface {
	mustEmbedUnimplementedClinicServiceServer()
}

func RegisterClinicServiceServer(s grpc.ServiceRegistrar, srv ClinicServiceServer) {
	s.RegisterService(&ClinicService_ServiceDesc, srv)
}

func _ClinicService_CreateClinic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateClinicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClinicServiceServer).CreateClinic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClinicService_CreateClinic_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClinicServiceServer).CreateClinic(ctx, req.(*CreateClinicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClinicService_GetClinic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetClinicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClinicServiceServer).GetClinic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClinicService_GetClinic_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClinicServiceServer).GetClinic(ctx, req.(*GetClinicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClinicService_ListClinic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListClinicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClinicServiceServer).ListClinic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClinicService_ListClinic_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClinicServiceServer).ListClinic(ctx, req.(*ListClinicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClinicService_ListClinicCustomers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListClinicCustomerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClinicServiceServer).ListClinicCustomers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClinicService_ListClinicCustomers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClinicServiceServer).ListClinicCustomers(ctx, req.(*ListClinicCustomerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClinicService_GetCustomerClinicNames_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCustomerClinicNamesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClinicServiceServer).GetCustomerClinicNames(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClinicService_GetCustomerClinicNames_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClinicServiceServer).GetCustomerClinicNames(ctx, req.(*GetCustomerClinicNamesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClinicService_GetClinicsClinicAccountIDs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetClinicsClinicAccountIDsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClinicServiceServer).GetClinicsClinicAccountIDs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClinicService_GetClinicsClinicAccountIDs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClinicServiceServer).GetClinicsClinicAccountIDs(ctx, req.(*GetClinicsClinicAccountIDsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClinicService_ListClinicCustomersByClinicID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListClinicCustomersByClinicIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClinicServiceServer).ListClinicCustomersByClinicID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClinicService_ListClinicCustomersByClinicID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClinicServiceServer).ListClinicCustomersByClinicID(ctx, req.(*ListClinicCustomersByClinicIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClinicService_GetClinicByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClinicID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClinicServiceServer).GetClinicByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClinicService_GetClinicByID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClinicServiceServer).GetClinicByID(ctx, req.(*ClinicID))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClinicService_UpdateClinicNPINumber_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateClinicNPINumberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClinicServiceServer).UpdateClinicNPINumber(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClinicService_UpdateClinicNPINumber_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClinicServiceServer).UpdateClinicNPINumber(ctx, req.(*UpdateClinicNPINumberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClinicService_EditClinicName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EditClinicNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClinicServiceServer).EditClinicName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClinicService_EditClinicName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClinicServiceServer).EditClinicName(ctx, req.(*EditClinicNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClinicService_ModifyClinicNames_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ModifyClinicNamesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClinicServiceServer).ModifyClinicNames(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClinicService_ModifyClinicNames_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClinicServiceServer).ModifyClinicNames(ctx, req.(*ModifyClinicNamesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClinicService_SearchClinicsByName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchClinicsNameRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClinicServiceServer).SearchClinicsByName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClinicService_SearchClinicsByName_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClinicServiceServer).SearchClinicsByName(ctx, req.(*SearchClinicsNameRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClinicService_SignUpClinicForExistingAccount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SignUpClinicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClinicServiceServer).SignUpClinicForExistingAccount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClinicService_SignUpClinicForExistingAccount_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClinicServiceServer).SignUpClinicForExistingAccount(ctx, req.(*SignUpClinicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClinicService_CheckClinicAttributes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckClinicAttributesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClinicServiceServer).CheckClinicAttributes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClinicService_CheckClinicAttributes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClinicServiceServer).CheckClinicAttributes(ctx, req.(*CheckClinicAttributesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClinicService_GetFirstCustomerOfClinic_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFirstCustomerOfClinicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClinicServiceServer).GetFirstCustomerOfClinic(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClinicService_GetFirstCustomerOfClinic_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClinicServiceServer).GetFirstCustomerOfClinic(ctx, req.(*GetFirstCustomerOfClinicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClinicService_FuzzySearchClinics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FuzzySearchClinicsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClinicServiceServer).FuzzySearchClinics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClinicService_FuzzySearchClinics_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClinicServiceServer).FuzzySearchClinics(ctx, req.(*FuzzySearchClinicsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ClinicService_GetClinicAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClinicID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClinicServiceServer).GetClinicAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ClinicService_GetClinicAddress_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClinicServiceServer).GetClinicAddress(ctx, req.(*ClinicID))
	}
	return interceptor(ctx, in, info, handler)
}

// ClinicService_ServiceDesc is the grpc.ServiceDesc for ClinicService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClinicService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "coresamples_service.ClinicService",
	HandlerType: (*ClinicServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateClinic",
			Handler:    _ClinicService_CreateClinic_Handler,
		},
		{
			MethodName: "GetClinic",
			Handler:    _ClinicService_GetClinic_Handler,
		},
		{
			MethodName: "ListClinic",
			Handler:    _ClinicService_ListClinic_Handler,
		},
		{
			MethodName: "ListClinicCustomers",
			Handler:    _ClinicService_ListClinicCustomers_Handler,
		},
		{
			MethodName: "GetCustomerClinicNames",
			Handler:    _ClinicService_GetCustomerClinicNames_Handler,
		},
		{
			MethodName: "GetClinicsClinicAccountIDs",
			Handler:    _ClinicService_GetClinicsClinicAccountIDs_Handler,
		},
		{
			MethodName: "ListClinicCustomersByClinicID",
			Handler:    _ClinicService_ListClinicCustomersByClinicID_Handler,
		},
		{
			MethodName: "GetClinicByID",
			Handler:    _ClinicService_GetClinicByID_Handler,
		},
		{
			MethodName: "UpdateClinicNPINumber",
			Handler:    _ClinicService_UpdateClinicNPINumber_Handler,
		},
		{
			MethodName: "EditClinicName",
			Handler:    _ClinicService_EditClinicName_Handler,
		},
		{
			MethodName: "modifyClinicNames",
			Handler:    _ClinicService_ModifyClinicNames_Handler,
		},
		{
			MethodName: "SearchClinicsByName",
			Handler:    _ClinicService_SearchClinicsByName_Handler,
		},
		{
			MethodName: "SignUpClinicForExistingAccount",
			Handler:    _ClinicService_SignUpClinicForExistingAccount_Handler,
		},
		{
			MethodName: "CheckClinicAttributes",
			Handler:    _ClinicService_CheckClinicAttributes_Handler,
		},
		{
			MethodName: "GetFirstCustomerOfClinic",
			Handler:    _ClinicService_GetFirstCustomerOfClinic_Handler,
		},
		{
			MethodName: "FuzzySearchClinics",
			Handler:    _ClinicService_FuzzySearchClinics_Handler,
		},
		{
			MethodName: "GetClinicAddress",
			Handler:    _ClinicService_GetClinicAddress_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/clinic.proto",
}
