package graph

import (
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/notificationpb"
)

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	devices       devicepb.DeviceServiceServer
	notifications notificationpb.NotificationServiceServer
}
