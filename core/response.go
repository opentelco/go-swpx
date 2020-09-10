package core

import (
	"git.liero.se/opentelco/go-swpx/proto/networkelement"
)

type Response struct {
	RequestObjectID string                      `json:"request_object_id,omitempty" bson:"request_object_id"`
	NetworkElement  *networkelement.Element     `json:"network_element,omitempty" bson:"network_element"`
	PhysicalPort    string                      `json:"physical_port,omitempty" bson:"physical_port"`
	Transceiver     *networkelement.Transceiver `json:"transceiver,omitempty" bson:"transceiver"`
	Error           error                       `json:"error,omitempty" bson:"error"`
}
