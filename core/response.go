package core

import (
	"git.liero.se/opentelco/go-swpx/proto/networkelement"
)

type Response struct {
	RequestObjectID string
	NetworkElement  *networkelement.Element
	PhysicalPort    string
	Transceiver     *networkelement.Transceiver
	Error           error `json:"error,omitempty"`
}
