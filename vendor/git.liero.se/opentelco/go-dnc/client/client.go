package client

import (
	"git.liero.se/opentelco/go-dnc/models/protobuf/dispatcher"
	"git.liero.se/opentelco/go-dnc/models/protobuf/transport"
)

// Client is the client for DNC
type Client interface {
	Ping() (*dispatcher.PingReply, error)
	Put(*transport.Message) (*transport.Message, error)

	Close() error
}
