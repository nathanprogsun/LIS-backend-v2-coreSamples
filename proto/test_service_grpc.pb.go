// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: proto/test_service.proto

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
	TestService_GetTest_FullMethodName                    = "/coresamples_service.TestService/GetTest"
	TestService_GetTestField_FullMethodName               = "/coresamples_service.TestService/GetTestField"
	TestService_CreateTest_FullMethodName                 = "/coresamples_service.TestService/CreateTest"
	TestService_GetTestTubeTypes_FullMethodName           = "/coresamples_service.TestService/GetTestTubeTypes"
	TestService_GetTestIDsFromTestCodes_FullMethodName    = "/coresamples_service.TestService/GetTestIDsFromTestCodes"
	TestService_GetDuplicateAssayGroupTest_FullMethodName = "/coresamples_service.TestService/GetDuplicateAssayGroupTest"
)

// TestServiceClient is the client API for TestService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TestServiceClient interface {
	GetTest(ctx context.Context, in *GetTestRequest, opts ...grpc.CallOption) (*GetTestResponse, error)
	GetTestField(ctx context.Context, in *GetTestFieldRequest, opts ...grpc.CallOption) (*GetTestResponse, error)
	CreateTest(ctx context.Context, in *CreateTestRequest, opts ...grpc.CallOption) (*CreateTestResponse, error)
	// TODO: implement this
	GetTestTubeTypes(ctx context.Context, in *GetTestTubeTypesRequest, opts ...grpc.CallOption) (*GetTestTubeTypesResponse, error)
	GetTestIDsFromTestCodes(ctx context.Context, in *GetTestIDsFromTestCodesRequest, opts ...grpc.CallOption) (*GetTestIDsFromTestCodesResponse, error)
	GetDuplicateAssayGroupTest(ctx context.Context, in *GetDuplicateAssayGroupTestRequest, opts ...grpc.CallOption) (*GetDuplicateAssayGroupTestResponse, error)
}

type testServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTestServiceClient(cc grpc.ClientConnInterface) TestServiceClient {
	return &testServiceClient{cc}
}

