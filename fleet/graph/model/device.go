package model

import (
	"git.liero.se/opentelco/go-swpx/fleet/internal"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
)

func (p *ListDevicesParams) ToProto() *devicepb.ListParameters {

	params := &devicepb.ListParameters{
		Limit:  internal.PointerIntToPointerInt64(p.Limit),
		Offset: internal.PointerIntToPointerInt64(p.Offset),
	}

	if p.Hostname != nil {
		params.Hostname = p.Hostname
	}

	if p.ManagementIP != nil {
		params.ManagementIp = p.ManagementIP
	}

	return params
}
