package notification

import (
	"fmt"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/notificationpb"
)

func NewDeviceCouldNotBeDiscovered(identifier, reason string) *notificationpb.CreateRequest {
	msg := fmt.Sprintf("The device %s could not be discovered. Please check the device and try again. Reason:,%s", identifier, reason)
	return &notificationpb.CreateRequest{
		Title:        fmt.Sprintf("Device %s could not be discovered", identifier),
		ResourceType: notificationpb.Notification_RESOURCE_TYPE_DEVICE,
		Message:      &msg,
	}
}

func NewDeviceDiscovered(deviceId, identifier string) *notificationpb.CreateRequest {
	msg := fmt.Sprintf("The device (%s) with ID: %s has been discovered", identifier, deviceId)
	return &notificationpb.CreateRequest{
		Title:        fmt.Sprintf("Device %s discovered", identifier),
		ResourceId:   deviceId,
		ResourceType: notificationpb.Notification_RESOURCE_TYPE_DEVICE,
		Message:      &msg,
	}
}

func NewDeviceConfigurationCollectionFailed(deviceId, identifier, reason string) *notificationpb.CreateRequest {
	msg := fmt.Sprintf("Failed to collect configuration for device %s with ID: %s. Reason: %s  ", identifier, deviceId, reason)
	return &notificationpb.CreateRequest{
		Title:        fmt.Sprintf("Failed to collect config for %s", identifier),
		ResourceId:   deviceId,
		ResourceType: notificationpb.Notification_RESOURCE_TYPE_DEVICE,
		Message:      &msg,
	}
}

func NewDeviceCollectionFailed(deviceId, identifier, reason string) *notificationpb.CreateRequest {
	msg := fmt.Sprintf("Failed to collect device %s with ID: %s, Reason: %s  ", identifier, deviceId, reason)
	return &notificationpb.CreateRequest{
		Title:        fmt.Sprintf("Failed to collect device: %s", identifier),
		ResourceId:   deviceId,
		ResourceType: notificationpb.Notification_RESOURCE_TYPE_DEVICE,
		Message:      &msg,
	}
}

func NewDeviceConfigScheduleCancelled(deviceId, identifier string) *notificationpb.CreateRequest {
	msg := fmt.Sprintf("To many failures to collect config for device (%s) with ID: %s ", identifier, deviceId)
	return &notificationpb.CreateRequest{
		Title:        fmt.Sprintf("Collect config schedule disabled for %s", identifier),
		ResourceId:   deviceId,
		ResourceType: notificationpb.Notification_RESOURCE_TYPE_DEVICE,
		Message:      &msg,
	}
}
