/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
*/

package resources

import (
	"encoding/base64"
	"log"
	"os"
	"syscall"
	"unsafe"

	pkgerrors "github.com/pkg/errors"
)

const (

	// The name of the device in different kernel version
	DEVICE_NODE_NAME_DEPRECATED = "/dev/tdx-attest"
	DEVICE_NODE_NAME_1_0        = "/dev/tdx-guest"
	DEVICE_NODE_NAME_1_5        = "/dev/tdx_guest"

	// The length of report data
	REPORT_DATA_LEN = 64
	// The length of TDX report
	TDX_REPORT_LEN = 1024
	// The length of Data can be extended into RTMR
	TDX_EXTEND_RTMR_DATA_LEN = 48

	/*The device operators for tdx v1.0
	  Reference: TDX_CMD_GET_REPORT = _IOWR('T', 0x01, __u64)
	  defined in arch/x86/include/uapi/asm/tdx.h in kernel source
	*/
	TDX_CMD_GET_REPORT_V1_0 = 0xc0085401

	/* The device operators for tdx v1.5
	   Reference: TDX_CMD_GET_REPORT0 = _IOWR('T', 1, struct tdx_report_req)
	   defined in include/uapi/linux/tdx-guest.h in kernel source
	*/
	TDX_CMD_GET_REPORT0_V1_5 = 0xc4405401

	RTMR_0_OFFSET = 0x2d0
	RTMR_1_OFFSET = 0x300
	RTMR_2_OFFSET = 0x330
	RTMR_3_OFFSET = 0x360

	RTMR_LEN = 0x30
)

var TdxGetReportErr = pkgerrors.New("Failed to get TDX report.")
var InvalidRtmrIndexErr = pkgerrors.New("Invalid RTMR index used.")

type TdxReportReq struct {
	SubType    uint8
	ReportData uint64
	RpdLen     uint32
	TdReport   uint64
	TdrLen     uint32
}

type TdxResource struct {
	BaseTeeResource
}

func NewTdxResource() *TdxResource {
	return &TdxResource{
		BaseTeeResource{
			Type: "Intel TDX",
		},
	}
}

func NewTdxReportReq(data string) TdxReportReq {
	d := make([]byte, REPORT_DATA_LEN)
	copy(d, []byte(data))

	r := make([]byte, TDX_REPORT_LEN)

	return TdxReportReq{
		SubType:    uint8(0),
		ReportData: uint64(uintptr(unsafe.Pointer(&d[0]))),
		RpdLen:     uint32(REPORT_DATA_LEN),
		TdReport:   uint64(uintptr(unsafe.Pointer(&r[0]))),
		TdrLen:     uint32(TDX_REPORT_LEN),
	}
}

func NewTdxReportReq0(data string) []byte {
	d := make([]byte, REPORT_DATA_LEN+TDX_REPORT_LEN)
	copy(d, []byte(data))
	return d
}

func (r *TdxResource) FindDeviceAvailable() (string, error) {

	if _, err := os.Stat(DEVICE_NODE_NAME_DEPRECATED); err == nil {
		log.Printf("Deprecated device node %s, please upgrade to use %s or %s",
			DEVICE_NODE_NAME_DEPRECATED, DEVICE_NODE_NAME_1_0, DEVICE_NODE_NAME_1_5)
		return "", DeviceNotFoundErr
	}

	if _, err := os.Stat(DEVICE_NODE_NAME_1_0); err == nil {
		return DEVICE_NODE_NAME_1_0, nil
	}

	if _, err := os.Stat(DEVICE_NODE_NAME_1_5); err == nil {
		return DEVICE_NODE_NAME_1_5, nil
	}

	return "", DeviceNotFoundErr
}

func (r *TdxResource) GetReport(device string, data string) (string, error) {

	var report string

	/* Open TDX device fd to get prepared for TDVM call*/
	deviceNode, err := os.OpenFile(device, os.O_RDWR, 0644)
	if err != nil {
		return "", err
	}

	if len(data) > REPORT_DATA_LEN {
		err = pkgerrors.New("Report data with invalid length.")
		return "", err
	}

	if device == DEVICE_NODE_NAME_1_0 {
		report, err = getTdxReport(deviceNode, data)
		if err != nil {
			return "", err
		}
	} else {
		report, err = getTdxReport0(deviceNode, data)
		if err != nil {
			return "", err
		}
	}

	return report, nil
}

func getTdxReport(deviceNode *os.File, data string) (string, error) {

	req := NewTdxReportReq(data)
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(deviceNode.Fd()),
		uintptr(TDX_CMD_GET_REPORT_V1_0), uintptr(unsafe.Pointer(&req)))
	if errno != 0 {
		return "", TdxGetReportErr
	}

	td_report_value := make([]byte, TDX_REPORT_LEN)
	copy(td_report_value, (*[1 << 30]byte)(unsafe.Pointer(uintptr(req.TdReport)))[:TDX_REPORT_LEN])
	report_string := base64.StdEncoding.EncodeToString(td_report_value)

	return report_string, nil
}

func getTdxReport0(deviceNode *os.File, data string) (string, error) {

	req := NewTdxReportReq0(data)
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(deviceNode.Fd()),
		uintptr(TDX_CMD_GET_REPORT0_V1_5), uintptr(unsafe.Pointer(&req[0])))
	if errno != 0 {
		return "", TdxGetReportErr
	}

	report := base64.StdEncoding.EncodeToString(req[REPORT_DATA_LEN:])
	return report, nil
}

func (r *TdxResource) GetRTMRMeasurement(device string, data string, index int) (string, error) {

	if index < 0 || index > 3 {
		return "", InvalidRtmrIndexErr
	}

	report, err := r.GetReport(device, data)
	if err != nil {
		return "", err
	}

	measurement, err := collectRtmrMeasurement(report, index)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(measurement), nil
}

func collectRtmrMeasurement(report string, index int) ([]byte, error) {

	r, err := base64.StdEncoding.DecodeString(report)
	if err != nil {
		return []byte{}, err
	}
	value := make([]byte, RTMR_LEN)

	switch index {
	case 0:
		value = r[RTMR_0_OFFSET : RTMR_0_OFFSET+RTMR_LEN]
	case 1:
		value = r[RTMR_1_OFFSET : RTMR_1_OFFSET+RTMR_LEN]
	case 2:
		value = r[RTMR_2_OFFSET : RTMR_2_OFFSET+RTMR_LEN]
	case 3:
		value = r[RTMR_3_OFFSET : RTMR_3_OFFSET+RTMR_LEN]
	}
	return value, nil
}
