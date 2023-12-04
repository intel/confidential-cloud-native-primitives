/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */

package measurement

import (
	"testing"

	pb "github.com/intel/confidential-cloud-native-primitives/sdk/golang/ccnp/measurement/proto"
)

const (
	EXPECTED_REPORT_DATA       = "abcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefghabcdefgh"
	EXPECTED_REPORT_DATA_SHORT = "abcd"
	TDREPORT_TYPE_LENGTH       = 4
	TDREPORT_TYPE_BYTE1        = 129
	CATEGORY_UNKNOWN           = 3
	TDX_RTMR_INDEX_UNKNOWN     = 4
	EXPECTED_TDX_REPORT_LEN    = 1024
	TEE_TYPE_TDX               = 129
	TDX_TCB_SVN_LENGTH         = 16
	TDX_MRSEAM_LENGTH          = 48
	TDX_MRSEAMSINGER_LENGTH    = 48
	TDX_SEAM_ATTRIBUTES_LENGTH = 8
	TDX_TD_ATTRIBUTES_LENGTH   = 8
	TDX_XFAM_LENGTH            = 8
	TDX_MRTD_LENGTH            = 48
	TDX_MRCONFIGID_LENGTH      = 48
	TDX_MROWNER_LENGTH         = 48
	TDX_MROWNERCONFIG_LENGTH   = 48
	TDX_RTMR_LENGTH            = 48
	TDX_RTMRS_LENGTH           = 192
	TDX_REPORT_DATA_LENGTH     = 64
)

func parseTDXReportAndEvaluate(r TDReportInfo, reportData string, t *testing.T) {
	if len(r.TDReportRaw) != EXPECTED_TDX_REPORT_LEN {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport size, retrieved: %v, expected: %v", len(r.TDReportRaw), EXPECTED_TDX_REPORT_LEN)
	}

	tdreport := r.TDReport

	if len(tdreport.ReportType) != TDREPORT_TYPE_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport Report Type length, retrieved: %v, expected: %v", len(tdreport.ReportType), TDREPORT_TYPE_LENGTH)
	}

	if tdreport.ReportType[0] != TDREPORT_TYPE_BYTE1 {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport ReportType[0], retrieved: %v, expected: %v", tdreport.ReportType[0], TDREPORT_TYPE_BYTE1)
	}

	if len(tdreport.TeeTcbSvn) != TDX_TCB_SVN_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport TEE TCB SVN length, retrieved: %v, expected: %v", len(tdreport.TeeTcbSvn), TDX_TCB_SVN_LENGTH)
	}

	if len(tdreport.Mrseam) != TDX_MRSEAM_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport Mrseam length, retrieved: %v, expected: %v", len(tdreport.Mrseam), TDX_MRSEAM_LENGTH)
	}

	if len(tdreport.Mrseamsigner) != TDX_MRSEAMSINGER_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport Mrseamsigner length, retrieved: %v, expected: %v", len(tdreport.Mrseamsigner), TDX_MRSEAMSINGER_LENGTH)
	}

	if len(tdreport.SeamAttributes) != TDX_SEAM_ATTRIBUTES_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport SeamAttributes length, retrieved: %v, expected: %v", len(tdreport.SeamAttributes), TDX_SEAM_ATTRIBUTES_LENGTH)
	}

	if len(tdreport.TdAttributes) != TDX_TD_ATTRIBUTES_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport TdAttributes length, retrieved: %v, expected: %v", len(tdreport.TdAttributes), TDX_TD_ATTRIBUTES_LENGTH)
	}

	if len(tdreport.Xfam) != TDX_XFAM_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport Xfam length, retrieved: %v, expected: %v", len(tdreport.Xfam), TDX_XFAM_LENGTH)
	}

	if len(tdreport.Mrtd) != TDX_MRTD_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport Mrtd length, retrieved: %v, expected: %v", len(tdreport.Mrtd), TDX_MRTD_LENGTH)
	}

	if len(tdreport.Mrconfigid) != TDX_MRCONFIGID_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport Mrconfigid length, retrieved: %v, expected: %v", len(tdreport.Mrconfigid), TDX_MRCONFIGID_LENGTH)
	}

	if len(tdreport.Mrowner) != TDX_MROWNER_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport Mrowner length, retrieved: %v, expected: %v", len(tdreport.Mrowner), TDX_MROWNER_LENGTH)
	}

	if len(tdreport.Mrownerconfig) != TDX_MROWNERCONFIG_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport Mrownerconfig length, retrieved: %v, expected: %v", len(tdreport.Mrownerconfig), TDX_MROWNERCONFIG_LENGTH)
	}

	if len(tdreport.Rtmrs) != TDX_RTMRS_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport Rtmrs length, retrieved: %v, expected: %v", len(tdreport.Rtmrs), TDX_RTMRS_LENGTH)
	}

	if len(tdreport.ReportData) != TDX_REPORT_DATA_LENGTH {
		t.Fatalf("[parseTDXReportAndEvaluate] wrong TDReport Data length, retrieved: %v, expected: %v", len(tdreport.ReportData), TDX_REPORT_DATA_LENGTH)
	}

	if len(reportData) != 0 {
		if string(tdreport.ReportData[:len(reportData)]) != reportData {
			t.Fatalf("[parseTDXReportAndEvaluate], report data retrieved = %s, expected %s",
				tdreport.ReportData, reportData)
		}
	} else {
		var empty_report_data [64]uint8
		if tdreport.ReportData != empty_report_data {
			t.Fatalf("[parseTDXReportAndEvaluate], report data retrieved = %v, expected empty string",
				tdreport.ReportData)
		}
	}
}

