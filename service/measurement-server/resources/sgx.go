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
	DEVICE_NODE_NAME = "/dev/sgx"
)

type SgxResource struct {
	BaseTeeResource
}

func NewSgxResource() *SgxResource {
	return &SgxResource{
		BaseTeeResource{
			Type: "Intel SGX",
		},
	}
}

func (r *SgxResource) FindDeviceAvailable() (string, error) {

	if _, err := os.Stat(DEVICE_NODE_NAME); err == nil {
		return DEVICE_NODE_NAME, nil
	}

	return "", DeviceNotFoundErr
}

func (r *SgxResource) GetReport(device string, data string) (string, error) {

	/*
	   // Open TDX device fd to get prepared for TDVM call
	   deviceNode, err := os.OpenFile(device, os.O_RDWR, 0644)
	   if err != nil {
	       return "", err
	   }

	   if len(data) > REPORT_DATA_LEN {
	       err = pkgerrors.New("Report data with invalid length.")
	       return "", err
	   }
	*/
	log.Println("Not implemented")
	return "", nil
}
