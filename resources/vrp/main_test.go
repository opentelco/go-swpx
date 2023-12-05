/*
 * Copyright (c) 2023. Liero AB
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
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelco.io/go-dnc/models/pb/transportpb"
	"go.opentelco.io/go-swpx/config"
	"go.opentelco.io/go-swpx/proto/go/resourcepb"
)

func TestMapInterface(t *testing.T) {

	driver := &VRPDriver{
		logger: nil,
		dnc:    &MockClient{},
		conf:   &config.ResourceVRP{},
	}

	req := &resourcepb.Request{
		Hostname: "some-host",
		Port:     "GigabitEthernet0/0/1",
	}

	msg := createLogicalPortIndex(req, driver.conf)
	msg, err := driver.dnc.Put(context.Background(), msg)
	if err != nil {
		t.Errorf(err.Error())
	}

}

type MockClient struct{}

func (m MockClient) Put(ctx context.Context, msg *transportpb.Message) (*transportpb.Message, error) {

	return &transportpb.Message{}, nil
}

func (m MockClient) Close() error {
	return nil
}

func TestInterfaceDescriptionIndexParser(t *testing.T) {
}

func Test_ParseConfig(t *testing.T) {
	f, err := os.ReadFile("tests/vrp_config.cfg")
	assert.NoError(t, err)

	strConf := string(f)
	fmt.Println(cleanConfig(strConf))

}
