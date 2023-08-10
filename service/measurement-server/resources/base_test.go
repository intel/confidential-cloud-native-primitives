/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */

package resources

import (
	"encoding/base64"
	"golang.org/x/exp/slices"
	"testing"
)

var (
	TDX_DEVICE_LIST = []string{DEVICE_NODE_NAME_1_0, DEVICE_NODE_NAME_1_5}
)

func TestFindDeviceAvailable(t *testing.T) {
	r := NewBaseTeeResource()

	device, err := r.FindDeviceAvailable()

	if err != nil || !slices.Contains(TDX_DEVICE_LIST, device) {
		t.Fatalf(`FindDeviceAvailable() = %s, %v want TDX device, %v`, device, err, nil)
	}
}

func TestGetReport(t *testing.T) {
	r := NewBaseTeeResource()

	data := "test"
	device, _ := r.FindDeviceAvailable()

	report, err := r.GetReport(device, data)
	if err != nil {
		t.Fatalf(`GetReport(device, data) = %s, %v want report and %v`,
			report, err, nil)
	}
	_, err = base64.StdEncoding.DecodeString(report)
	if err != nil {
		t.Fatalf(`GetReport(device, data) = %s, %v want base64 encoded report and %v`,
			report, err, nil)
	}
}
