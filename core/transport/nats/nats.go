package nats

import (
	"log"
	"strings"

	"github.com/nats-io/go-nats"
)

func init() {}

// New Returns a Nats connection

type Connection struct {
	*nats.EncodedConn
}

func (c *Connection) Ping() {}

func New(uri []string) (*Connection, error) {

	nc, err := nats.Connect(strings.Join(uri, ","))
	c, _ := nats.NewEncodedConn(nc, "json")

	if err != nil {
		log.Fatal(err)
	}

	cc := &Connection{c}
	return cc, nil

}
