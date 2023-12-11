/*
* Copyright (c) 2023, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */

package eventlog

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	pb "github.com/intel/confidential-cloud-native-primitives/sdk/golang/ccnp/eventlog/proto"
	el "github.com/intel/confidential-cloud-native-primitives/service/eventlog-server/resources"
	pkgerrors "github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	UDS_PATH = "unix:/run/ccnp/uds/eventlog.sock"
	TYPE_TDX = "TDX"
	TYPE_TPM = "TPM"
)

type CCEventLogEntry struct {
	RegIdx  uint32
	EvtType uint32
	EvtSize uint32
	AlgId   uint16
	Event   []uint8
	Digest  []uint8
}

type GetPlatformEventlogOptions struct {
	eventlogCategory pb.CATEGORY
	startPosition    int32
	count            int32
}

func WithEventlogCategory(eventlogCategory pb.CATEGORY) func(*GetPlatformEventlogOptions) {
	return func(opts *GetPlatformEventlogOptions) {
		opts.eventlogCategory = eventlogCategory
	}
}

func WithStartPosition(startPosition int32) func(*GetPlatformEventlogOptions) {
	return func(opts *GetPlatformEventlogOptions) {
		opts.startPosition = startPosition
	}
}

func WithCount(count int32) func(*GetPlatformEventlogOptions) {
	return func(opts *GetPlatformEventlogOptions) {
		opts.count = count
	}
}

func isEventlogCategoryValid(eventlogCategory pb.CATEGORY) bool {
	return eventlogCategory == pb.CATEGORY_TDX_EVENTLOG || eventlogCategory == pb.CATEGORY_TPM_EVENTLOG
}

func getRawEventlogs(response *pb.GetEventlogReply) ([]byte, error) {
	path := response.EventlogDataLoc
	if path == "" {
		log.Fatalf("[getRawEventlogs] Failed to get eventlog from server")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("[getRawEventlogs] Error reading data from  %v: %v", path, err)
	}

	return data, nil
}

func parseTdxEventlog(rawEventlog []byte) ([]CCEventLogEntry, error) {
	var jsonEventlog = el.TDEventLogs{}
	err := json.Unmarshal(rawEventlog, &jsonEventlog)
	if err != nil {
		log.Fatalf("[parseEventlog] Error unmarshal raw eventlog: %v", err)
	}

	rawEventLogList := jsonEventlog.EventLogs
	var parsedEventLogList []CCEventLogEntry
	for i := 0; i < len(rawEventLogList); i++ {
		rawEventlog := rawEventLogList[i]
		eventLog := CCEventLogEntry{}

		if rawEventlog.DigestCount < 1 {
			continue
		}

		eventLog.RegIdx = rawEventlog.Rtmr
		eventLog.EvtType = rawEventlog.Etype
		eventLog.EvtSize = rawEventlog.EventSize
		eventLog.AlgId = rawEventlog.AlgorithmId
		eventLog.Event = rawEventlog.Event
		eventLog.Digest = []uint8(rawEventlog.Digests[rawEventlog.DigestCount-1])
		parsedEventLogList = append(parsedEventLogList, eventLog)

	}

	return parsedEventLogList, nil
}

func GetPlatformEventlog(opts ...func(*GetPlatformEventlogOptions)) ([]CCEventLogEntry, error) {

	input := GetPlatformEventlogOptions{eventlogCategory: pb.CATEGORY_TDX_EVENTLOG, startPosition: 0, count: 0}
	for _, opt := range opts {
		opt(&input)
	}

	if !isEventlogCategoryValid(input.eventlogCategory) {
		log.Fatalf("[GetPlatformEventlog] Invalid eventlogCategory specified")
	}

	if input.eventlogCategory == pb.CATEGORY_TPM_EVENTLOG {
		log.Fatalf("[GetPlatformEventlog] TPM to be supported later")
	}

	if input.startPosition < 0 {
		log.Fatalf("[GetPlatformEventlog] Invalid startPosition specified")
	}

	if input.count < 0 {
		log.Fatalf("[GetPlatformEventlog] Invalid count specified")
	}

	channel, err := grpc.Dial(UDS_PATH, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[GetPlatformEventlog] can not connect to UDS: %v", err)
	}
	defer channel.Close()

	client := pb.NewEventlogClient(channel)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	response, err := client.GetEventlog(ctx, &pb.GetEventlogRequest{
		EventlogLevel:    pb.LEVEL_PAAS,
		EventlogCategory: input.eventlogCategory,
		StartPosition:    input.startPosition,
		Count:            input.count,
	})
	if err != nil {
		log.Fatalf("[GetPlatformEventlog] fail to get Platform Eventlog: %v", err)
	}

	switch input.eventlogCategory {
	case pb.CATEGORY_TDX_EVENTLOG:
		rawEventlog, err := getRawEventlogs(response)
		if err != nil {
			log.Fatalf("[GetPlatformEventlog] fail to get raw eventlog: %v", err)
		}

		return parseTdxEventlog(rawEventlog)

	case pb.CATEGORY_TPM_EVENTLOG:
		return nil, pkgerrors.New("[GetPlatformEventlog] vTPM to be supported later")
	default:
		log.Fatalf("[GetPlatformEventlog] unknown TEE enviroment!")
	}

	return nil, nil
}
