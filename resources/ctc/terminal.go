package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"go.opentelco.io/go-swpx/proto/go/devicepb"
)

const (
	MacRegex = "^([[:xdigit:]]{2}[:.-]?){5}[[:xdigit:]]{2}$"
)

func parseMacTable(data string) ([]*devicepb.MACEntry, error) {
	if data == "" {
		return nil, fmt.Errorf("no data found")
	}

	dataRows := strings.Split(data, "\n")
	rows := make([]*devicepb.MACEntry, 0)

	for _, row := range dataRows {
		fields := strings.Fields(row)

		// skip rows without mac data
		if !isMacAddressRow(fields) {
			continue
		}
		vlan, err := strconv.Atoi(strings.TrimRight(fields[1], "/-"))
		if err != nil {
			logger.Error("can't convert VLAN for mac entry table: ", err.Error())
			return nil, err
		}

		rows = append(rows, &devicepb.MACEntry{
			HardwareAddress: fields[2],
			Vlan:            int64(vlan),
			Vendor:          fields[3],
		})
	}

	return rows, nil
}

func isMacAddressRow(fields []string) bool {
	if len(fields) != 6 {
		return false
	}

	match, parseErr := regexp.Match(MacRegex, []byte(fields[2]))
	if parseErr != nil {
		return false
	}

	return match
}
