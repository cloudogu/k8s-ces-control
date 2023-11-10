// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: debugMode.proto

package debug

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

// DebugModeClient is the client API for DebugMode service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DebugModeClient interface {
	Enable(ctx context.Context, in *ToggleDebugModeRequest, opts ...grpc.CallOption) (*types.BasicResponse, error)
	Disable(ctx context.Context, in *ToggleDebugModeRequest, opts ...grpc.CallOption) (*types.BasicResponse, error)
	Status(ctx context.Context, in *types.BasicRequest, opts ...grpc.CallOption) (*DebugModeStatusResponse, error)
}

type debugModeClient struct {
	cc grpc.ClientConnInterface
}

func NewDebugModeClient(cc grpc.ClientConnInterface) DebugModeClient {
	return &debugModeClient{cc}
}

func (c *debugModeClient) Enable(ctx context.Context, in *ToggleDebugModeRequest, opts ...grpc.CallOption) (*types.BasicResponse, error) {
	out := new(types.BasicResponse)
	err := c.cc.Invoke(ctx, "/debug.DebugMode/Enable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debugModeClient) Disable(ctx context.Context, in *ToggleDebugModeRequest, opts ...grpc.CallOption) (*types.BasicResponse, error) {
	out := new(types.BasicResponse)
	err := c.cc.Invoke(ctx, "/debug.DebugMode/Disable", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *debugModeClient) Status(ctx context.Context, in *types.BasicRequest, opts ...grpc.CallOption) (*DebugModeStatusResponse, error) {
	out := new(DebugModeStatusResponse)
	err := c.cc.Invoke(ctx, "/debug.DebugMode/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DebugModeServer is the server API for DebugMode service.
// All implementations must embed UnimplementedDebugModeServer
// for forward compatibility
type DebugModeServer interface {
	Enable(context.Context, *ToggleDebugModeRequest) (*types.BasicResponse, error)
	Disable(context.Context, *ToggleDebugModeRequest) (*types.BasicResponse, error)
	Status(context.Context, *types.BasicRequest) (*DebugModeStatusResponse, error)
	mustEmbedUnimplementedDebugModeServer()
}

// UnimplementedDebugModeServer must be embedded to have forward compatible implementations.
type UnimplementedDebugModeServer struct {
}

func (UnimplementedDebugModeServer) Enable(context.Context, *ToggleDebugModeRequest) (*types.BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Enable not implemented")
}
func (UnimplementedDebugModeServer) Disable(context.Context, *ToggleDebugModeRequest) (*types.BasicResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Disable not implemented")
}
func (UnimplementedDebugModeServer) Status(context.Context, *types.BasicRequest) (*DebugModeStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedDebugModeServer) mustEmbedUnimplementedDebugModeServer() {}

// UnsafeDebugModeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DebugModeServer will
// result in compilation errors.
type UnsafeDebugModeServer interface {
	mustEmbedUnimplementedDebugModeServer()
}

func RegisterDebugModeServer(s grpc.ServiceRegistrar, srv DebugModeServer) {
	s.RegisterService(&DebugMode_ServiceDesc, srv)
}

func _DebugMode_Enable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ToggleDebugModeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebugModeServer).Enable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/debug.DebugMode/Enable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebugModeServer).Enable(ctx, req.(*ToggleDebugModeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebugMode_Disable_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ToggleDebugModeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebugModeServer).Disable(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/debug.DebugMode/Disable",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebugModeServer).Disable(ctx, req.(*ToggleDebugModeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DebugMode_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(types.BasicRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DebugModeServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/debug.DebugMode/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DebugModeServer).Status(ctx, req.(*types.BasicRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DebugMode_ServiceDesc is the grpc.ServiceDesc for DebugMode service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DebugMode_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "debug.DebugMode",
	HandlerType: (*DebugModeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Enable",
			Handler:    _DebugMode_Enable_Handler,
		},
		{
			MethodName: "Disable",
			Handler:    _DebugMode_Disable_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _DebugMode_Status_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "debugMode.proto",
}