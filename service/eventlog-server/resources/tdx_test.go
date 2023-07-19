/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */

package resources

import (
	pkgerrors "github.com/pkg/errors"
	"testing"
)

const (
	ERR_MSG = "Exceed valid length"
)

func TestGetUintObjectInvalidLength(t *testing.T) {
	data := make([]byte, 0)
	index := 0
	var err error
	var value1 uint32
	var value2 uint16
	var value3 uint8

	value1, index, err = getUint32Object(data, index)
	if err.Error() != ERR_MSG || value1 != uint32(0) || index != 0 {
		t.Fatalf(`getUint32Object([], 0) = %d, %d, %v want match for %d, %d, %v`,
			value1, index, err, uint32(0), 0, pkgerrors.New("Exceed valid length"))
	}

	value2, index, err = getUint16Object(data, index)
	if err.Error() != ERR_MSG || value2 != uint16(0) || index != 0 {
		t.Fatalf(`getUint16Object([], 0) = %d, %d, %v want match for %d, %d, %v`,
			value2, index, err, uint16(0), 0, pkgerrors.New("Exceed valid length"))
	}

	value3, index, err = getUint8Object(data, index)
	if err.Error() != ERR_MSG || value3 != uint8(0) || index != 0 {
		t.Fatalf(`getUint8Object([], 0) = %d, %d, %v want match for %d, %d, %v`,
			value3, index, err, uint8(0), 0, pkgerrors.New("Exceed valid length"))
	}
}

func TestGetUintObjectValidLength(t *testing.T) {
	data := make([]byte, 8)
	index := 0
	var err error
	var value1 uint32
	var value2 uint16
	var value3 uint8

	value1, index, err = getUint32Object(data, index)
	if err != nil || value1 != uint32(0) || index != 4 {
		t.Fatalf(`getUint32Object([], 0) = %d, %d, %v want match for %d, %d, %v`,
			value1, index, err, uint32(0), 4, nil)
	}

	value2, index, err = getUint16Object(data, index)
	if err != nil || value2 != uint16(0) || index != 6 {
		t.Fatalf(`getUint16Object([], 0) = %d, %d, %v want match for %d, %d, %v`,
			value2, index, err, uint16(0), 6, nil)
	}

	value3, index, err = getUint8Object(data, index)
	if err != nil || value3 != uint8(0) || index != 7 {
		t.Fatalf(`getUint8Object([], 0) = %d, %d, %v want match for %d, %d, %v`,
			value3, index, err, uint8(0), 7, nil)
	}
}

func TestGetEventLogDigestInfo(t *testing.T) {
	var flag bool

	digestCount := uint32(1)
	data := make([]byte, 5)
	index := 0
	digestSizes := make(map[uint16]uint16)
	digestSizes[uint16(0)] = 1
	flag = true

	digests, algId, index, err := getEventLogDigestInfo(data, index, digestCount, digestSizes)
	if len(digests) != 1 {
		flag = false
	}

	for _, item := range digests {
		if item != "0" {
			flag = false
			break
		}
	}

	if err != nil || flag || algId != uint16(0) || index != 3 {
		t.Fatalf(`getEventLogDigestInfo([0,0,0,0,0], 0, 1, map[0]=1) = %s, %d, %d, %v want match for %s, %d, %d, %v`,
			digests, algId, index, err, []string{"0"}, uint16(0), 3, nil)
	}
}

func TestGetHeaderDigestInfo(t *testing.T) {
	data := make([]byte, 8)
	index := 0

	// Success
	data[0] = uint8(1)
	data[6] = uint8(48)

	digestsizes, index, err := getHeaderDigestInfo(data, index)
	if err != nil || digestsizes[0] != uint16(48) || index != 8 {
		t.Fatalf(`getHeaderDigestInfo([1,0,0,0,0,0,48,0], 0) = %d, %d, %v want match for %d, %d, %v`,
			digestsizes[0], index, err, uint16(48), 8, nil)
	}

	// Failure
	data[0] = uint8(0)
	data[1] = uint8(1)

	digestsizes, index, err = getHeaderDigestInfo(data, index)
	if err.Error() != ERR_MSG || digestsizes[0] != uint16(0) || index != 0 {
		t.Fatalf(`getHeaderDigestInfo([0,1,0,0,0,0,48,0], 0) = %d, %d, %v want match for %d, %d, %v`,
			digestsizes[0], index, err, uint16(0), 0, ERR_MSG)
	}
}

func TestGetBasicInfo(t *testing.T) {
	data := make([]byte, 12)
	data[0] = uint8(1)

	rtmr, etype, digestCount, index, err := getBasicInfo(data)
	if err != nil || index != 12 || etype != uint32(0) || digestCount != uint32(0) {
		t.Fatalf(`getBasicInfo([1,0,0,0,0,0,0,0,0,0,0,0]) = %d, %d, %d, %d, %v want %d, %d, %d, %d, %v`,
			rtmr, etype, digestCount, index, err, uint32(0), uint32(0), uint32(0), uint32(0), nil)
	}
}

func TestFetchEventlog(t *testing.T) {
	_, count, err := fetchEventlogs()
	if err != nil || count == 0 {
		t.Fatalf(`fetchEventlog() = %d, %v want %s, %v`, count, err, "large then 0", nil)
	}
}

func TestGetTdxEventlog(t *testing.T) {
	start_position := 0
	count := 1

	log, err := GetTdxEventlog(start_position, count)
	if err != nil || log == "" {
		t.Fatalf(`GetTdxEventlog(0, 1) = %s, %v want %s, %v`, log, err, "eventlogs exist", nil)
	}
}
