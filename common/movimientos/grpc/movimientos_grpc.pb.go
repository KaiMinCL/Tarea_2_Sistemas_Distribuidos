// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package grpc

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

// MovimientosServiceClient is the client API for MovimientosService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MovimientosServiceClient interface {
	RegistrarMovimiento(ctx context.Context, in *MovimientoRequest, opts ...grpc.CallOption) (*MovimientoResponse, error)
}

type movimientosServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMovimientosServiceClient(cc grpc.ClientConnInterface) MovimientosServiceClient {
	return &movimientosServiceClient{cc}
}

func (c *movimientosServiceClient) RegistrarMovimiento(ctx context.Context, in *MovimientoRequest, opts ...grpc.CallOption) (*MovimientoResponse, error) {
	out := new(MovimientoResponse)
	err := c.cc.Invoke(ctx, "/movimientos.MovimientosService/RegistrarMovimiento", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MovimientosServiceServer is the server API for MovimientosService service.
// All implementations must embed UnimplementedMovimientosServiceServer
// for forward compatibility
type MovimientosServiceServer interface {
	RegistrarMovimiento(context.Context, *MovimientoRequest) (*MovimientoResponse, error)
	mustEmbedUnimplementedMovimientosServiceServer()
}

// UnimplementedMovimientosServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMovimientosServiceServer struct {
}

func (UnimplementedMovimientosServiceServer) RegistrarMovimiento(context.Context, *MovimientoRequest) (*MovimientoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegistrarMovimiento not implemented")
}
func (UnimplementedMovimientosServiceServer) mustEmbedUnimplementedMovimientosServiceServer() {}

// UnsafeMovimientosServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MovimientosServiceServer will
// result in compilation errors.
type UnsafeMovimientosServiceServer interface {
	mustEmbedUnimplementedMovimientosServiceServer()
}

func RegisterMovimientosServiceServer(s grpc.ServiceRegistrar, srv MovimientosServiceServer) {
	s.RegisterService(&MovimientosService_ServiceDesc, srv)
}

func _MovimientosService_RegistrarMovimiento_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MovimientoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MovimientosServiceServer).RegistrarMovimiento(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/movimientos.MovimientosService/RegistrarMovimiento",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MovimientosServiceServer).RegistrarMovimiento(ctx, req.(*MovimientoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MovimientosService_ServiceDesc is the grpc.ServiceDesc for MovimientosService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MovimientosService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "movimientos.MovimientosService",
	HandlerType: (*MovimientosServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegistrarMovimiento",
			Handler:    _MovimientosService_RegistrarMovimiento_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "movimientos.proto",
}