package mappers

import (
	"git.liero.se/opentelco/go-swpx/fleet/graph/model"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/notificationpb"
)

type Notification struct{ *notificationpb.Notification }
type Notifications []*notificationpb.Notification

func (n Notifications) ToGQL() []*model.Notification {
	respo := []*model.Notification{}
	for _, v := range n {
		respo = append(respo, Notification{v}.ToGQL())
	}
	return respo
}

func (n Notification) ToGQL() *model.Notification {
	respo := &model.Notification{
		ID:           n.Id,
		ResourceID:   n.ResourceId,
		ResourceType: NotificationResourceType(n.ResourceType).ToGQL(),
		Message:      n.Message,
		Read:         n.Read,
		Timestamp:    n.Timestamp.AsTime(),
	}
	return respo
}

type NotificationResourceType notificationpb.Notification_ResourceType

func (n NotificationResourceType) ToGQL() model.NotificationResourceType {
	switch notificationpb.Notification_ResourceType(n) {
	case notificationpb.Notification_RESOURCE_TYPE_DEVICE:
		return model.NotificationResourceTypeDevice
	case notificationpb.Notification_RESOURCE_TYPE_CONFIG:
		return model.NotificationResourceTypeConfig
	case notificationpb.Notification_RESOURCE_TYPE_UNSPECIFIED:
		return model.NotificationResourceTypeUnspecified
	default:
		return model.NotificationResourceTypeUnspecified
	}
}

type ListNotificationResponse struct{ *notificationpb.ListResponse }

func (l ListNotificationResponse) ToGQL() *model.ListNotificationsResponse {
	respo := &model.ListNotificationsResponse{
		Notifications: Notifications(l.Notifications).ToGQL(),
		PageInfo:      ToPageInfo(l.PageInfo).ToGQL(),
	}
	return respo

}

func ToListNotificationResponse(in *notificationpb.ListResponse) *ListNotificationResponse {
	return &ListNotificationResponse{in}
}
