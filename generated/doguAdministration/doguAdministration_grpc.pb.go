// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: doguAdministration.proto

package doguAdministration

import (
	context "context"
	types "github.com/cloudogu/k8s-ces-control/generated/types"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DoguAdministrationClient is the client API for DoguAdministration service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DoguAdministrationClient interface {
	// getDogusToAdminList returns the list of dogus to administrate
	GetDoguList(ctx context.Context, in *DoguListRequest, opts ...grpc.CallOption) (*DoguListResponse, error)
	// StartDogu starts a dogu
	StartDogu(ctx context.Context, in *DoguAdministrationRequest, opts ...grpc.CallOption) (*types.BasicResponse, error)
	// StopDogu stops a dogu
	StopDogu(ctx context.Context, in *DoguAdministrationRequest, opts ...grpc.CallOption) (*types.BasicResponse, error)
	// RestartDogu restarts a dogu
	RestartDogu(ctx context.Context, in *DoguAdministrationRequest, opts ...grpc.CallOption) (*types.BasicResponse, error)
}

type doguAdministrationClient struct {
	cc grpc.ClientConnInterface
}

func NewDoguAdministrationClient(cc grpc.ClientConnInterface) DoguAdministrationClient {
	return &doguAdministrationClient{cc}
}

func (c *doguAdministrationClient) GetDoguList(ctx context.Context, in *DoguListRequest, opts ...grpc.CallOption) (*DoguListResponse, error) {
	out := new(DoguListResponse)
	err := c.cc.Invoke(ctx, "/doguAdministration.DoguAdministration/GetDoguList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *doguAdministrationClient) StartDogu(ctx context.Context, in *DoguAdministrationRequest, opts ...grpc.CallOption) (*types.BasicResponse, error) {
	out := new(types.BasicResponse)
	err := c.cc.Invoke(ctx, "/doguAdministration.DoguAdministration/StartDogu", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *doguAdministrationClient) StopDogu(ctx context.Context, in *DoguAdministrationRequest, opts ...grpc.CallOption) (*types.BasicResponse, error) {
	out := new(types.BasicResponse)
	err := c.cc.Invoke(ctx, "/doguAdministration.DoguAdministration/StopDogu", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *doguAdministrationClient) RestartDogu(ctx context.Context, in *DoguAdministrationRequest, opts ...grpc.CallOption) (*types.BasicResponse, error) {
	out := new(types.BasicResponse)
	err := c.cc.Invoke(ctx, "/doguAdministration.DoguAdministration/RestartDogu", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DoguAdministrationServer is the server API for DoguAdministration service.
// All implementations must embed UnimplementedDoguAdministrationServer
// for forward compatibility
type DoguAdministrationServer interface {
	// getDogusToAdminList returns the list of dogus to administrate
	GetDoguList(context.Context, *DoguListRequest) (*DoguListResponse, error)
	// StartDogu starts a dogu
	StartDogu(context.Context, *DoguAdministrationRequest) (*types.BasicResponse, error)
	// StopDogu stops a dogu
	StopDogu(context.Context, *DoguAdministrationRequest) (*types.BasicResponse, error)
	// RestartDogu restarts a dogu
	RestartDogu(context.Context, *DoguAdministrationRequest) (*types.BasicResponse, error)
	mustEmbedUnimplementedDoguAdministrationServer()
}

// UnimplementedDoguAdministrationServer must be embedded to have forward compatible implementations.
type UnimplementedDoguAdministrationServer struct {
}

func (UnimplementedDoguAdministrationServer) GetDoguList(context.Context, *DoguListRequest) (*DoguListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDoguList not implemented")
}
func (UnimplementedDoguAdministrationServer) StartDogu(context.Context, *DoguAdministrationRequest) (*types.BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartDogu not implemented")
}
func (UnimplementedDoguAdministrationServer) StopDogu(context.Context, *DoguAdministrationRequest) (*types.BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopDogu not implemented")
}
func (UnimplementedDoguAdministrationServer) RestartDogu(context.Context, *DoguAdministrationRequest) (*types.BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RestartDogu not implemented")
}
func (UnimplementedDoguAdministrationServer) mustEmbedUnimplementedDoguAdministrationServer() {}

// UnsafeDoguAdministrationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DoguAdministrationServer will
// result in compilation errors.
type UnsafeDoguAdministrationServer interface {
	mustEmbedUnimplementedDoguAdministrationServer()
}

func RegisterDoguAdministrationServer(s grpc.ServiceRegistrar, srv DoguAdministrationServer) {
	s.RegisterService(&DoguAdministration_ServiceDesc, srv)
}

func _DoguAdministration_GetDoguList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DoguListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DoguAdministrationServer).GetDoguList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/doguAdministration.DoguAdministration/GetDoguList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DoguAdministrationServer).GetDoguList(ctx, req.(*DoguListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoguAdministration_StartDogu_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DoguAdministrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DoguAdministrationServer).StartDogu(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/doguAdministration.DoguAdministration/StartDogu",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DoguAdministrationServer).StartDogu(ctx, req.(*DoguAdministrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoguAdministration_StopDogu_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DoguAdministrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DoguAdministrationServer).StopDogu(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/doguAdministration.DoguAdministration/StopDogu",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DoguAdministrationServer).StopDogu(ctx, req.(*DoguAdministrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoguAdministration_RestartDogu_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DoguAdministrationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DoguAdministrationServer).RestartDogu(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/doguAdministration.DoguAdministration/RestartDogu",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DoguAdministrationServer).RestartDogu(ctx, req.(*DoguAdministrationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DoguAdministration_ServiceDesc is the grpc.ServiceDesc for DoguAdministration service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DoguAdministration_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "doguAdministration.DoguAdministration",
	HandlerType: (*DoguAdministrationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDoguList",
			Handler:    _DoguAdministration_GetDoguList_Handler,
		},
		{
			MethodName: "StartDogu",
			Handler:    _DoguAdministration_StartDogu_Handler,
		},
		{
			MethodName: "StopDogu",
			Handler:    _DoguAdministration_StopDogu_Handler,
		},
		{
			MethodName: "RestartDogu",
			Handler:    _DoguAdministration_RestartDogu_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "doguAdministration.proto",
}
