/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */

package resources

import (
	"io"
	"log"
	"os"

	pkgerrors "github.com/pkg/errors"
)

const (
	//The location of CCEL table
	TPM_EVENT_LOG_LOCATION = "/sys/kernel/security/tpm0/binary_bios_measurements"

	CHUNK_SIZE = 16384
)

var (
	TpmGetEventlogErr      = pkgerrors.New("Failed to get eventlog in TPM.")
	TpmEventlogNotFoundErr = pkgerrors.New("TPM eventlog not found.")
)

func GetTpmEventlog(start_position int, length int) (string, error) {

	var eventlog string

	/* Check if the tpm eventlog file exists*/
	if _, err := os.Stat(TPM_EVENT_LOG_LOCATION); err != nil {
		return "", err
	}

	/* Read eventlog data*/
	object, err := os.OpenFile(TPM_EVENT_LOG_LOCATION, os.O_RDONLY, 0644)
	if err != nil {
		return "", TpmEventlogNotFoundErr
	}

	data, err := io.ReadAll(object)
	if err != nil {
		return "", err
	}

	if len(data) == 0 {
		return "", TpmEventlogNotFoundErr
	}

	eventlog, err = parseTpmEventlog(data, start_position, length)
	if err != nil {
		return "", err
	}

	return eventlog, nil
}

func parseTpmEventlog(data []byte, position int, length int) (string, error) {
	//not implemented, refer to https://github.com/tpm2-software/tpm2-tools/blob/master/tools/misc/tpm2_eventlog.c#L62
	log.Println("Not implemented.")

	return "", nil

}
