package main

import (
	"testing"

	"github.com/opentelco/go-dnc/models/protobuf/transport"
	proto "github.com/opentelco/go-swpx/proto/resource"
)

func TestMapInterface(t *testing.T) {

	driver := &VRPDriver{
		logger: nil,
		dnc:    &MockClient{},
	}

	req := &proto.NetworkElement{
		Hostname:  "någon-host",
		Interface: "GigabitEthernet0/0/1",
	}

	msg := createDiscoveryMsg(req)
	msg, err := driver.dnc.Put(msg)
	if err != nil {
		t.Errorf(err.Error())
	}

}

type MockClient struct{}

func (m MockClient) Ping() (string, error) {
	return "", nil
}

func (m MockClient) Put(*transport.Message) (*transport.Message, error) {

	return &transport.Message{}, nil
}

func (m MockClient) Close() error {
	return nil
}

func TestInterfaceDescriptionIndexParser(t *testing.T) {
	return nil
}
