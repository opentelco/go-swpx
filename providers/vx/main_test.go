package main

import (
	"context"
	"os"
	"testing"

	"git.liero.se/opentelco/go-swpx/proto/go/core"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

func TestProvider_ResolveSessionRequest(t *testing.T) {
	logger := hclog.NewNullLogger()
	ctx := context.Background()
	p := Provider{logger: logger}
	os.Setenv("REGION", "VX_SA1")
	req := &core.SessionRequest{
		Hostname: "testa-a1",
	}
	req, err := p.ResolveSessionRequest(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, req.Hostname, "testa-a1.joburg.net.venturanext.se")
}

func TestProvider_ResolveSessionRequest_1(t *testing.T) {
	logger := hclog.NewNullLogger()
	ctx := context.Background()
	p := Provider{logger: logger}
	os.Setenv("REGION", "")
	req := &core.SessionRequest{
		Hostname: "testa-a1",
	}
	req, err := p.ResolveSessionRequest(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, req.Hostname, "testa-a1")
}

func TestProvider_ResolveSessionRequest_Ip(t *testing.T) {
	logger := hclog.NewNullLogger()
	ctx := context.Background()
	p := Provider{logger: logger}
	os.Setenv("REGION", "VX_SA1")
	req := &core.SessionRequest{
		Hostname: "172.26.11.232",
	}
	req, err := p.ResolveSessionRequest(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, req.Hostname, "172.26.11.232")
}

func Test_setEnv(t *testing.T) {
	err := setupEnv()
	assert.NoError(t, err)
}

func Test_ResolveHostname(t *testing.T) {
	fqdn := "testa-a1.joburg.net.venturanext.se"
	host := HostFromFQDN(fqdn)
	assert.Equal(t, host, "testa-a1")

	fqdn = "testa-a1"
	host = HostFromFQDN(fqdn)
	assert.Equal(t, host, "testa-a1")

}