func parseTDXRtmrAndEvaluate(r TDXRtmrInfo, t *testing.T) {
	if len(r.TDXRtmrRaw) != TDX_RTMR_LENGTH {
		t.Fatalf("[parseTDXRtmrAndEvaluate] wrong RTMT size, retrieved: %v, expected: %v", len(r.TDXRtmrRaw), TDX_RTMR_LENGTH)
	}
}

func TestGetPlatformMeasurementTDReportDefault(t *testing.T) {
	ret, err := GetPlatformMeasurement()
	if err != nil {
		t.Fatalf("[TestGetPlatformMeasurementTDReportDefault] get Platform Measurement error: %v", err)
	}

	switch ret.(type) {
	case TDReportInfo:
		var r, _ = ret.(TDReportInfo)
		parseTDXReportAndEvaluate(r, "", t)
	default:
		t.Fatalf("[TestGetPlatformMeasurementTDReportDefault] unknown TEE enviroment!")
	}
}

func TestGetPlatformMeasurementTDReportCategoryOnly(t *testing.T) {
	ret, err := GetPlatformMeasurement(WithMeasurementType(pb.CATEGORY_TEE_REPORT))
	if err != nil {
		t.Fatalf("[TestGetPlatformMeasurementTDReportCategoryOnly] get Platform Measurement error: %v", err)
	}

	switch ret.(type) {
	case TDReportInfo:
		var r, _ = ret.(TDReportInfo)
		parseTDXReportAndEvaluate(r, "", t)

	default:
		t.Fatalf("[TestGetPlatformMeasurementTDReportCategoryOnly] unknown TEE enviroment!")
	}
}

func TestGetPlatformMeasurementTDReportCategoryAndEmptyReportData(t *testing.T) {
	ret, err := GetPlatformMeasurement(WithMeasurementType(pb.CATEGORY_TEE_REPORT), WithReportData(""))
	if err != nil {
		t.Fatalf("[TestGetPlatformMeasurementTDReportCategoryAndReportData] get Platform Measurement error: %v", err)
	}

	switch ret.(type) {
	case TDReportInfo:
		var r, _ = ret.(TDReportInfo)
		parseTDXReportAndEvaluate(r, "", t)

	default:
		t.Fatalf("[TestGetPlatformMeasurementTDReportCategoryAndReportData] unknown TEE enviroment!")
	}
}

func TestGetPlatformMeasurementTDReportCategoryAndReportData(t *testing.T) {
	ret, err := GetPlatformMeasurement(WithMeasurementType(pb.CATEGORY_TEE_REPORT), WithReportData(EXPECTED_REPORT_DATA))
	if err != nil {
		t.Fatalf("[TestGetPlatformMeasurementTDReportCategoryAndReportData] get Platform Measurement error: %v", err)
	}

	switch ret.(type) {
	case TDReportInfo:
		var r, _ = ret.(TDReportInfo)
		parseTDXReportAndEvaluate(r, EXPECTED_REPORT_DATA, t)

	default:
		t.Fatalf("[TestGetPlatformMeasurementTDReportCategoryAndReportData] unknown TEE enviroment!")
	}
}

func TestGetPlatformMeasurementTDReportCategoryAndShortReportData(t *testing.T) {
	ret, err := GetPlatformMeasurement(WithMeasurementType(pb.CATEGORY_TEE_REPORT), WithReportData(EXPECTED_REPORT_DATA_SHORT))
	if err != nil {
		t.Fatalf("[TestGetPlatformMeasurementTDReportCategoryAndShortReportData] get Platform Measurement error: %v", err)
	}

	switch ret.(type) {
	case TDReportInfo:
		var r, _ = ret.(TDReportInfo)
		parseTDXReportAndEvaluate(r, EXPECTED_REPORT_DATA_SHORT, t)

	default:
		t.Fatalf("[TestGetPlatformMeasurementTDReportCategoryAndShortReportData] unknown TEE enviroment!")
	}
}

func TestGetPlatformMeasurementRTMRWithMeasurementType(t *testing.T) {
	ret, err := GetPlatformMeasurement(WithMeasurementType(pb.CATEGORY_TDX_RTMR))
	if err != nil {
		t.Fatalf("[TestGetPlatformMeasurementRTMRWithMeasurementType] get Platform Measurement error: %v", err)
	}

	switch ret.(type) {
	case TDXRtmrInfo:
		var r, _ = ret.(TDXRtmrInfo)
		parseTDXRtmrAndEvaluate(r, t)

	default:
		t.Fatalf("[TestGetPlatformMeasurementRTMRWithMeasurementType] unknown TEE enviroment!")
	}
}

func TestGetPlatformMeasurementRTMRWithMeasurementTypeAndIndex(t *testing.T) {
	ret, err := GetPlatformMeasurement(WithMeasurementType(pb.CATEGORY_TDX_RTMR), WithRegisterIndex(1))
	if err != nil {
		t.Fatalf("[TestGetPlatformMeasurementRTMRWithMeasurementTypeAndIndex] get Platform Measurement error: %v", err)
	}

	switch ret.(type) {
	case TDXRtmrInfo:
		var r, _ = ret.(TDXRtmrInfo)
		parseTDXRtmrAndEvaluate(r, t)

	default:
		t.Fatalf("[TestGetPlatformMeasurementRTMRWithMeasurementTypeAndIndex] unknown TEE enviroment!")
	}
}
