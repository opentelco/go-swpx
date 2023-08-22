package model

import (
	"git.liero.se/opentelco/go-swpx/fleet/internal"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (p *ListDevicesParams) ToProto() *devicepb.ListParameters {
	if p == nil {
		return &devicepb.ListParameters{}
	}

	params := &devicepb.ListParameters{
		Search:       p.Search,
		Hostname:     p.Hostname,
		ManagementIp: p.ManagementIP,
		Limit:        internal.PointerIntToPointerInt64(p.Limit),
		Offset:       internal.PointerIntToPointerInt64(p.Offset),
	}

	if p.Hostname != nil {
		params.Hostname = p.Hostname
	}

	if p.ManagementIP != nil {
		params.ManagementIp = p.ManagementIP
	}

	return params
}

func (p *ListDeviceChangesParams) ToProto() *devicepb.ListChangesParameters {
	if p == nil {
		return &devicepb.ListChangesParameters{}
	}
	params := &devicepb.ListChangesParameters{
		DeviceId: p.DeviceID,
		Limit:    internal.PointerIntToPointerInt64(p.Limit),
		Offset:   internal.PointerIntToPointerInt64(p.Offset),
	}

	if p.Order != nil {
		switch *p.Order {
		case ListOrderAsc:
			params.OrderAsc = internal.ToPtr(true)
		}
	}
	if p.OrderBy != nil {
		switch *p.OrderBy {
		case ListDeviceChangesOrderByCreatedAt:
			params.OrderBy = []devicepb.ListChangesParameters_OrderBy{devicepb.ListChangesParameters_ORDER_BY_CREATED}
		}
	}

	if p.After != nil {
		params.After = timestamppb.New(*p.After)
	}
	if p.Before != nil {
		params.Before = timestamppb.New(*p.Before)
	}

	return params
}

func (p *ListDeviceEventsParams) ToProto() *devicepb.ListEventsParameters {
	if p == nil {
		return &devicepb.ListEventsParameters{}
	}

	params := &devicepb.ListEventsParameters{
		DeviceId: p.DeviceID,
		Limit:    internal.PointerIntToPointerInt64(p.Limit),
		Offset:   internal.PointerIntToPointerInt64(p.Offset),
	}

	if p.Order != nil {
		switch *p.Order {
		case ListOrderAsc:
			params.OrderAsc = internal.ToPtr(true)
		}
	}
	if p.OrderBy != nil {
		switch *p.OrderBy {
		case ListDeviceEventsOrderByCreatedAt:
			params.OrderBy = []devicepb.ListEventsParameters_OrderBy{devicepb.ListEventsParameters_ORDER_BY_CREATED}
		}
	}

	if p.After != nil {
		params.After = timestamppb.New(*p.After)
	}
	if p.Before != nil {
		params.Before = timestamppb.New(*p.Before)
	}

	return params
}
