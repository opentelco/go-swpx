package core

import (
	"git.liero.se/opentelco/go-swpx/proto/networkelement"
)

type Response struct {
	RequestObjectID string
	NetworkElement  *networkelement.Element
	PhysicalPort    *networkelement.PhysicalPortInformation
	Transceiver     *networkelement.Transceiver
	Error           error
}
