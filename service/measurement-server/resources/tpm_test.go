/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */

package resources

import (
	"testing"
)

func TestFindTPMDeviceAvailable(t *testing.T) {

	_, err := findDeviceAvailable()
	if err.Error() != DeviceNotFoundErr.Error() {
		t.Fatalf(`FindDeviceAvailable() = %v want %v`,
			err, DeviceNotFoundErr)
	}
}

func TestGetTpmMeasurement(t *testing.T) {
	index := 0

	_, err := GetTpmMeasurement(index)
	if err.Error() != DeviceNotFoundErr.Error() {
		t.Fatalf(`GetReport() = %v want %v`,
			err, DeviceNotFoundErr)
	}
}
