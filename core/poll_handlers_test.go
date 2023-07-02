package core

import (
	"fmt"
	"testing"

	"git.liero.se/opentelco/go-swpx/proto/go/corepb"
)

func Test_RequestResolver(t *testing.T) {

	req := &corepb.PollRequest{
		Session: &corepb.SessionRequest{
			Hostname:      "hostname-a1",
			Port:          "Port1",
			NetworkRegion: "VX_DE1",
		},
		Settings: &corepb.Settings{
			ProviderPlugin: []string{"vx"},
			ResourcePlugin: "",
		},
		Type: corepb.PollRequest_GET_BASIC_INFO,
	}

	test_x(req.Session)

	fmt.Println(req.Session.Hostname)

}

func test_x(sess *corepb.SessionRequest) {
	sess.Hostname = "hostname-a1-proccessed"
	sess.Port = "Port1-processed"
	sess.NetworkRegion = "VX_DE1-processed"
}
