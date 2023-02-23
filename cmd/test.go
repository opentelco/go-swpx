package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	pb "git.liero.se/opentelco/go-swpx/proto/go/core"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	TestRootCmd.AddCommand(TestBulkCmd)

	TestBulkCmd.Flags().StringP("target", "t", "", "the target to test")
	if err := TestBulkCmd.MarkFlagRequired("target"); err != nil {
		panic(err)
	}

	TestBulkCmd.Flags().StringP("port-name", "n", "GigabitEthernet0/0/", "the port name to test")
	TestBulkCmd.Flags().Int("start", 1, "the first port to test")
	TestBulkCmd.Flags().Int("stop", 1, "the last port to test")
	TestBulkCmd.Flags().BoolP("concurrent", "c", false, "to run request concurrent or not")
	TestBulkCmd.Flags().String("ttl", "90s", "how long will we wait on each request")
	TestBulkCmd.Flags().StringSlice("providers", []string{""}, "specify which provider plugins to use")
	TestBulkCmd.Flags().StringP("resource", "r", "", "specify which resource plugin to use")

	Root.AddCommand(TestRootCmd)

}

// Test is a testing command that should be removed..
var TestRootCmd = &cobra.Command{
	Use:   "test",
	Short: "run some testing",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Test is a testing command that should be removed..
var TestBulkCmd = &cobra.Command{
	Use:   "bulk",
	Short: "run some bulk testing",
	Long:  `test is a command used under development to test libraries and other`,
	Run: func(cmd *cobra.Command, args []string) {

		target, _ := cmd.Flags().GetString("target")
		portName, _ := cmd.Flags().GetString("port-name")
		start, _ := cmd.Flags().GetInt("start")
		stop, _ := cmd.Flags().GetInt("stop")
		concurrent, _ := cmd.Flags().GetBool("concurrent")
		ttlString, _ := cmd.Flags().GetString("ttl")
		providers, _ := cmd.Flags().GetStringSlice("providers")
		resource, _ := cmd.Flags().GetString("resource")

		cmd.Println("selected providers: ", providers)

		_, err := time.ParseDuration(ttlString)
		if err != nil {
			cmd.Println("could not parse ttl: ", err)
			os.Exit(1)
		}

		conn, err := grpc.Dial("127.0.0.1:1338", grpc.WithInsecure())
		if err != nil {
			cmd.Println("could not connect to swpx: ", err)
			os.Exit(1)
		}
		swpx := pb.NewCoreClient(conn)
		wg := &sync.WaitGroup{}
		startTime := time.Now()

		for i := start; i <= stop; i++ {
			wg.Add(1)

			p := func(i int) *pb.Response {
				defer wg.Done()
				resp, err := swpx.Poll(cmd.Context(), &pb.Request{
					Settings: &pb.Request_Settings{
						ProviderPlugin: providers,
						ResourcePlugin: resource,
						RecreateIndex:  false,
						Timeout:        ttlString,
						CacheTtl:       "0s",
					},
					Type:     pb.Request_GET_BASIC_INFO,
					Hostname: target,
					Port:     fmt.Sprintf("%s%d", portName, i),
				})
				if err != nil {
					cmd.Printf("could not complete poll to %s (%s%d) reason: %s\n", target, portName, i, err)
				}
				if resp != nil {
					cmd.Printf("completed request %s (%s%d) in: %s (time since start: %s) \n", target, portName, i, resp.GetExecutionTime(), time.Since(startTime))
				}

				return resp

			}

			if concurrent {
				go p(i)
			} else {
				r := p(i)
				bs, _ := json.MarshalIndent(r, "", "  ")
				cmd.Println(string(bs))
			}
		}

		wg.Wait()
		cmd.Printf("completed all requests (%d) in: %s \n", stop, time.Since(startTime))
	},
}
