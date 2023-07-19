/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */

package resources

import (
	"testing"
)

const (
	TPM_ERR_MSG = "stat /sys/kernel/security/tpm0/binary_bios_measurements: no such file or directory"
)

func TestGetTpmEventlog(t *testing.T) {
	start_position := 0
	count := 1

	_, err := GetTpmEventlog(start_position, count)
	if err.Error() != TPM_ERR_MSG {
		t.Fatalf(`GetTpmEventlog(0,1) get error %s want %s`, err.Error(), TPM_ERR_MSG)
	}
}
