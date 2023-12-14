/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */

package eventlog

import (
	"testing"

	pb "github.com/intel/confidential-cloud-native-primitives/sdk/golang/ccnp/eventlog/proto"
)

func TestGetPlatformEventlogDefault(t *testing.T) {
	eventlogs, err := GetPlatformEventlog()

	if err != nil {
		t.Fatalf("[TestGetPlatformEventlogDefault] get Platform Eventlog error: %v", err)
	}

	if len(eventlogs) == 0 {
		t.Fatalf("[TestGetPlatformEventlogDefault] error: no eventlog returns")
	}

}

func TestGetPlatformEventlogWithEventlogCategory(t *testing.T) {

	eventlogs, err := GetPlatformEventlog(WithEventlogCategory(pb.CATEGORY_TDX_EVENTLOG))

	if err != nil {
		t.Fatalf("[TestGetPlatformEventlogWithEventlogCategory] get Platform Eventlog error: %v", err)
	}

	if len(eventlogs) == 0 {
		t.Fatalf("[TestGetPlatformEventlogWithEventlogCategory] error: no eventlog returns")
	}

}

func TestGetPlatformEventlogWithStartPosition(t *testing.T) {

	eventlogs, err := GetPlatformEventlog(WithStartPosition(2))

	if err != nil {
		t.Fatalf("[TestGetPlatformEventlogWithEventlogCategory] get Platform Eventlog error: %v", err)
	}

	if len(eventlogs) == 0 {
		t.Fatalf("[TestGetPlatformEventlogWithEventlogCategory] error: no eventlog returns")
	}

}

func TestGetPlatformEventlogWithStartPositionAndCount(t *testing.T) {

	eventlogs, err := GetPlatformEventlog(WithStartPosition(2), WithCount(5))

	if err != nil {
		t.Fatalf("[TestGetPlatformEventlogWithStartPositionAndCount] get Platform Eventlog error: %v", err)
	}

	if len(eventlogs) != 5 {
		t.Fatalf("[TestGetPlatformEventlogWithStartPositionAndCount] error: expected number of logs is 5, retrieved %v", len(eventlogs))
	}

}

func TestGetPlatformEventlogWithAllOptions(t *testing.T) {

	eventlogs, err := GetPlatformEventlog(WithEventlogCategory(pb.CATEGORY_TDX_EVENTLOG), WithStartPosition(2), WithCount(3))

	if err != nil {
		t.Fatalf("[TestGetPlatformEventlogWithAllOptions] get Platform Eventlog error: %v", err)
	}

	if len(eventlogs) != 3 {
		t.Fatalf("[TestGetPlatformEventlogWithAllOptions] error: expected number of logs is 3, retrieved %v", len(eventlogs))
	}

}
