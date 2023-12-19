package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"go.opentelco.io/go-swpx/proto/go/corepb"
	"go.opentelco.io/go-swpx/proto/go/stanzapb"
	"google.golang.org/grpc"
)

func init() {

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
	TestRootCmd.AddCommand(TestBulkCmd)

	// get info
	TestDeviceInfo.Flags().StringP("target", "t", "", "the target to test")
	if err := TestDeviceInfo.MarkFlagRequired("target"); err != nil {
		panic(err)
	}
	TestDeviceInfo.Flags().String("ttl", "90s", "how long will we wait on each request")
	TestDeviceInfo.Flags().StringSlice("providers", []string{""}, "specify which provider plugins to use")
	TestDeviceInfo.Flags().StringP("resource", "r", "", "specify which resource plugin to use")
	TestRootCmd.AddCommand(TestDeviceInfo)

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

	diagnosticRunCmd.Flags().Int32("poll", 2, "how many times diagnostic should poll the device for data (10s between polls)")
	diagnosticRunCmd.Flags().StringP("target", "t", "", "the target to test")
	if err := diagnosticRunCmd.MarkFlagRequired("target"); err != nil {
		panic(err)
	}

	diagnosticRunCmd.Flags().StringP("port", "p", "GigabitEthernet0/0/1", "the port to check for")
	if err := diagnosticRunCmd.MarkFlagRequired("port"); err != nil {
		panic(err)
	}

	diagnosticRunCmd.Flags().StringSlice("providers", []string{""}, "specify which provider plugins to use")
	diagnosticRunCmd.Flags().StringP("resource", "r", "generic", "specify which resource plugin to use")
	diagnosticRunCmd.Flags().String("region", "default", "specify which region to poll the device in")
	diagnosticRunCmd.Flags().StringP("fingerprint", "f", "", "specify which fingerprint to use")

	diagnosticListCmd.Flags().Int64("limit", 10, "how many to fetch")
	diagnosticListCmd.Flags().Int64("offset", 0, "how many to skip")

	diagnosticCmd.AddCommand(diagnosticRunCmd)
	diagnosticCmd.AddCommand(diagnosticGetCmd)
	diagnosticCmd.AddCommand(diagnosticListCmd)

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
var TestDeviceInfo = &cobra.Command{
	Use:     "basic-device-info",
	Aliases: []string{"bdi"},
	Short:   "get basic device information from the poller",
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
		swpx := corepb.NewPollerClient(conn)

		resp, err := swpx.CollectBasicDeviceInformation(cmd.Context(), &corepb.CollectBasicDeviceInformationRequest{
			Settings: &corepb.Settings{
				ProviderPlugin: providers,
				ResourcePlugin: resource,
				RecreateIndex:  true,
				Timeout:        ttlString,
				CacheTtl:       "0s",
			},
			Session: &corepb.SessionRequest{
				Hostname: target,
			},
		})
		if err != nil {
			cmd.Printf("could not get device info for: %s,  reason: %s\n", target, err)
		}

		bs, _ := json.MarshalIndent(resp, "", "  ")
		cmd.Println(string(bs))
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
		swpx := corepb.NewPollerClient(conn)
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
		swpx := corepb.NewPollerClient(conn)

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
		swpx := corepb.NewPollerClient(conn)

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
			cmd.Println(prettyPrintJSON(resp.Device))
		} else {
			fmt.Println("failed to discover")
		}

	},
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "configure a device with a set of commands",
	Run: func(cmd *cobra.Command, args []string) {

		ttlString, _ := cmd.Flags().GetString("ttl")
		target, _ := cmd.Flags().GetString("target")
		providers, _ := cmd.Flags().GetStringSlice("providers")
		resource, _ := cmd.Flags().GetString("resource")
		region, _ := cmd.Flags().GetString("region")
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
		swpx := corepb.NewCommanderClient(conn)

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
				Hostname:      target,
				NetworkRegion: region,
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
	Short: "run and get diagnostic on network elements port",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var diagnosticRunCmd = &cobra.Command{
	Use:   "run",
	Short: "run diagnostic on network elements port",
	Run: func(cmd *cobra.Command, args []string) {

		target, _ := cmd.Flags().GetString("target")
		ttlString, _ := cmd.Flags().GetString("ttl")
		pollTimes, _ := cmd.Flags().GetInt32("poll")
		providers, _ := cmd.Flags().GetStringSlice("providers")
		resource, _ := cmd.Flags().GetString("resource")
		region, _ := cmd.Flags().GetString("resource")
		port, _ := cmd.Flags().GetString("port")
		fp, _ := cmd.Flags().GetString("fingerprint")

		if pollTimes < 2 {
			cmd.Println("poll must be 2 or more")
			os.Exit(1)
		}

		cmd.Println("selected providers: ", providers)

		conn, err := grpc.Dial("127.0.0.1:1338", grpc.WithInsecure())
		if err != nil {
			cmd.Println("could not connect to swpx: ", err)
			os.Exit(1)
		}
		swpx := corepb.NewPollerClient(conn)

		resp, err := swpx.RunDiagnostic(cmd.Context(), &corepb.RunDiagnosticRequest{
			Settings: &corepb.Settings{
				ProviderPlugin: providers,
				ResourcePlugin: resource,
				Timeout:        ttlString,
			},
			Session: &corepb.SessionRequest{
				Hostname:      target,
				Port:          port,
				NetworkRegion: region,
			},
			PollTimes:   pollTimes,
			Fingerprint: fp,
		})
		if err != nil {
			cmd.PrintErr(err)
		}
		if resp != nil {
			cmd.Println(prettyPrintJSON(resp))
		}

	},
}

var diagnosticGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "get diagnostic",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		id := args[0]
		if id == "" {
			cmd.Println("id cannot be empty")
			os.Exit(1)
		}

		conn, err := grpc.Dial("127.0.0.1:1338", grpc.WithInsecure())
		if err != nil {
			cmd.Println("could not connect to swpx: ", err)
			os.Exit(1)
		}
		swpx := corepb.NewPollerClient(conn)

		resp, err := swpx.GetDiagnostic(cmd.Context(), &corepb.GetDiagnosticRequest{
			Id: id,
		})
		if err != nil {
			cmd.PrintErr(err)
		}
		if resp != nil {
			cmd.Println(prettyPrintJSON(resp))
		}

	},
}

var diagnosticListCmd = &cobra.Command{
	Use:   "list [fingerprint]",
	Short: "list diagnostic",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		limit, _ := cmd.Flags().GetInt64("limit")
		offset, _ := cmd.Flags().GetInt64("offset")

		id := args[0]
		if id == "" {
			cmd.Println("fingerprint cannot be empty")
			os.Exit(1)
		}

		conn, err := grpc.Dial("127.0.0.1:1338", grpc.WithInsecure())
		if err != nil {
			cmd.Println("could not connect to swpx: ", err)
			os.Exit(1)
		}
		swpx := corepb.NewPollerClient(conn)

		params := &corepb.ListDiagnosticsRequest{
			Fingerprint: id,
		}
		if limit != 0 {
			params.Limit = limit
		}
		if offset != 0 {
			params.Offset = offset
		}

		resp, err := swpx.ListDiagnostics(cmd.Context(), params)
		if err != nil {
			cmd.PrintErr(err)
		}
		if resp != nil {
			cmd.Println(prettyPrintJSON(resp))
		}

	},
}
