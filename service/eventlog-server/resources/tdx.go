/* SPDX-license-identifier: Apache-2.0*/

package resources

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	pkgerrors "github.com/pkg/errors"
)

const (
	//The location of CCEL table
	CCEL_FILE_LOCATION = "/sys/firmware/acpi/tables/CCEL"
	CCEL_DATA_LOCATION = "/sys/firmware/acpi/tables/data/CCEL"
	//The location of mounted CCEL table
	CCEL_FILE_MOUNT_LOCATION = "/run/firmware/acpi/tables/CCEL"
	CCEL_DATA_MOUNT_LOCATION = "/run/firmware/acpi/tables/data/CCEL"

	EVENT_TYPE_EV_NO_ACTION = 0x3
)

var (
	TdxGetEventlogErr     = pkgerrors.New("Failed to get eventlog in CCEL table.")
	CcelTableNotFoundErr  = pkgerrors.New("CCEL table not found.")
	InvalidCcelTableErr   = pkgerrors.New("CCEL table with invalid data")
	FetchCcelTableAttrErr = pkgerrors.New("Failed to get the base address of CCEL table")
)

func GetTdxEventlog(start_position int, count int) (string, error) {

	var eventlog string
	var object *os.File
	var err error

	/* Check if the ccel file exists*/
	if _, err = os.Stat(CCEL_FILE_MOUNT_LOCATION); err != nil {
		log.Println("Checking CCEL file in host path")
		if _, err = os.Stat(CCEL_FILE_LOCATION); err != nil {
			return "", err
		}
	}

	/* Open TDX device fd to get prepared for TDVM call*/
	object, err = os.OpenFile(CCEL_FILE_MOUNT_LOCATION, os.O_RDONLY, 0644)
	if err != nil {
		object, err = os.OpenFile(CCEL_FILE_LOCATION, os.O_RDONLY, 0644)
		if err != nil {
			return "", CcelTableNotFoundErr
		}
	}

	data, err := io.ReadAll(object)
	if err != nil {
		return "", err
	}

	if len(data) == 0 || !bytes.Equal(data[0:4], []byte("CCEL")) {
		return "", InvalidCcelTableErr
	}

	eventlog, err = parseTdxEventlog(data, start_position, count)
	if err != nil {
		return "", err
	}

	return eventlog, nil
}

func parseTdxEventlog(data []byte, position int, count int) (string, error) {

	var baseAddr, currLen uint64
	var num int
	var eventlogs TDEventLogs

	if len(data) > 56 {
		return "", InvalidCcelTableErr
	}

	err := binary.Read(bytes.NewReader(data[48:56]), binary.LittleEndian, &baseAddr)
	if err != nil {
		return "", FetchCcelTableAttrErr
	}

	err = binary.Read(bytes.NewReader(data[40:48]), binary.LittleEndian, &currLen)
	if err != nil {
		return "", FetchCcelTableAttrErr
	}

	eventlogs, num, err = fetchEventlogs()
	if err != nil {
		return "", err
	}

	if position+count >= num {
		return "", pkgerrors.New("Invalid count exceeds event log length")
	}

	if count != 0 {
		eventlogs = TDEventLogs{
			Header:    eventlogs.Header,
			EventLogs: eventlogs.EventLogs[position : position+count],
		}
	}

	eventlogs_str, err := json.Marshal(eventlogs)
	if err != nil {
		log.Println("Error in marshaling event logs")
		return "", err
	}
	return string(eventlogs_str), nil
}

func fetchEventlogs() (TDEventLogs, int, error) {

	var index int
	var object *os.File
	var err error

	/* Check if the ccel file exists in either host or container*/
	if _, err = os.Stat(CCEL_DATA_MOUNT_LOCATION); err != nil {
		log.Println("Checking CCEL data in host path")
		_, err = os.Stat(CCEL_DATA_LOCATION)
		if err != nil {
			return TDEventLogs{}, 0, err
		}
	}

	/* Open TDX device fd to get prepared for TDVM call*/
	object, err = os.OpenFile(CCEL_DATA_MOUNT_LOCATION, os.O_RDONLY, 0644)
	if err != nil {
		object, err = os.OpenFile(CCEL_DATA_LOCATION, os.O_RDONLY, 0644)
		if err != nil {
			return TDEventLogs{}, 0, CcelTableNotFoundErr
		}
	}

	data, err := io.ReadAll(object)
	if err != nil {
		return TDEventLogs{}, 0, err
	}

	if len(data) == 0 {
		return TDEventLogs{}, 0, CcelTableNotFoundErr
	}

	eventLogs := TDEventLogs{}
	specidHeader := TDEventLogSpecIdHeader{}
	index = 0
	count := 0

	for index < len(data) {
		var rtmr, etype uint32
		start := index

		rtmr, index, err = getUint32Object(data, index)
		if err != nil {
			log.Println("Error in getting RTMR value")
			return TDEventLogs{}, 0, err
		}

		if rtmr == 0xFFFFFFFF {
			break
		}

		etype, _, err = getUint32Object(data, index)
		if err != nil {
			log.Println("Error in getting event type")
			return TDEventLogs{}, 0, err
		}

		if etype == EVENT_TYPE_EV_NO_ACTION {

			specidHeader.Rtmr, specidHeader.Etype, specidHeader.DigestCount, index, err = getBasicInfo(data[start:])
			if err != nil {
				log.Println("Error in getting basic info")
				return TDEventLogs{}, 0, err
			}

			index += 20
			index += 24

			specidHeader.DigestSizes, index, err = getHeaderDigestInfo(data[start:], index)
			if err != nil {
				log.Println("Error in getting header digest info")
				return TDEventLogs{}, 0, err
			}

			var vendorSize uint8
			vendorSize, index, err = getUint8Object(data[start:], index)
			if err != nil {
				log.Println("Error in getting vendor size")
				return TDEventLogs{}, 0, err
			}

			index = index + int(vendorSize)
			specidHeader.Length = int(index)
			specidHeader.HeaderData = data[0:index]
			eventLogs.Header = specidHeader

			index = start + specidHeader.Length
			count += 1
			continue

		}

		eventLog := TDEventLog{}
		eventLog.Rtmr, eventLog.Etype, eventLog.DigestCount, index, err = getBasicInfo(data[start:])
		if err != nil {
			log.Println("Error in getting basic info")
			return TDEventLogs{}, 0, err
		}

		eventLog.Digests, eventLog.AlgorithmId, index, err = getEventLogDigestInfo(data[start:], index, eventLog.DigestCount, eventLogs.Header.DigestSizes)
		if err != nil {
			log.Println("Error in getting event log digest info")
			return TDEventLogs{}, 0, err
		}

		eventLog.EventSize, index, err = getUint32Object(data[start:], index)
		if err != nil {
			log.Println("Error in getting event size")
			return TDEventLogs{}, 0, err
		}

		eventLog.Event = data[int(start)+index : int(start)+index+int(eventLog.EventSize)]
		index = index + int(eventLog.EventSize)
		eventLog.Length = index
		eventLog.Data = data[start : start+index]
		eventLogs.EventLogs = append(eventLogs.EventLogs, eventLog)
		index = start + eventLog.Length

		count += 1
	}

	return eventLogs, count, nil
}

