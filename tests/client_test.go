package tests

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"go.opentelco.io/go-swpx/proto/go/corepb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func Test_Grpc(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "127.0.0.1:1338", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)

	}

	client := corepb.NewCoreServiceClient(conn)
	r, err := client.Poll(ctx, &corepb.PollRequest{
		Settings: &corepb.Settings{
			ProviderPlugin: []string{"vx"},
			ResourcePlugin: "vrp",
			RecreateIndex:  false,
			Timeout:        "60s",
			CacheTtl:       "0s",
		},
		Session: &corepb.SessionRequest{
			Hostname: "bryanston-west-a39",
			Port:     "GigabitEthernet0/0/1",
		},
	})
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}

	fmt.Printf("result: %+v", r)

	ne := r.NetworkElement
	b, err := protojson.MarshalOptions{
		Multiline:       false,
		Indent:          "  ",
		AllowPartial:    false,
		UseProtoNames:   false,
		UseEnumNumbers:  false,
		EmitUnpopulated: false,
		Resolver:        nil,
	}.Marshal(ne)

	if err != nil {
		fmt.Println("could not marshal", err.Error())
	}

	fmt.Println(string(b))
}
