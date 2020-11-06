package tests

import (
	"context"
	"encoding/json"
	"git.liero.se/opentelco/go-swpx/proto/go/resource"
	"google.golang.org/grpc"
	"log"
	"testing"
)

// This test is run along with the cache and DNC instances.
// Also, the gRPC server needs to be used in the Start command in root.go instead of the regular one
func Test_GRPCEndpoint(t *testing.T) {
	req := &resource.TechnicalInformationRequest{
		Hostname:      "10.5.5.100",
		Port:          "GigabitEthernet0/0/1",
		Provider:      "zitius_provider",
		Driver:        "resource_vrp_plug",
		Region:        "queue-area1-1",
		Timeout:       20,
		RecreateIndex: false,
		CacheTtl:      30,
	}

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":1337", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s. Make sure that the DNC instances are running", err)
	}
	defer conn.Close()
	c := resource.NewTechnicalInformationClient(conn)

	response, err := c.TechnicalPortInformation(context.Background(), req)
	if err != nil {
		log.Printf("Error when calling endpoint: %s", err)
		t.Fail()
	}

	fmt, _ := json.MarshalIndent(response, "", "\t")
	log.Printf("Response from server: %s", fmt)
}
