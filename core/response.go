package core

import "github.com/opentelco/go-swpx/proto/networkelement"

type Response struct {
	RequestObjectID string
	NetworkElement  *networkelement.Element
	Error           error
}
