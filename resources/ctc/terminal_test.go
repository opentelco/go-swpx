package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var macs = `Type VID MAC Address Vendor Ports
---- --- ----------- ------ -----
Dynamic 296 00:02:ab:c1:30:aa CTCUnion Gi 1/1
Dynamic 921 00:02:ab:c1:30:aa CTCUnion Gi 1/1
Dynamic 921 5c:64:8e:5e:55:73 Unknown Gi 1/1`

func Test_isMacAddressRow(t *testing.T) {
	dataRows := strings.Split(macs, "\n")

	for i, row := range dataRows {
		fields := strings.Fields(row)

		if i == 0 || i == 1 {
			assert.Equal(t, false, isMacAddressRow(fields))
		} else {
			fmt.Println("else", fields)
			assert.Equal(t, true, isMacAddressRow(fields))
		}

	}
}

func Test_parseMacTable(t *testing.T) {
	entries, err := parseMacTable(macs)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(entries))

	assert.Equal(t, "00:02:ab:c1:30:aa", entries[0].HardwareAddress)
	assert.Equal(t, "CTCUnion", entries[0].Vendor)
	assert.Equal(t, int64(296), entries[0].Vlan)

	assert.Equal(t, "00:02:ab:c1:30:aa", entries[1].HardwareAddress)
	assert.Equal(t, "CTCUnion", entries[1].Vendor)
	assert.Equal(t, int64(921), entries[1].Vlan)

	assert.Equal(t, "5c:64:8e:5e:55:73", entries[2].HardwareAddress)
	assert.Equal(t, "Unknown", entries[2].Vendor)
	assert.Equal(t, int64(921), entries[2].Vlan)

}