func (c *testServiceClient) GetTest(ctx context.Context, in *GetTestRequest, opts ...grpc.CallOption) (*GetTestResponse, error) {
	out := new(GetTestResponse)
	err := c.cc.Invoke(ctx, TestService_GetTest_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testServiceClient) GetTestField(ctx context.Context, in *GetTestFieldRequest, opts ...grpc.CallOption) (*GetTestResponse, error) {
	out := new(GetTestResponse)
	err := c.cc.Invoke(ctx, TestService_GetTestField_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testServiceClient) CreateTest(ctx context.Context, in *CreateTestRequest, opts ...grpc.CallOption) (*CreateTestResponse, error) {
	out := new(CreateTestResponse)
	err := c.cc.Invoke(ctx, TestService_CreateTest_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testServiceClient) GetTestTubeTypes(ctx context.Context, in *GetTestTubeTypesRequest, opts ...grpc.CallOption) (*GetTestTubeTypesResponse, error) {
	out := new(GetTestTubeTypesResponse)
	err := c.cc.Invoke(ctx, TestService_GetTestTubeTypes_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testServiceClient) GetTestIDsFromTestCodes(ctx context.Context, in *GetTestIDsFromTestCodesRequest, opts ...grpc.CallOption) (*GetTestIDsFromTestCodesResponse, error) {
	out := new(GetTestIDsFromTestCodesResponse)
	err := c.cc.Invoke(ctx, TestService_GetTestIDsFromTestCodes_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *testServiceClient) GetDuplicateAssayGroupTest(ctx context.Context, in *GetDuplicateAssayGroupTestRequest, opts ...grpc.CallOption) (*GetDuplicateAssayGroupTestResponse, error) {
	out := new(GetDuplicateAssayGroupTestResponse)
	err := c.cc.Invoke(ctx, TestService_GetDuplicateAssayGroupTest_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TestServiceServer is the server API for TestService service.
// All implementations must embed UnimplementedTestServiceServer
// for forward compatibility
type TestServiceServer interface {
	GetTest(context.Context, *GetTestRequest) (*GetTestResponse, error)
	GetTestField(context.Context, *GetTestFieldRequest) (*GetTestResponse, error)
	CreateTest(context.Context, *CreateTestRequest) (*CreateTestResponse, error)
	// TODO: implement this
	GetTestTubeTypes(context.Context, *GetTestTubeTypesRequest) (*GetTestTubeTypesResponse, error)
	GetTestIDsFromTestCodes(context.Context, *GetTestIDsFromTestCodesRequest) (*GetTestIDsFromTestCodesResponse, error)
	GetDuplicateAssayGroupTest(context.Context, *GetDuplicateAssayGroupTestRequest) (*GetDuplicateAssayGroupTestResponse, error)
	mustEmbedUnimplementedTestServiceServer()
}

// UnimplementedTestServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTestServiceServer struct {
}

func (UnimplementedTestServiceServer) GetTest(context.Context, *GetTestRequest) (*GetTestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTest not implemented")
}
func (UnimplementedTestServiceServer) GetTestField(context.Context, *GetTestFieldRequest) (*GetTestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTestField not implemented")
}
func (UnimplementedTestServiceServer) CreateTest(context.Context, *CreateTestRequest) (*CreateTestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTest not implemented")
}
func (UnimplementedTestServiceServer) GetTestTubeTypes(context.Context, *GetTestTubeTypesRequest) (*GetTestTubeTypesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTestTubeTypes not implemented")
}
func (UnimplementedTestServiceServer) GetTestIDsFromTestCodes(context.Context, *GetTestIDsFromTestCodesRequest) (*GetTestIDsFromTestCodesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTestIDsFromTestCodes not implemented")
}
func (UnimplementedTestServiceServer) GetDuplicateAssayGroupTest(context.Context, *GetDuplicateAssayGroupTestRequest) (*GetDuplicateAssayGroupTestResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDuplicateAssayGroupTest not implemented")
}
func (UnimplementedTestServiceServer) mustEmbedUnimplementedTestServiceServer() {}

// UnsafeTestServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TestServiceServer will
// result in compilation errors.
type UnsafeTestServiceServer interface {
	mustEmbedUnimplementedTestServiceServer()
}

func RegisterTestServiceServer(s grpc.ServiceRegistrar, srv TestServiceServer) {
	s.RegisterService(&TestService_ServiceDesc, srv)
}

func _TestService_GetTest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).GetTest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TestService_GetTest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).GetTest(ctx, req.(*GetTestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestService_GetTestField_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTestFieldRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).GetTestField(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TestService_GetTestField_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).GetTestField(ctx, req.(*GetTestFieldRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestService_CreateTest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateTestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).CreateTest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TestService_CreateTest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).CreateTest(ctx, req.(*CreateTestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestService_GetTestTubeTypes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTestTubeTypesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).GetTestTubeTypes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TestService_GetTestTubeTypes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).GetTestTubeTypes(ctx, req.(*GetTestTubeTypesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestService_GetTestIDsFromTestCodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTestIDsFromTestCodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).GetTestIDsFromTestCodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TestService_GetTestIDsFromTestCodes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).GetTestIDsFromTestCodes(ctx, req.(*GetTestIDsFromTestCodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _TestService_GetDuplicateAssayGroupTest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDuplicateAssayGroupTestRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TestServiceServer).GetDuplicateAssayGroupTest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: TestService_GetDuplicateAssayGroupTest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TestServiceServer).GetDuplicateAssayGroupTest(ctx, req.(*GetDuplicateAssayGroupTestRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// TestService_ServiceDesc is the grpc.ServiceDesc for TestService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TestService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "coresamples_service.TestService",
	HandlerType: (*TestServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTest",
			Handler:    _TestService_GetTest_Handler,
		},
		{
			MethodName: "GetTestField",
			Handler:    _TestService_GetTestField_Handler,
		},
		{
			MethodName: "CreateTest",
			Handler:    _TestService_CreateTest_Handler,
		},
		{
			MethodName: "GetTestTubeTypes",
			Handler:    _TestService_GetTestTubeTypes_Handler,
		},
		{
			MethodName: "GetTestIDsFromTestCodes",
			Handler:    _TestService_GetTestIDsFromTestCodes_Handler,
		},
		{
			MethodName: "GetDuplicateAssayGroupTest",
			Handler:    _TestService_GetDuplicateAssayGroupTest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/test_service.proto",
}
