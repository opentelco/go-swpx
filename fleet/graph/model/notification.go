package model

import (
	"git.liero.se/opentelco/go-swpx/fleet/internal"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/notificationpb"
)

func (n *ListNotificationsParams) ToProto() *notificationpb.ListRequest {

	if n == nil {
		return &notificationpb.ListRequest{}
	}

	params := &notificationpb.ListRequest{
		ResourceIds: n.ResourceIds,
		Filter:      ListNotificationFilters(n.Filter).ToProto(),
	}

	if n.Limit != nil {
		params.Limit = internal.PointerIntToPointerInt64(n.Limit)
	}
	if n.Offset != nil {
		params.Offset = internal.PointerIntToPointerInt64(n.Offset)
	}

	return params
}

type ListNotificationFilters []ListNotificationFilter

func (f ListNotificationFilters) ToProto() []notificationpb.ListRequest_Filter {
	respo := []notificationpb.ListRequest_Filter{}
	for _, v := range f {
		respo = append(respo, v.ToProto())
	}
	return respo

}

func (f ListNotificationFilter) ToProto() notificationpb.ListRequest_Filter {
	switch f {
	case ListNotificationFilterIncludeRead:
		return notificationpb.ListRequest_INCLUDE_READ
	case ListNotificationFilterResourceTypeConfig:
		return notificationpb.ListRequest_RESOURCE_TYPE_CONFIG
	case ListNotificationFilterResourceTypeDevice:
		return notificationpb.ListRequest_RESOURCE_TYPE_DEVICE
	}
	return notificationpb.ListRequest_LIST_REQUEST_FILTER_UNSET
}

func (p MarkNotificationsAsReadParams) ToProto() *notificationpb.MarkAsReadRequest {
	return &notificationpb.MarkAsReadRequest{
		Ids: p.Ids,
	}
}
