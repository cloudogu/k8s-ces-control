// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: health.proto

package health

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

// DoguHealthClient is the client API for DoguHealth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DoguHealthClient interface {
	GetByName(ctx context.Context, in *DoguHealthRequest, opts ...grpc.CallOption) (*DoguHealthResponse, error)
	GetByNames(ctx context.Context, in *DoguHealthListRequest, opts ...grpc.CallOption) (*DoguHealthMapResponse, error)
	GetAll(ctx context.Context, in *DoguHealthAllRequest, opts ...grpc.CallOption) (*DoguHealthMapResponse, error)
}

type doguHealthClient struct {
	cc grpc.ClientConnInterface
}

func NewDoguHealthClient(cc grpc.ClientConnInterface) DoguHealthClient {
	return &doguHealthClient{cc}
}

func (c *doguHealthClient) GetByName(ctx context.Context, in *DoguHealthRequest, opts ...grpc.CallOption) (*DoguHealthResponse, error) {
	out := new(DoguHealthResponse)
	err := c.cc.Invoke(ctx, "/health.DoguHealth/GetByName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *doguHealthClient) GetByNames(ctx context.Context, in *DoguHealthListRequest, opts ...grpc.CallOption) (*DoguHealthMapResponse, error) {
	out := new(DoguHealthMapResponse)
	err := c.cc.Invoke(ctx, "/health.DoguHealth/GetByNames", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *doguHealthClient) GetAll(ctx context.Context, in *DoguHealthAllRequest, opts ...grpc.CallOption) (*DoguHealthMapResponse, error) {
	out := new(DoguHealthMapResponse)
	err := c.cc.Invoke(ctx, "/health.DoguHealth/GetAll", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DoguHealthServer is the server API for DoguHealth service.
// All implementations must embed UnimplementedDoguHealthServer
// for forward compatibility
type DoguHealthServer interface {
	GetByName(context.Context, *DoguHealthRequest) (*DoguHealthResponse, error)
	GetByNames(context.Context, *DoguHealthListRequest) (*DoguHealthMapResponse, error)
	GetAll(context.Context, *DoguHealthAllRequest) (*DoguHealthMapResponse, error)
	mustEmbedUnimplementedDoguHealthServer()
}

// UnimplementedDoguHealthServer must be embedded to have forward compatible implementations.
type UnimplementedDoguHealthServer struct {
}

func (UnimplementedDoguHealthServer) GetByName(context.Context, *DoguHealthRequest) (*DoguHealthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByName not implemented")
}
func (UnimplementedDoguHealthServer) GetByNames(context.Context, *DoguHealthListRequest) (*DoguHealthMapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetByNames not implemented")
}
func (UnimplementedDoguHealthServer) GetAll(context.Context, *DoguHealthAllRequest) (*DoguHealthMapResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAll not implemented")
}
func (UnimplementedDoguHealthServer) mustEmbedUnimplementedDoguHealthServer() {}

// UnsafeDoguHealthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DoguHealthServer will
// result in compilation errors.
type UnsafeDoguHealthServer interface {
	mustEmbedUnimplementedDoguHealthServer()
}

func RegisterDoguHealthServer(s grpc.ServiceRegistrar, srv DoguHealthServer) {
	s.RegisterService(&DoguHealth_ServiceDesc, srv)
}

func _DoguHealth_GetByName_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DoguHealthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DoguHealthServer).GetByName(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/health.DoguHealth/GetByName",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DoguHealthServer).GetByName(ctx, req.(*DoguHealthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoguHealth_GetByNames_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DoguHealthListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DoguHealthServer).GetByNames(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/health.DoguHealth/GetByNames",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DoguHealthServer).GetByNames(ctx, req.(*DoguHealthListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DoguHealth_GetAll_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DoguHealthAllRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DoguHealthServer).GetAll(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/health.DoguHealth/GetAll",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DoguHealthServer).GetAll(ctx, req.(*DoguHealthAllRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DoguHealth_ServiceDesc is the grpc.ServiceDesc for DoguHealth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DoguHealth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "health.DoguHealth",
	HandlerType: (*DoguHealthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetByName",
			Handler:    _DoguHealth_GetByName_Handler,
		},
		{
			MethodName: "GetByNames",
			Handler:    _DoguHealth_GetByNames_Handler,
		},
		{
			MethodName: "GetAll",
			Handler:    _DoguHealth_GetAll_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "health.proto",
}