type TDEventLogSpecIdHeader struct {
	Address     uint64
	Length      int
	HeaderData  []byte
	Rtmr        uint32
	Etype       uint32
	DigestCount uint32
	DigestSizes map[uint16]uint16
}

type TDEventLog struct {
	Rtmr        uint32
	Etype       uint32
	DigestCount uint32
	Digests     []string
	Data        []byte
	Event       []byte
	Length      int
	EventSize   uint32
	AlgorithmId uint16
}

type TDEventLogs struct {
	Header    TDEventLogSpecIdHeader
	EventLogs []TDEventLog
}

func getBasicInfo(data []byte) (uint32, uint32, uint32, int, error) {

	var rtmr, etype, digestCount uint32
	var err error
	i := 0

	rtmr, i, err = getUint32Object(data, i)
	if err != nil {
		log.Println("Error in getting RTMR")
		return uint32(0), uint32(0), uint32(0), 0, err
	}

	etype, i, err = getUint32Object(data, i)
	if err != nil {
		log.Println("Error in getting event log type")
		return uint32(0), uint32(0), uint32(0), 0, err
	}

	digestCount, i, err = getUint32Object(data, i)
	if err != nil {
		log.Println("Error in getting digest count")
		return uint32(0), uint32(0), uint32(0), 0, err
	}

	return rtmr - 1, etype, digestCount, i, nil
}

func getHeaderDigestInfo(data []byte, index int) (map[uint16]uint16, int, error) {

	var algNum uint32
	var algId, digestSize uint16
	var err error

	digestSizes := make(map[uint16]uint16)
	i := index

	algNum, i, err = getUint32Object(data, i)
	if err != nil {
		log.Println("Error in getting algorithm number")
		return digestSizes, 0, err
	}

	for j := 0; j < int(algNum); j++ {
		algId, i, err = getUint16Object(data, i)
		if err != nil {
			log.Println("Error in getting algorithm id")
			return digestSizes, 0, err
		}

		digestSize, i, err = getUint16Object(data, i)
		if err != nil {
			log.Println("Error in getting algorithm id")
			return digestSizes, 0, err
		}

		digestSizes[algId] = digestSize
	}

	return digestSizes, i, nil
}

func getEventLogDigestInfo(data []byte, index int, digestCount uint32, digestSizes map[uint16]uint16) ([]string, uint16, int, error) {

	var algId uint16
	var err error
	var digests []string

	i := index

	for j := 0; j < int(digestCount); j++ {
		algId, i, err = getUint16Object(data, i)
		if err != nil {
			log.Println("Error in getting algorithm id")
			return digests, uint16(0), 0, err
		}

		for k := range digestSizes {
			if k == algId {
				digestSize := digestSizes[k]
				digestData := data[i : i+int(digestSize)]
				i = i + int(digestSize)
				digests = append(digests, fmt.Sprintf("%v", digestData))
			}
		}
	}

	return digests, algId, i, nil
}

func getUint32Object(data []byte, index int) (uint32, int, error) {
	var value uint32

	if index+4 > len(data) {
		return uint32(0), index, pkgerrors.New("Exceed valid length")
	}
	err := binary.Read(bytes.NewReader(data[index:index+4]), binary.LittleEndian, &value)
	if err != nil {
		log.Println("Error in reading uint32 object")
		return uint32(0), index, err
	}

	return value, index + 4, nil
}

func getUint16Object(data []byte, index int) (uint16, int, error) {
	var value uint16

	if index+2 > len(data) {
		return uint16(0), index, pkgerrors.New("Exceed valid length")
	}
	err := binary.Read(bytes.NewReader(data[index:index+2]), binary.LittleEndian, &value)
	if err != nil {
		log.Println("Error in reading uint16 object")
		return uint16(0), index, err
	}

	return value, index + 2, nil
}

func getUint8Object(data []byte, index int) (uint8, int, error) {
	var value uint8

	if index+1 > len(data) {
		return uint8(0), index, pkgerrors.New("Exceed valid length")
	}
	err := binary.Read(bytes.NewReader(data[index:index+1]), binary.LittleEndian, &value)
	if err != nil {
		log.Println("Error in reading uint8 object")
		return uint8(0), index, err
	}

	return value, index + 1, nil
}
