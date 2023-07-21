package main

import (
	"context"

	"git.liero.se/opentelco/go-swpx/proto/go/resourcepb"
	"git.liero.se/opentelco/go-swpx/proto/go/stanzapb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (d *driver) ConfigureStanza(ctx context.Context, req *resourcepb.ConfigureStanzaRequest) (*stanzapb.ConfigureResponse, error) {
	return nil, status.New(codes.Unimplemented, "not implemented").Err()
}
