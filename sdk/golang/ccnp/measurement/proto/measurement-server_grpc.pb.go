// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.11.4
// source: proto/measurement-server.proto

package getMeasurement

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

// MeasurementClient is the client API for Measurement service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MeasurementClient interface {
	GetMeasurement(ctx context.Context, in *GetMeasurementRequest, opts ...grpc.CallOption) (*GetMeasurementReply, error)
}

type measurementClient struct {
	cc grpc.ClientConnInterface
}

func NewMeasurementClient(cc grpc.ClientConnInterface) MeasurementClient {
	return &measurementClient{cc}
}

func (c *measurementClient) GetMeasurement(ctx context.Context, in *GetMeasurementRequest, opts ...grpc.CallOption) (*GetMeasurementReply, error) {
	out := new(GetMeasurementReply)
	err := c.cc.Invoke(ctx, "/measurement.Measurement/GetMeasurement", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MeasurementServer is the server API for Measurement service.
// All implementations must embed UnimplementedMeasurementServer
// for forward compatibility
type MeasurementServer interface {
	GetMeasurement(context.Context, *GetMeasurementRequest) (*GetMeasurementReply, error)
	mustEmbedUnimplementedMeasurementServer()
}

// UnimplementedMeasurementServer must be embedded to have forward compatible implementations.
type UnimplementedMeasurementServer struct {
}

func (UnimplementedMeasurementServer) GetMeasurement(context.Context, *GetMeasurementRequest) (*GetMeasurementReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMeasurement not implemented")
}
func (UnimplementedMeasurementServer) mustEmbedUnimplementedMeasurementServer() {}

// UnsafeMeasurementServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MeasurementServer will
// result in compilation errors.
type UnsafeMeasurementServer interface {
	mustEmbedUnimplementedMeasurementServer()
}

func RegisterMeasurementServer(s grpc.ServiceRegistrar, srv MeasurementServer) {
	s.RegisterService(&Measurement_ServiceDesc, srv)
}

func _Measurement_GetMeasurement_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMeasurementRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MeasurementServer).GetMeasurement(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/measurement.Measurement/GetMeasurement",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MeasurementServer).GetMeasurement(ctx, req.(*GetMeasurementRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Measurement_ServiceDesc is the grpc.ServiceDesc for Measurement service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Measurement_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "measurement.Measurement",
	HandlerType: (*MeasurementServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetMeasurement",
			Handler:    _Measurement_GetMeasurement_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/measurement-server.proto",
}
