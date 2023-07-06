/*
 *
 * Copyright 2023 Intel authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
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
	pkgerrors "github.com/pkg/errors"
)

var (
	InvalidRequestErr = pkgerrors.New("Invalid Request")

	/*
		    tls               = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
			certFile          = flag.String("cert_file", "", "The TLS cert file")
			keyFile           = flag.String("key_file", "", "The TLS key file")
	*/
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
	case pb.CATEGORY_TDX_RTMR:
	case pb.CATEGORY_TPM:
	default:
		log.Println("Invalid measurement category.")
		return "", InvalidRequestErr
	}
	return measurement, err
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

	/*
		    //reserved for TLS
		    var opts []grpc.ServerOption
			if *tls {
				if *certFile == "" {
					*certFile = data.Path("x509/server_cert.pem")
				}
				if *keyFile == "" {
					*keyFile = data.Path("x509/server_key.pem")
				}
				creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
				if err != nil {
					log.Fatalf("Failed to generate credentials: %v", err)
				}
				opts = []grpc.ServerOption{grpc.Creds(creds)}
			}
			grpcServer := grpc.NewServer(opts...)
	*/

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
