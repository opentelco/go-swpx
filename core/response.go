package core

import (
	"git.liero.se/opentelco/go-swpx/proto/networkelement"
	proto "git.liero.se/opentelco/go-swpx/proto/resource"
)

type Response struct {
	RequestObjectID string
	NetworkElement  *networkelement.Element
	PhysicalPort    *proto.NetworkElementInterface
	Transceiver     *networkelement.Transceiver
	Error           error
}
