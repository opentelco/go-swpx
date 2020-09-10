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
	"git.liero.se/opentelco/go-swpx/proto/networkelement"
	"testing"
)

func Test_ParseFullMacTable(t *testing.T) {
	table := `display mac-address
-------------------------------------------------------------------------------
MAC Address    VLAN/VSI                          Learned-From        Type
-------------------------------------------------------------------------------
0035-1a6d-4ebf 55/-                              XGE0/0/1            dynamic
20f1-7cb0-ecbf 55/-                              XGE0/0/1            dynamic
4431-92fa-d95d 55/-                              XGE0/0/1            dynamic
4ac0-6c7e-6f46 55/-                              XGE0/0/1            dynamic
744d-28f0-1354 55/-                              XGE0/0/1            dynamic
8478-ac07-f461 55/-                              XGE0/0/1            dynamic
9c7d-a360-93a8 55/-                              XGE0/0/1            dynamic
9c7d-a360-93c5 55/-                              XGE0/0/1            dynamic
d257-bc69-1da1 55/-                              XGE0/0/1            dynamic
0848-2c20-15a1 999/-                             GE0/0/1             dynamic

-------------------------------------------------------------------------------
Total items displayed = 10

<liero-test-a1>
`

	got, err := ParseMacTable(table)
	want := []*networkelement.MACEntry{
		{HardwareAddress: "0035-1a6d-4ebf", Vlan: 55},
		{HardwareAddress: "20f1-7cb0-ecbf", Vlan: 55},
		{HardwareAddress: "4431-92fa-d95d", Vlan: 55},
		{HardwareAddress: "4ac0-6c7e-6f46", Vlan: 55},
		{HardwareAddress: "744d-28f0-1354", Vlan: 55},
		{HardwareAddress: "8478-ac07-f461", Vlan: 55},
		{HardwareAddress: "9c7d-a360-93a8", Vlan: 55},
		{HardwareAddress: "9c7d-a360-93c5", Vlan: 55},
		{HardwareAddress: "d257-bc69-1da1", Vlan: 55},
		{HardwareAddress: "0848-2c20-15a1", Vlan: 999},
	}

	if err != nil {
		t.Errorf("unexpected error: %v", err.Error())
	}
	if !compareMAC(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func Test_ParseSingleValueMacTable(t *testing.T) {
	table := `<liero-test-a1>display mac-address GigabitEthernet 0/0/1
-------------------------------------------------------------------------------
MAC Address    VLAN/VSI                          Learned-From        Type
-------------------------------------------------------------------------------
0848-2c20-15a1 999/-                             GE0/0/1             dynamic

-------------------------------------------------------------------------------
Total items displayed = 1
`

	got, err := ParseMacTable(table)
	want := []*networkelement.MACEntry{
		{HardwareAddress: "0848-2c20-15a1", Vlan: 999},
	}

	if err != nil {
		t.Errorf("unexpected error: %v", err.Error())
	}
	if !compareMAC(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func Test_ParseEmptyMacTable(t *testing.T) {
	table := `<liero-test-a1>display mac-address GigabitEthernet 0/0/12
-------------------------------------------------------------------------------
MAC Address    VLAN/VSI                          Learned-From        Type
-------------------------------------------------------------------------------

-------------------------------------------------------------------------------
Total items displayed = 0
`

	got, err := ParseMacTable(table)
	want := make([]*networkelement.MACEntry, 0)

	if err != nil {
		t.Errorf("unexpected error: %v", err.Error())
	}
	if !compareMAC(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func Test_ParseSingleValueIPTable(t *testing.T) {
	table := `<liero-test-a1>display dhcp snooping user-bind interface GigabitEthernet 0/0/1
DHCP Dynamic Bind-table:
Flags:O - outer vlan ,I - inner vlan ,P - Vlan-mapping 
IP Address       MAC Address     VSI/VLAN(O/I/P) Interface      Lease           
--------------------------------------------------------------------------------
192.168.112.19   0848-2c20-1599  296 /--  /--    GE0/0/1        2020.08.10-12:27
--------------------------------------------------------------------------------
Print count:           1          Total count:           1         
`

	got, err := ParseIPTable(table)
	want := []*networkelement.DHCPEntry{
		{IpAddress: "192.168.112.19", HardwareAddress: "0848-2c20-1599", Vlan: 296, Timestamp: "2020.08.10-12:27"},
	}

	if err != nil {
		t.Errorf("unexpected error: %v", err.Error())
	}
	if !compareDHCP(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func compareMAC(x, y []*networkelement.MACEntry) bool {
	if len(x) != len(y) {
		return false
	}

	for i := range x {
		if x[i].HardwareAddress != y[i].HardwareAddress || x[i].Vlan != y[i].Vlan {
			return false
		}
	}

	return true
}

func compareDHCP(x, y []*networkelement.DHCPEntry) bool {
	if len(x) != len(y) {
		return false
	}

	for i := range x {
		if x[i].HardwareAddress != y[i].HardwareAddress || x[i].Vlan != y[i].Vlan || x[i].IpAddress != y[i].IpAddress || x[i].Timestamp != y[i].Timestamp {
			return false
		}
	}

	return true
}
