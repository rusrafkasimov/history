// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package operation_history

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

// OperationHistoryClient is the client API for OperationHistory service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OperationHistoryClient interface {
	CreateHistory(ctx context.Context, in *OperationHistoryRequest, opts ...grpc.CallOption) (*OperationHistoryResponse, error)
}

type operationHistoryClient struct {
	cc grpc.ClientConnInterface
}

func NewOperationHistoryClient(cc grpc.ClientConnInterface) OperationHistoryClient {
	return &operationHistoryClient{cc}
}

func (c *operationHistoryClient) CreateHistory(ctx context.Context, in *OperationHistoryRequest, opts ...grpc.CallOption) (*OperationHistoryResponse, error) {
	out := new(OperationHistoryResponse)
	err := c.cc.Invoke(ctx, "/proto_operation_history.OperationHistory/CreateHistory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OperationHistoryServer is the server API for OperationHistory service.
// All implementations must embed UnimplementedOperationHistoryServer
// for forward compatibility
type OperationHistoryServer interface {
	CreateHistory(context.Context, *OperationHistoryRequest) (*OperationHistoryResponse, error)
	mustEmbedUnimplementedOperationHistoryServer()
}

// UnimplementedOperationHistoryServer must be embedded to have forward compatible implementations.
type UnimplementedOperationHistoryServer struct {
}

func (UnimplementedOperationHistoryServer) CreateHistory(context.Context, *OperationHistoryRequest) (*OperationHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateHistory not implemented")
}
func (UnimplementedOperationHistoryServer) mustEmbedUnimplementedOperationHistoryServer() {}

// UnsafeOperationHistoryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OperationHistoryServer will
// result in compilation errors.
type UnsafeOperationHistoryServer interface {
	mustEmbedUnimplementedOperationHistoryServer()
}

func RegisterOperationHistoryServer(s grpc.ServiceRegistrar, srv OperationHistoryServer) {
	s.RegisterService(&OperationHistory_ServiceDesc, srv)
}

func _OperationHistory_CreateHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OperationHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OperationHistoryServer).CreateHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto_operation_history.OperationHistory/CreateHistory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OperationHistoryServer).CreateHistory(ctx, req.(*OperationHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OperationHistory_ServiceDesc is the grpc.ServiceDesc for OperationHistory service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OperationHistory_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto_operation_history.OperationHistory",
	HandlerType: (*OperationHistoryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateHistory",
			Handler:    _OperationHistory_CreateHistory_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "operation_history.proto",
}
