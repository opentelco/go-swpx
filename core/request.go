package core

import "context"

// ReqestType is the possible requsts that can be made to swpx
// they maps to a switchcase on each request
type RequestType uint

// The request possible to use in swpx right now
const (
	GetTechnicalInformationPort RequestType = iota
	GetTechnicalInformationElement
)

// Request is the internal representation of a incoming request
// it is passed between the api and the core
type Request struct {
	ObjectID                string
	NetworkElement          string
	NetworkElementInterface *string
	Provider                string
	Resource                string
	DontUseIndex            bool

	// metadata to handle the request
	Response chan *Response
	Context  context.Context
	Type     RequestType
}
