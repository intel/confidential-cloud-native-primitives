/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */

package resources

import (
	"testing"
)

func TestFindSEVDeviceAvailable(t *testing.T) {
	r := NewSevResource()

	_, err := r.FindDeviceAvailable()
	if err.Error() != DeviceNotFoundErr.Error() {
		t.Fatalf(`FindDeviceAvailable() = %v want %v`,
			err, DeviceNotFoundErr)
	}
}

func TestGetSEVReport(t *testing.T) {
	r := NewSevResource()
	device := DEVICE_NODE_NAME_1
	data := ""

	_, err := r.GetReport(device, data)
	if err != nil {
		t.Fatalf(`GetReport() = %v want %v`,
			err, nil)
	}
}
