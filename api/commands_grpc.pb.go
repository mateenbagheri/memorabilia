// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             v5.27.1
// source: api/commands.proto

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	Commands_Echo_FullMethodName = "/commands.Commands/Echo"
)

// CommandsClient is the client API for Commands service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommandsClient interface {
	Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoResponse, error)
}

type commandsClient struct {
	cc grpc.ClientConnInterface
}

func NewCommandsClient(cc grpc.ClientConnInterface) CommandsClient {
	return &commandsClient{cc}
}

func (c *commandsClient) Echo(ctx context.Context, in *EchoRequest, opts ...grpc.CallOption) (*EchoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EchoResponse)
	err := c.cc.Invoke(ctx, Commands_Echo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommandsServer is the server API for Commands service.
// All implementations must embed UnimplementedCommandsServer
// for forward compatibility
type CommandsServer interface {
	Echo(context.Context, *EchoRequest) (*EchoResponse, error)
	mustEmbedUnimplementedCommandsServer()
}

// UnimplementedCommandsServer must be embedded to have forward compatible implementations.
type UnimplementedCommandsServer struct {
}

func (UnimplementedCommandsServer) Echo(context.Context, *EchoRequest) (*EchoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Echo not implemented")
}
func (UnimplementedCommandsServer) mustEmbedUnimplementedCommandsServer() {}

// UnsafeCommandsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommandsServer will
// result in compilation errors.
type UnsafeCommandsServer interface {
	mustEmbedUnimplementedCommandsServer()
}

func RegisterCommandsServer(s grpc.ServiceRegistrar, srv CommandsServer) {
	s.RegisterService(&Commands_ServiceDesc, srv)
}

func _Commands_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommandsServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Commands_Echo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommandsServer).Echo(ctx, req.(*EchoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Commands_ServiceDesc is the grpc.ServiceDesc for Commands service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Commands_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "commands.Commands",
	HandlerType: (*CommandsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _Commands_Echo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/commands.proto",
}
