package main

import (
	"context"
	"git.liero.se/opentelco/go-dnc/models/protobuf/dispatcher"
	"git.liero.se/opentelco/go-dnc/models/protobuf/transport"
	"git.liero.se/opentelco/go-swpx/shared"
	"testing"

	proto "git.liero.se/opentelco/go-swpx/proto/resource"
)

func TestMapInterface(t *testing.T) {

	driver := &VRPDriver{
		logger: nil,
		dnc:    &MockClient{},
		conf:   shared.Configuration{},
	}

	req := &proto.NetworkElement{
		Hostname:  "någon-host",
		Interface: "GigabitEthernet0/0/1",
	}

	msg := createDiscoveryMsg(req, driver.conf)
	msg, err := driver.dnc.Put(context.Background(), msg)
	if err != nil {
		t.Errorf(err.Error())
	}

}

type MockClient struct{}

func (m MockClient) Ping() (*dispatcher.PingReply, error) {
	return &dispatcher.PingReply{}, nil
}

func (m MockClient) Put(ctx context.Context, msg *transport.Message) (*transport.Message, error) {

	return &transport.Message{}, nil
}

func (m MockClient) Close() error {
	return nil
}

func TestInterfaceDescriptionIndexParser(t *testing.T) {
}
