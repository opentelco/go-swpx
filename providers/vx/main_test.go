package main

import (
	"context"
	"os"
	"testing"

	"git.liero.se/opentelco/go-swpx/proto/go/core"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func TestProvider_PreHandler(t *testing.T) {
	logger := hclog.NewNullLogger()
	ctx := context.Background()
	p := Provider{logger: logger}
	os.Setenv("REGION", "VX_SA1")
	req := &core.Request{
		Hostname: "testa-a1",
	}
	req, err := p.PreHandler(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, req.Hostname, "testa-a1.joburg.net.venturanext.se")
}

func TestProvider_PreHandler_1(t *testing.T) {
	logger := hclog.NewNullLogger()
	ctx := context.Background()
	p := Provider{logger: logger}
	os.Setenv("REGION", "")
	req := &core.Request{
		Hostname: "testa-a1",
	}
	req, err := p.PreHandler(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, req.Hostname, "testa-a1")
}

func TestProvider_PreHandler_Ip(t *testing.T) {
	logger := hclog.NewNullLogger()
	ctx := context.Background()
	p := Provider{logger: logger}
	os.Setenv("REGION", "VX_SA1")
	req := &core.Request{
		Hostname: "172.26.11.232",
	}
	req, err := p.PreHandler(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, req.Hostname, "172.26.11.232")
}
