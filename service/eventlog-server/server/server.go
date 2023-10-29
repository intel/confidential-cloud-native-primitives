/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	pb "github.com/intel/confidential-cloud-native-primitives/service/eventlog-server/proto"
	resources "github.com/intel/confidential-cloud-native-primitives/service/eventlog-server/resources"
	pkgerrors "github.com/pkg/errors"
)

var (
	InvalidRequestErr = pkgerrors.New("Invalid Request")
)

const (
	RUNTIME_EVENT_LOG_DIR  = "/run/ccnp-eventlog/"
	FILENAME               = "eventlog.log"
	protocol               = "unix"
	sockAddr               = "/run/ccnp/uds/eventlog.sock"
	MAX_CONCURRENT_STREAMS = 100
)

type eventlogServer struct {
	pb.UnimplementedEventlogServer
}

func getContainerLevelEventlog(eventlogReq *pb.GetEventlogRequest) (string, error) {
	// not implemented
	return "", nil
}

func getPaasLevelEventlog(eventlogReq *pb.GetEventlogRequest) (string, error) {
	var category pb.CATEGORY
	var eventlog string
	var err error

	category = eventlogReq.EventlogCategory

	switch category {
	case pb.CATEGORY_TPM_EVENTLOG:
		eventlog, err = resources.GetTpmEventlog(int(eventlogReq.StartPosition), int(eventlogReq.Count))
	case pb.CATEGORY_TDX_EVENTLOG:
		eventlog, err = resources.GetTdxEventlog(int(eventlogReq.StartPosition), int(eventlogReq.Count))
	default:
		log.Println("Invalid eventlog category.")
		return "", InvalidRequestErr
	}
	return eventlog, err
}

func (*eventlogServer) GetEventlog(ctx context.Context, eventlogReq *pb.GetEventlogRequest) (*pb.GetEventlogReply, error) {
	var eventlog_level pb.LEVEL
	var eventlog string
	var err error

	eventlog_level = eventlogReq.EventlogLevel

	switch eventlog_level {
	case pb.LEVEL_SAAS:
		eventlog, err = getContainerLevelEventlog(eventlogReq)
	case pb.LEVEL_PAAS:
		eventlog, err = getPaasLevelEventlog(eventlogReq)
	default:
		log.Println("Invalid eventlog level.")
		return &pb.GetEventlogReply{}, InvalidRequestErr
	}

	if err != nil {
		return &pb.GetEventlogReply{}, err
	}

	if _, err := os.Stat(RUNTIME_EVENT_LOG_DIR); os.IsNotExist(err) {
		err := os.MkdirAll(RUNTIME_EVENT_LOG_DIR, os.ModePerm)
		if err != nil {
			return &pb.GetEventlogReply{}, err
		}
	}

	file, err := os.Create(fmt.Sprintf("%s%s", RUNTIME_EVENT_LOG_DIR, FILENAME))
	if err != nil {
		log.Println("Error creating event log file in", RUNTIME_EVENT_LOG_DIR)
		return &pb.GetEventlogReply{}, err
	}
	defer file.Close()

	_, err = file.WriteString(eventlog)
	if err != nil {
		log.Println("Error writing event log file in", RUNTIME_EVENT_LOG_DIR+FILENAME)
		return &pb.GetEventlogReply{}, err
	}

	return &pb.GetEventlogReply{EventlogDataLoc: RUNTIME_EVENT_LOG_DIR + FILENAME}, nil
}

func (*eventlogServer) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

func (*eventlogServer) Watch(in *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	return nil
}

func newServer() *eventlogServer {
	s := &eventlogServer{}
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

	opts := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(MAX_CONCURRENT_STREAMS),
	}

	grpcServer := grpc.NewServer(opts...)
	healthServer := health.NewServer()

	pb.RegisterEventlogServer(grpcServer, newServer())
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)

	log.Printf("server listening at %v", lis.Addr())
	reflection.Register(grpcServer)
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
