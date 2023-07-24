/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
*/

package resources

import (
	"log"
	"os"
)

const (
	// The device fd for AMD SEV
	DEVICE_NODE_NAME_TPM = "/dev/tpm0"
)

func findDeviceAvailable() (string, error) {

	if _, err := os.Stat(DEVICE_NODE_NAME_TPM); err == nil {
		return DEVICE_NODE_NAME_TPM, nil
	}

	return "", DeviceNotFoundErr
}

func GetTpmMeasurement(index int) (string, error) {

	_, err := findDeviceAvailable()
	if err != nil {
		return "", err
	}

	/*
	   // Open TPM device fd to get prepared
	   deviceNode, err := os.OpenFile(device, os.O_RDWR, 0644)
	   if err != nil {
	       return "", err
	   }
	*/
	log.Println("Not implemented")
	return "", nil

}
