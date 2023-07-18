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
	pkgerrors "github.com/pkg/errors"
	"strings"
)

const (
	TDX_FLAG = "tdx"
	SGX_FLAG = "sgx"
	SEV_FLAG = "sev"
)

var DeviceNotFoundErr = pkgerrors.New("No applicable TDX device found.")

type BaseTeeInterface interface {
	GetType() string
	FindDeviceAvailable() (string, error)
	GetReport(device string, data string) (string, error)
}

type BaseTeeResource struct {
	Type string
}

func NewBaseTeeResource() BaseTeeResource {
	return BaseTeeResource{
		Type: "Base",
	}
}

func (r *BaseTeeResource) GetType() string {
	return r.Type
}

func (r *BaseTeeResource) FindDeviceAvailable() (string, error) {
	tdxResource := NewTdxResource()
	device, err := tdxResource.FindDeviceAvailable()
	if err == nil {
		return device, nil
	}

	return "", DeviceNotFoundErr
}

func (r *BaseTeeResource) GetReport(device string, data string) (string, error) {

	var report string
	var err error

	if strings.Contains(device, TDX_FLAG) {
		tdx := NewTdxResource()
		report, err = tdx.GetReport(device, data)
		if err != nil {
			return "", err
		}
	}

	return report, nil
}
