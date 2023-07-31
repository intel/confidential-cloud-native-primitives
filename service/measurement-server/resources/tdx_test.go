/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */

package resources

import (
	"encoding/base64"
	pkgerrors "github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"testing"
)

func TestFindTDXDeviceAvailable(t *testing.T) {
	r := NewTdxResource()

	device, err := r.FindDeviceAvailable()
	if err != nil || !slices.Contains(TDX_DEVICE_LIST, device) {
		t.Fatalf(`FindDeviceAvailable() = %s, %v want valid tdx device and nil`,
			device, err)
	}
}

func TestTDXGetReport(t *testing.T) {
	r := NewTdxResource()

	device, _ := r.FindDeviceAvailable()
	report, err := r.GetReport(device, "")
	if err != nil {
		t.Fatalf(`GetReport(device, "") = %s, %v want report and nil`,
			report, err)
	}
	_, err = base64.StdEncoding.DecodeString(report)
	if err != nil {
		t.Fatalf(`GetReport(device, "") = %s, %v want base64 encoded report and nil`,
			report, err)
	}
}

func TestGetRTMRMeasurement(t *testing.T) {
	r := NewTdxResource()

	device, _ := r.FindDeviceAvailable()
	index := 0

	measurement, err := r.GetRTMRMeasurement(device, "", index)
	if err != nil {
		t.Fatalf(`GetRTMRMeasurement(device, "", 0) = %s, %v want measurement, %v`,
			measurement, err, nil)
	}
	_, err = base64.StdEncoding.DecodeString(measurement)
	if err != nil {
		t.Fatalf(`GetRTMRMeasurement(device, "", 0) = %s, %v want base64 encoded measurement, %v`,
			measurement, err, nil)
	}
}

func TestCollectRtmrMeasurement(t *testing.T) {
	report := "gQAAAAAAAAAAAAAAAAAAAAEBAQEB/wABAAAAAAAAAAA2OBU1hr5liB7Hcj2fQyCokisfxMHJ1OaXAQCQ/DTfh/IWWhhEH9TLXg4wlHIYCs8Ku50kbzi8ZyJBxnA+PZyuM3ulswD3Zxhhzt9AelPfSZX06h7XRu2m1aOGrtD3tQ4AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADrxjc/2ZBFYiP3mNWUILLmoXf+lHqdGIMjj5Sv/hn8C/8BAwAAAAAAAAEBAAAAAAAAAAAAAAAAAFi1VbaJLemWgQThKktgTVRGjKyORNj10YBYB8YItDdufnvvDf5alim7S2hZcvwDIgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEBAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAQAAAAAOcaBgAAAAAApKADNGxaGab9JQRx6HK9Bx2MktdDGr2kY0F4CKFzg6oNQph4FLyS9fWcYES2d/UUAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAnxfdBYPsEXhJQLRMcCfs30lC9XhkJEd1WAGCYlRFWrk9kOE8AD1XDT557A451xHsbR8pW9kYqecaRNV4+V4HjBX5yzX1YJhMUx5ErBFMTU4HF4Bsk9Gk3raULoSQFNCt6NfjSQvWDZug3rNW8HHa1PEUTfNH1mVexFvyDJ7rFpZ5MHrjg8aslqTue2ojc/uAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="
	index := 0
	sample_measurement, _ := base64.StdEncoding.DecodeString("nxfdBYPsEXhJQLRMcCfs30lC9XhkJEd1WAGCYlRFWrk9kOE8AD1XDT557A451xHs")

	measurement, err := collectRtmrMeasurement(report, index)
	encoded_measurement := base64.StdEncoding.EncodeToString(measurement)

	if err != nil || encoded_measurement != "nxfdBYPsEXhJQLRMcCfs30lC9XhkJEd1WAGCYlRFWrk9kOE8AD1XDT557A451xHs" {
		t.Fatalf(`collectRtmrMeasurement(report, 0) = %v, %v want %v, %v`,
			measurement, err, sample_measurement, nil)
	}
}

func TestGetRTMRMeasurementWithInvalidIndex(t *testing.T) {
	r := NewTdxResource()

	device, _ := r.FindDeviceAvailable()
	index := 5

	_, err := r.GetRTMRMeasurement(device, "", index)
	if err.Error() != "Invalid RTMR index used." {
		t.Fatalf(`GetRTMRMeasurement(device, "", 5) = %v want, %v`,
			err, pkgerrors.New("Invalid RTMR index used."))
	}
}
