/* SPDX-license-identifier: Apache-2.0*/

package resources

import (
	pkgerrors "github.com/pkg/errors"
	"strings"
)

const (
	TDX_FLAG = "tdx"
	SEV_FLAG = "sev"
)

var DeviceNotFoundErr = pkgerrors.New("No applicable device found.")

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

	sevResource := NewSevResource()
	device, err = sevResource.FindDeviceAvailable()
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
	} else if strings.Contains(device, SEV_FLAG) {
		sev := NewSevResource()
		report, err = sev.GetReport(device, data)
		if err != nil {
			return "", err
		}
	}

	return report, nil
}
