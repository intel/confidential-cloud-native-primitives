/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */
package main

import (
	"context"
	"encoding/base64"
	"errors"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/intel/confidential-cloud-native-primitives/service/measurement-server/proto"
)

const (
	INVALID_MEASUREMENT_TYPE     pb.TYPE     = 9
	INVALID_MEASUREMENT_CATEGORY pb.CATEGORY = 9
)

var lis *bufconn.Listener

func initTestServer(ctx context.Context) {
	buffer := 1024 * 1024
	lis = bufconn.Listen(buffer)

	server := grpc.NewServer()
	pb.RegisterMeasurementServer(server, &measurementServer{})
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()
}

func TestMeasurementServerGetMeasurement(t *testing.T) {
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
	client := pb.NewMeasurementClient(conn)

	type expectation struct {
		err error
	}

	tests := map[string]struct {
		in       *pb.GetMeasurementRequest
		expected expectation
	}{
		"Empty_Request": {
			in: &pb.GetMeasurementRequest{},
			expected: expectation{
				err: nil,
			},
		},
		"Invalid_Measurement_Type": {
			in: &pb.GetMeasurementRequest{
				MeasurementType: INVALID_MEASUREMENT_TYPE,
			},
			expected: expectation{
				err: errors.New("rpc error: code = Unknown desc = Invalid Request"),
			},
		},
		"Invalid_Measurement_Category": {
			in: &pb.GetMeasurementRequest{
				MeasurementType:     pb.TYPE_PAAS,
				MeasurementCategory: INVALID_MEASUREMENT_CATEGORY,
			},
			expected: expectation{
				err: errors.New("rpc error: code = Unknown desc = Invalid Request"),
			},
		},
		"Request_on_SAAS_Measurement": {
			in: &pb.GetMeasurementRequest{
				MeasurementType: pb.TYPE_SAAS,
			},
			expected: expectation{
				err: nil,
			},
		},
		"Request_on_TPM_Measurement_without_TPM_Support": {
			in: &pb.GetMeasurementRequest{
				MeasurementType:     pb.TYPE_PAAS,
				MeasurementCategory: pb.CATEGORY_TPM,
			},
			expected: expectation{
				err: errors.New("rpc error: code = Unknown desc = No applicable device found."),
			},
		},
		"Request_on_TEE_Report_Measurement_with_Options_TDX": {
			in: &pb.GetMeasurementRequest{
				MeasurementType:     pb.TYPE_PAAS,
				MeasurementCategory: pb.CATEGORY_TEE_REPORT,
			},
			expected: expectation{
				err: nil,
			},
		},
		"Request_on_TDX_RTMR_Measurement": {
			in: &pb.GetMeasurementRequest{
				MeasurementType:     pb.TYPE_PAAS,
				MeasurementCategory: pb.CATEGORY_TDX_RTMR,
				RegisterIndex:       0,
			},
			expected: expectation{
				err: nil,
			},
		},
		"Request_on_TDX_Measurement_with_Invalid_Register_Index": {
			in: &pb.GetMeasurementRequest{
				MeasurementType:     pb.TYPE_PAAS,
				MeasurementCategory: pb.CATEGORY_TDX_RTMR,
				RegisterIndex:       5,
			},
			expected: expectation{
				err: errors.New("rpc error: code = Unknown desc = Invalid RTMR index used."),
			},
		},
	}

	for scenario, tt := range tests {
		t.Run(scenario, func(t *testing.T) {
			out, err := client.GetMeasurement(ctx, tt.in)
			if err != nil {
				if tt.expected.err == nil {
					t.Errorf("Err -> \nWant: nil\nGot: %q\n", err)
				} else {
					if tt.expected.err.Error() != err.Error() {
						t.Errorf("Err -> \nWant: %q\nGot: %q\n", tt.expected.err, err)
					}
				}
			} else {
				if tt.in.MeasurementType != pb.TYPE_SAAS && out.Measurement == "" {
					t.Errorf("Out -> \nWant measurement\nGot : %q\n", out)
				} else {
					_, err = base64.StdEncoding.DecodeString(out.Measurement)
					if err != nil {
						t.Errorf("Out -> \nWant base64 encoded measurement\nGot: %q\n", out.Measurement)
					}
				}
			}

		})
	}
}
