/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"context"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	pb "github.com/intel/confidential-cloud-native-primitives/service/measurement-server/proto"
	resources "github.com/intel/confidential-cloud-native-primitives/service/measurement-server/resources"
	pkgerrors "github.com/pkg/errors"
)

var (
	InvalidRequestErr = pkgerrors.New("Invalid Request")
)

const (
	protocol = "unix"
	sockAddr = "/run/ccnp/uds/measurement.sock"
)

type measurementServer struct {
	pb.UnimplementedMeasurementServer
}

func getContainerMeasurement(measurementReq *pb.GetMeasurementRequest) (string, error) {
	// not implemented
	log.Println("Not implemented.")
	return "", nil
}

func getPaasMeasurement(measurementReq *pb.GetMeasurementRequest) (string, error) {
	var category pb.CATEGORY
	var measurement string
	var err error

	category = measurementReq.MeasurementCategory

	switch category {
	case pb.CATEGORY_TEE_REPORT:
		measurement, err = getTeeReport(measurementReq)
	case pb.CATEGORY_TDX_RTMR:
		var device string
		r := resources.NewTdxResource()
		device, err = r.FindDeviceAvailable()
		if err != nil {
			return "", err
		}
		measurement, err = r.GetRTMRMeasurement(device, measurementReq.ReportData, int(measurementReq.RegisterIndex))
	case pb.CATEGORY_TPM:
		measurement, err = resources.GetTpmMeasurement(int(measurementReq.RegisterIndex))
	default:
		log.Println("Invalid measurement category.")
		return "", InvalidRequestErr
	}
	return measurement, err
}

func getTeeReport(measurementReq *pb.GetMeasurementRequest) (string, error) {

	reportData := measurementReq.ReportData

	r := resources.NewBaseTeeResource()
	device, err := r.FindDeviceAvailable()
	if err != nil {
		return "", err
	}

	report, err := r.GetReport(device, reportData)
	if err != nil {
		return "", err
	}

	return report, nil
}

func (*measurementServer) GetMeasurement(ctx context.Context, measurementReq *pb.GetMeasurementRequest) (*pb.GetMeasurementReply, error) {
	var measurement_type pb.TYPE
	var measurement string
	var err error

	measurement_type = measurementReq.MeasurementType

	switch measurement_type {
	case pb.TYPE_SAAS:
		measurement, err = getContainerMeasurement(measurementReq)
	case pb.TYPE_PAAS:
		measurement, err = getPaasMeasurement(measurementReq)
	default:
		log.Println("Invalid measurement type.")
		return &pb.GetMeasurementReply{}, InvalidRequestErr

	}

	if err != nil {
		return &pb.GetMeasurementReply{}, err
	}
	return &pb.GetMeasurementReply{Measurement: measurement}, nil
}

func (*measurementServer) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (*measurementServer) Watch(in *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	return nil
}

func newServer() *measurementServer {
	s := &measurementServer{}
	return s
}

func main() {
	if _, err := os.Stat(sockAddr); !os.IsNotExist(err) {
		if err := os.RemoveAll(sockAddr); err != nil {
			log.Fatal(err)
		}
	}

	lis, err := net.Listen(protocol, sockAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	healthServer := health.NewServer()

	pb.RegisterMeasurementServer(grpcServer, newServer())
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	log.Printf("server listening at %v", lis.Addr())
	reflection.Register(grpcServer)
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
