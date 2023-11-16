/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */

package measurement

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"log"
	"time"

	pb "github.com/intel/confidential-cloud-native-primitives/sdk/golang/ccnp/measurement/proto"
	pkgerrors "github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	UDS_PATH       = "unix:/run/ccnp/uds/measurement.sock"
	TDX_REPORT_LEN = 1024
)

type GetPlatformMeasurementOptions struct {
	measurementType pb.CATEGORY
	reportData      string
	registerIndex   int32
}

type TDReportInfo struct {
	TDReportRaw [TDX_REPORT_LEN]uint8 // full TD report
	TDReport    TDReportStruct
}

type TDReportStruct struct {
	//REPORTMACSTRUCT
	ReportType     [4]uint8
	Reserved1      [12]uint8
	CpuSvn         [16]uint8
	TeeTcbInfoHash [48]uint8
	TeeInfoHash    [48]uint8
	ReportData     [64]uint8
	Reserved2      [32]uint8
	Mac            [32]uint8

	//TEE_TCB_INFO
	Mrseam         [48]uint8
	Mrseamsigner   [48]uint8
	TeeTcbSvn      [16]uint8
	SeamAttributes [8]uint8

	//RESERVED
	Reserved3 [17]uint8

	//TDINFO_STRUCT
	TdAttributes  [8]uint8
	Xfam          [8]uint8
	Mrtd          [48]uint8
	Mrconfigid    [48]uint8
	Mrowner       [48]uint8
	Mrownerconfig [48]uint8
	Rtmrs         [192]uint8
	Reserved4     [112]uint8
}

type TDXRtmrInfo struct {
	TDXRtmrRaw []uint8
}

type TPMReportInfo struct {
	TPMReportRaw []uint8
	TPMReport    TPMReportStruct
}

type TPMReportStruct struct{}

func checkMeasurementType(measurementType pb.CATEGORY) bool {
	return measurementType == pb.CATEGORY_TEE_REPORT || measurementType == pb.CATEGORY_TDX_RTMR || measurementType == pb.CATEGORY_TPM
}

func WithMeasurementType(measurementType pb.CATEGORY) func(*GetPlatformMeasurementOptions) {
	return func(opts *GetPlatformMeasurementOptions) {
		opts.measurementType = measurementType
	}
}

func WithReportData(reportData string) func(*GetPlatformMeasurementOptions) {
	return func(opts *GetPlatformMeasurementOptions) {
		opts.reportData = reportData
	}
}

func WithRegisterIndex(registerIndex int32) func(*GetPlatformMeasurementOptions) {
	return func(opts *GetPlatformMeasurementOptions) {
		opts.registerIndex = registerIndex
	}
}

func GetPlatformMeasurement(opts ...func(*GetPlatformMeasurementOptions)) (interface{}, error) {
	input := GetPlatformMeasurementOptions{measurementType: pb.CATEGORY_TEE_REPORT, reportData: "", registerIndex: 0}
	for _, opt := range opts {
		opt(&input)
	}

	if !checkMeasurementType(input.measurementType) {
		log.Fatalf("[GetPlatformMeasurement] Invalid measurementType specified")
	}

	if input.measurementType == pb.CATEGORY_TPM {
		log.Fatalf("[GetPlatformMeasurement] TPM to be supported later")
	}

	if len(input.reportData) > 64 {
		log.Fatalf("[GetPlatformMeasurement] Invalid reportData specified")
	}

	if input.registerIndex < 0 || input.registerIndex > 16 {
		log.Fatalf("[GetPlatformMeasurement] Invalid registerIndex specified")
	}

	channel, err := grpc.Dial(UDS_PATH, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[GetPlatformMeasurement] can not connect to UDS: %v", err)
	}
	defer channel.Close()

	client := pb.NewMeasurementClient(channel)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.GetMeasurement(ctx, &pb.GetMeasurementRequest{
		MeasurementType:     pb.TYPE_PAAS,
		MeasurementCategory: input.measurementType,
		ReportData:          input.reportData,
		RegisterIndex:       input.registerIndex,
	})

	if err != nil {
		log.Fatalf("[GetPlatformMeasurement] fail to get Platform Measurement: %v", err)
	}

	measurement, err := base64.StdEncoding.DecodeString(response.Measurement)
	if err != nil {
		log.Fatalf("[GetPlatformMeasurement] decode tdreport error: %v", err)
	}

	switch input.measurementType {
	case pb.CATEGORY_TEE_REPORT:
		//TODO: need to get the type of TEE: TDX, SEV, vTPM etc.
		var tdReportInfo = TDReportInfo{}
		err = binary.Read(bytes.NewReader(measurement[0:TDX_REPORT_LEN]), binary.LittleEndian, &tdReportInfo.TDReportRaw)
		if err != nil {
			log.Fatalf("[parseTDXQuote] fail to parse quote cert data: %v", err)
		}
		tdReportInfo.TDReport = parseTDXReport(measurement)
		return tdReportInfo, nil
	case pb.CATEGORY_TDX_RTMR:
		var tdxRtmrInfo = TDXRtmrInfo{}
		tdxRtmrInfo.TDXRtmrRaw = measurement
		return tdxRtmrInfo, nil
	case pb.CATEGORY_TPM:
		return "", pkgerrors.New("[GetPlatformMeasurement] TPM to be supported later")
	default:
		log.Fatalf("[GetPlatformMeasurement] unknown TEE enviroment!")
	}

	return "", pkgerrors.New("[GetPlatformMeasurement] unknown TEE enviroment!")
}

func parseTDXReport(report []byte) TDReportStruct {
	var tdreport = TDReportStruct{}
	err := binary.Read(bytes.NewReader(report[0:TDX_REPORT_LEN]), binary.LittleEndian, &tdreport)
	if err != nil {
		log.Fatalf("[parseTDXReport] fail to parse tdreport: %v", err)
	}

	return tdreport
}

func parseTPMReport(report []byte) (interface{}, error) {
	return nil, pkgerrors.New("TPM to be supported later.")
}

func GetContainerMeasurement() (interface{}, error) {
	// TODO: add Container Measurement support later
	return nil, pkgerrors.New("Container Measurement support to be implemented later.")
}
