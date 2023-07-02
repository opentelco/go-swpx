package core

import (
	"fmt"
	"testing"

	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
)

func Test_RequestResolver(t *testing.T) {

	req := &pb_core.PollRequest{
		Session: &pb_core.SessionRequest{
			Hostname:      "hostname-a1",
			Port:          "Port1",
			NetworkRegion: "VX_DE1",
		},
		Settings: &pb_core.Settings{
			ProviderPlugin: []string{"vx"},
			ResourcePlugin: "",
		},
		Type: pb_core.PollRequest_GET_BASIC_INFO,
	}

	test_x(req.Session)

	fmt.Println(req.Session.Hostname)

}

func test_x(sess *pb_core.SessionRequest) {
	sess.Hostname = "hostname-a1-proccessed"
	sess.Port = "Port1-processed"
	sess.NetworkRegion = "VX_DE1-processed"
}
