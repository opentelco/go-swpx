package model

import (
	"git.liero.se/opentelco/go-swpx/fleet/internal"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/configurationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (p *ListConfigurationsParams) ToProto() *configurationpb.ListParameters {
	params := &configurationpb.ListParameters{
		DeviceId: p.DeviceID,
		Limit:    internal.PointerIntToPointerInt64(p.Limit),
		Offset:   internal.PointerIntToPointerInt64(p.Offset),
	}
	if p.After != nil {
		params.CreatedAfter = timestamppb.New(*p.After)
	}
	if p.Before != nil {
		params.CreatedBefore = timestamppb.New(*p.Before)
	}

	if p.Order != nil {
		switch *p.Order {
		case ListOrderAsc:
			params.OrderAsc = internal.ToPtr(true)
		}
	}
	if p.OrderBy != nil {
		switch *p.OrderBy {
		case ConfigurationOrderByCreated:
			params.OrderBy = []configurationpb.ListParameters_OrderBy{configurationpb.ListParameters_ORDER_BY_CREATED}
		}
	}

	return params

}
