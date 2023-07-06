/*
*
* Copyright 2023 Intel authors.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*
 */

package resources

import (
	"log"
	"os"
)

const (
	// The device fd for AMD SEV
	DEVICE_NODE_NAME_1 = "/dev/sev-guest"
	DEVICE_NODE_NAME_2 = "/dev/sev"
)

type SevResource struct {
	BaseTeeResource
}

func NewSevResource() *SevResource {
	return &SevResource{
		BaseTeeResource{
			Type: "AMD SEV",
		},
	}
}

func (r *SevResource) FindDeviceAvailable() (string, error) {

	if _, err := os.Stat(DEVICE_NODE_NAME_1); err == nil {
		return DEVICE_NODE_NAME_1, nil
	}

	if _, err := os.Stat(DEVICE_NODE_NAME_2); err == nil {
		return DEVICE_NODE_NAME_2, nil
	}

	return "", DeviceNotFoundErr
}

func (r *SevResource) GetReport(device string, data string) (string, error) {

	/*
	   // Open TDX device fd to get prepared for TDVM call
	   deviceNode, err := os.OpenFile(device, os.O_RDWR, 0644)
	   if err != nil {
	       return "", err
	   }

	   // TODO: check if there can be a data attached to SEV report
	   if len(data) > REPORT_DATA_LEN {
	       err = pkgerrors.New("Report data with invalid length.")
	       return "", err
	   }*/
	log.Println("Not implemented")
	return "", nil
}
