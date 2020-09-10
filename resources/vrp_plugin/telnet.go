/*
 * Copyright (c) 2020. Liero AB
 *
 * Permission is hereby granted, free of charge, to any person obtaining
 * a copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software
 * is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
 * OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
 * IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
 * CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
 * TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"fmt"
	"git.liero.se/opentelco/go-swpx/proto/networkelement"
	"regexp"
	"strconv"
	"strings"
)

const (
	MacRegex = "^([[:xdigit:]]{2}[:.-]?){5}[[:xdigit:]]{2}$"
	IPRegex  = "(\\b(?:\\d{1,3}\\.){3}\\d{1,3}\\b)\\s+([0-9A-Fa-f]{4}[-][a-f0-9A-F]{4}[-][a-f0-9A-F]{4})\\s+([0-9]{1,3}).*([1-9][0-9]{3}.[0-9]{2}.[0-9]{2}-[0-9]{2}:[0-9]{2})"
)

func ParseMacTable(data string) ([]*networkelement.MACEntry, error) {
	dataRows := strings.Split(data, "\n")
	rows := make([]*networkelement.MACEntry, 0)

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

		rows = append(rows, &networkelement.MACEntry{
			HardwareAddress: fields[0],
			Vlan:            int64(vlan),
		})
	}

	return rows, nil
}

func isMacAddressRow(fields []string) bool {
	if len(fields) != 4 {
		return false
	}

	match, parseErr := regexp.Match(MacRegex, []byte(fields[0]))
	if parseErr != nil {
		return false
	}

	return match
}

func ParseIPTable(data string) ([]*networkelement.DHCPEntry, error) {
	dataRows := strings.Split(data, "\n")
	rows := make([]*networkelement.DHCPEntry, 0)

	for _, row := range dataRows {
		fields := strings.Fields(row)

		// skip rows without IP data
		if !isIPAddressRow(fields) {
			continue
		}

		vlan, err := getVLAN(fields)
		if err != nil || vlan == 0 {
			return nil, fmt.Errorf("can't convert vlan string to int: %v", err)
		}

		rows = append(rows, &networkelement.DHCPEntry{
			IpAddress:       fields[0],
			HardwareAddress: fields[1],
			Vlan:            int64(vlan),
			Timestamp:       fields[6],
		})
	}

	return rows, nil
}

// looks for populated VLAN field
func getVLAN(fields []string) (int, error) {
	var err error

	for i := 2; i < 5; i++ {
		if fields[i] == "/--" {
			continue
		}
		val, err := strconv.Atoi(fields[i])

		if err == nil && val != 0 {
			return val, nil
		}
	}

	return 0, err
}

func isIPAddressRow(fields []string) bool {
	if len(fields) != 7 {
		return false
	}

	match, parseErr := regexp.Match(IPRegex, []byte(strings.Join(fields, " ")))
	if parseErr != nil {
		return false
	}

	return match
}
