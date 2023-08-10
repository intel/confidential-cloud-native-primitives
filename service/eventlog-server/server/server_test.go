/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */
package main

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/intel/confidential-cloud-native-primitives/service/eventlog-server/proto"
)

const (
	INVALID_EVENTLOG_LEVEL    pb.LEVEL    = 9
	INVALID_EVENTLOG_CATEGORY pb.CATEGORY = 9
)

var lis *bufconn.Listener

func initTestServer(ctx context.Context) {
	buffer := 1024 * 1024
	lis = bufconn.Listen(buffer)

	server := grpc.NewServer()
	pb.RegisterEventlogServer(server, &eventlogServer{})
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()
}

func TestEventlogServerGetEventlog(t *testing.T) {
	ctx := context.Background()
	initTestServer(ctx)

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(
		func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("failed to connect to server: %v", err)
		return
	}
	defer conn.Close()
	client := pb.NewEventlogClient(conn)

	type expectation struct {
		out *pb.GetEventlogReply
		err error
	}

	tests := map[string]struct {
		in       *pb.GetEventlogRequest
		expected expectation
	}{
		"Empty_Request": {
			in: &pb.GetEventlogRequest{},
			expected: expectation{
				out: &pb.GetEventlogReply{
					EventlogDataLoc: "/run/ccnp-eventlog/eventlog.log",
				},
				err: nil,
			},
		},
		"Invalid_Eventlog_Level": {
			in: &pb.GetEventlogRequest{
				EventlogLevel: INVALID_EVENTLOG_LEVEL,
			},
			expected: expectation{
				out: &pb.GetEventlogReply{},
				err: errors.New("rpc error: code = Unknown desc = Invalid Request"),
			},
		},
		"Invalid_Eventlog_Category": {
			in: &pb.GetEventlogRequest{
				EventlogLevel:    pb.LEVEL_PAAS,
				EventlogCategory: INVALID_EVENTLOG_CATEGORY,
			},
			expected: expectation{
				out: &pb.GetEventlogReply{},
				err: errors.New("rpc error: code = Unknown desc = Invalid Request"),
			},
		},
		"Request_on_SAAS_Level_Eventlog": {
			in: &pb.GetEventlogRequest{
				EventlogLevel: pb.LEVEL_SAAS,
			},
			expected: expectation{
				out: &pb.GetEventlogReply{
					EventlogDataLoc: "/run/ccnp-eventlog/eventlog.log",
				},
				err: nil,
			},
		},
		"Request_on_TPM_Eventlog_without_TPM_Support": {
			in: &pb.GetEventlogRequest{
				EventlogLevel:    pb.LEVEL_PAAS,
				EventlogCategory: pb.CATEGORY_TPM_EVENTLOG,
			},
			expected: expectation{
				out: &pb.GetEventlogReply{},
				err: errors.New("rpc error: code = Unknown desc = stat /sys/kernel/security/tpm0/binary_bios_measurements: no such file or directory"),
			},
		},
		"Request_on_TPM_Eventlog_with_Options_without_TPM_support": {
			in: &pb.GetEventlogRequest{
				EventlogLevel:    pb.LEVEL_PAAS,
				EventlogCategory: pb.CATEGORY_TPM_EVENTLOG,
				StartPosition:    1,
				Count:            5,
			},
			expected: expectation{
				out: &pb.GetEventlogReply{},
				err: errors.New("rpc error: code = Unknown desc = stat /sys/kernel/security/tpm0/binary_bios_measurements: no such file or directory"),
			},
		},
		"Request_on_Basic_TDX_Eventlog": {
			in: &pb.GetEventlogRequest{
				EventlogLevel:    pb.LEVEL_PAAS,
				EventlogCategory: pb.CATEGORY_TDX_EVENTLOG,
			},
			expected: expectation{
				out: &pb.GetEventlogReply{
					EventlogDataLoc: "/run/ccnp-eventlog/eventlog.log",
				},
				err: nil,
			},
		},
		"Request_on_TDX_Eventlog_with_options": {
			in: &pb.GetEventlogRequest{
				EventlogLevel:    pb.LEVEL_PAAS,
				EventlogCategory: pb.CATEGORY_TDX_EVENTLOG,
				StartPosition:    1,
				Count:            5,
			},
			expected: expectation{
				out: &pb.GetEventlogReply{
					EventlogDataLoc: "/run/ccnp-eventlog/eventlog.log",
				},
				err: nil,
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.GetEventlog(ctx, tt.in)
			if err != nil {
				if tt.expected.err == nil {
					t.Errorf("Err -> \nWant: nil\nGot: %q\n", err)
				} else {
					if tt.expected.err.Error() != err.Error() {
						t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.err, err)
					}
				}
			} else {
				if tt.expected.out.EventlogDataLoc != out.EventlogDataLoc {
					t.Errorf("Out -> \nWant: %q\nGot : %q", tt.expected.out, out)
				}
			}

		})
	}
}
