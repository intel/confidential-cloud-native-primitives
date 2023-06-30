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
    "io"

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
