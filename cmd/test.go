package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"git.liero.se/opentelco/go-swpx/proto/go/corepb"
	"git.liero.se/opentelco/go-swpx/proto/go/stanzapb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
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

	// collect
	collectConfigCmd.Flags().String("ttl", "90s", "how long will we wait on each request")
	collectConfigCmd.Flags().StringP("target", "t", "", "the target to test")
	if err := collectConfigCmd.MarkFlagRequired("target"); err != nil {
		panic(err)
	}
	collectConfigCmd.Flags().StringSlice("providers", []string{"default"}, "specify which provider plugins to use")
	collectConfigCmd.Flags().StringP("resource", "r", "vrp", "specify which resource plugin to use")
	Root.AddCommand(collectConfigCmd)

	discoverCmd.Flags().String("ttl", "90s", "how long will we wait on each request")
	discoverCmd.Flags().StringP("target", "t", "", "the target to test")
	if err := discoverCmd.MarkFlagRequired("target"); err != nil {
		panic(err)
	}
	discoverCmd.Flags().StringSlice("providers", []string{""}, "specify which provider plugins to use")
	discoverCmd.Flags().StringP("resource", "r", "generic", "specify which resource plugin to use")

	TestRootCmd.AddCommand(discoverCmd)

	// -

	configureCmd.Flags().String("ttl", "90s", "how long will we wait on each request")
	configureCmd.Flags().StringP("target", "t", "", "the target to test")
	if err := configureCmd.MarkFlagRequired("target"); err != nil {
		panic(err)
	}
	configureCmd.Flags().StringSlice("providers", []string{""}, "specify which provider plugins to use")
	configureCmd.Flags().StringP("resource", "r", "generic", "specify which resource plugin to use")
	configureCmd.Flags().StringSliceP("line", "l", []string{}, "specify which line to configure")
	if err := configureCmd.MarkFlagRequired("line"); err != nil {
		panic(err)
	}
	TestRootCmd.AddCommand(configureCmd)

	// -

	diagnosticCmd.Flags().String("ttl", "20s", "how long the diagnostc will run")
	diagnosticCmd.Flags().StringP("target", "t", "", "the target to test")
	if err := diagnosticCmd.MarkFlagRequired("target"); err != nil {
		panic(err)
	}
	diagnosticCmd.Flags().StringSlice("providers", []string{""}, "specify which provider plugins to use")
	diagnosticCmd.Flags().StringP("resource", "r", "generic", "specify which resource plugin to use")
	TestRootCmd.AddCommand(diagnosticCmd)

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
		swpx := corepb.NewCoreServiceClient(conn)
		wg := &sync.WaitGroup{}
		startTime := time.Now()

		for i := start; i <= stop; i++ {
			wg.Add(1)

			p := func(i int) *corepb.PollResponse {
				defer wg.Done()
				resp, err := swpx.Poll(cmd.Context(), &corepb.PollRequest{
					Settings: &corepb.Settings{
						ProviderPlugin: providers,
						ResourcePlugin: resource,
						RecreateIndex:  false,
						Timeout:        ttlString,
						CacheTtl:       "0s",
					},
					Type: corepb.PollRequest_GET_BASIC_INFO,
					Session: &corepb.SessionRequest{
						Hostname: target,
						Port:     fmt.Sprintf("%s%d", portName, i),
					},
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

var collectConfigCmd = &cobra.Command{
	Use:   "collect-config",
	Short: "get running config from network element",
	Run: func(cmd *cobra.Command, args []string) {

		target, _ := cmd.Flags().GetString("target")
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
		swpx := corepb.NewCoreServiceClient(conn)

		resp, err := swpx.CollectConfig(cmd.Context(), &corepb.CollectConfigRequest{
			Settings: &corepb.Settings{
				ProviderPlugin: providers,
				ResourcePlugin: resource,
				Timeout:        ttlString,
			},
			Session: &corepb.SessionRequest{
				Hostname: target,
			},
		})
		if err != nil {
			cmd.PrintErr(err)
		}
		if resp != nil {
			fmt.Println("Collected config")
			fmt.Println(resp.GetConfig())
		} else {
			fmt.Println("Failed to collect config")
		}

	},
}

var discoverCmd = &cobra.Command{
	Use:   "discover",
	Short: "discover the device and return basic info",
	Run: func(cmd *cobra.Command, args []string) {

		target, _ := cmd.Flags().GetString("target")
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
		swpx := corepb.NewCoreServiceClient(conn)

		resp, err := swpx.Discover(cmd.Context(), &corepb.DiscoverRequest{
			Settings: &corepb.Settings{
				ProviderPlugin: providers,
				ResourcePlugin: resource,
				Timeout:        ttlString,
			},
			Session: &corepb.SessionRequest{
				Hostname: target,
			},
		})
		if err != nil {
			cmd.PrintErr(err)
		}
		if resp != nil {
			cmd.Println(prettyPrintJSON(resp.NetworkElement))
		} else {
			fmt.Println("failed to discover")
		}

	},
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "configure a device with a set of commands",
	Run: func(cmd *cobra.Command, args []string) {

		target, _ := cmd.Flags().GetString("target")
		ttlString, _ := cmd.Flags().GetString("ttl")
		providers, _ := cmd.Flags().GetStringSlice("providers")
		resource, _ := cmd.Flags().GetString("resource")
		lines, _ := cmd.Flags().GetStringSlice("line")

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
		swpx := corepb.NewCommanderServiceClient(conn)

		stanza := make([]*stanzapb.ConfigurationLine, len(lines))
		for i, line := range lines {
			stanza[i] = &stanzapb.ConfigurationLine{
				Content: line,
			}
		}

		resp, err := swpx.ConfigureStanza(cmd.Context(), &corepb.ConfigureStanzaRequest{
			Settings: &corepb.Settings{
				ProviderPlugin: providers,
				ResourcePlugin: resource,
				Timeout:        ttlString,
			},
			Session: &corepb.SessionRequest{
				Hostname: target,
			},
			Stanza: stanza,
		})
		if err != nil {
			cmd.PrintErr(err)
		}
		if resp != nil {
			cmd.Println(prettyPrintJSON(resp))
		} else {
			fmt.Println("failed to configure")
		}

	},
}

var diagnosticCmd = &cobra.Command{
	Use:   "diagnostic",
	Short: "run diagnostic on network elements port",
	Run: func(cmd *cobra.Command, args []string) {

		target, _ := cmd.Flags().GetString("target")
		ttlString, _ := cmd.Flags().GetString("ttl")
		providers, _ := cmd.Flags().GetStringSlice("providers")
		resource, _ := cmd.Flags().GetString("resource")

		cmd.Println("selected providers: ", providers)

		ttlDur, err := time.ParseDuration(ttlString)
		if err != nil {
			cmd.Println("could not parse ttl: ", err)
			os.Exit(1)
		}

		conn, err := grpc.Dial("127.0.0.1:1338", grpc.WithInsecure())
		if err != nil {
			cmd.Println("could not connect to swpx: ", err)
			os.Exit(1)
		}
		swpx := corepb.NewCoreServiceClient(conn)

		resp, err := swpx.Diagnostic(cmd.Context(), &corepb.DiagnosticRequest{
			Settings: &corepb.Settings{
				ProviderPlugin: providers,
				ResourcePlugin: resource,
				Timeout:        ttlString,
			},
			Session: &corepb.SessionRequest{
				Hostname: target,
			},
			DiagnosicLength: durationpb.New(ttlDur),
		})
		if err != nil {
			cmd.PrintErr(err)
		}
		if resp != nil {
			cmd.Println(prettyPrintJSON(resp))
		} else {
			fmt.Println("failed to configure")
		}

	},
}
