/*
 * Copyright (c) 2020. Liero AB
 *
 * Permission is hereby granted, free of charge, to any person obtaining
 * a copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software
 * is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
 * OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
 * IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
 * CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
 * TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"

	"git.liero.se/opentelco/go-swpx/api"
	"git.liero.se/opentelco/go-swpx/core"
	pb "git.liero.se/opentelco/go-swpx/proto/go/core"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var httpPort string
var grpcPort string
var logger hclog.Logger

func init() {
	Root.AddCommand(Version)
	Root.AddCommand(Start)

	TestRootCmd.AddCommand(TestBulkCmd)

	TestBulkCmd.Flags().StringP("target", "t", "", "the target to test")
	if err := TestBulkCmd.MarkFlagRequired("target"); err != nil {
		panic(err)
	}

	TestBulkCmd.Flags().StringP("port-name", "n", "GigabitEthernet0/0/", "the port name to test")
	TestBulkCmd.Flags().Int("start", 1, "the first port to test")
	TestBulkCmd.Flags().Int("stop", 1, "the last port to test")
	TestBulkCmd.Flags().BoolP("concurrent", "c", false, "to run request concurrent or not")

	Root.AddCommand(TestRootCmd)

	Start.Flags().StringVarP(&httpPort, "port", "p", "1337", "the port to use for http")
	Start.Flags().StringVarP(&grpcPort, "grpc-port", "g", "1338", "the port to use for grpc")
}

const APP_NAME = "go-swpx"

var Root = &cobra.Command{Use: APP_NAME}

// Start starts Switchpoller daemon/server
var Start = &cobra.Command{
	Use:   "start",
	Short: "start the swpx daemon",
	Long:  `switchpoller x. the long description of the application`,
	Run: func(cmd *cobra.Command, args []string) {
		logger = hclog.New(&hclog.LoggerOptions{
			Name:            APP_NAME,
			Level:           hclog.Debug,
			Color:           hclog.AutoColor,
			IncludeLocation: true,
		})

		c, err := core.New(logger)
		if err != nil {
			panic(err)
		}
		if err := c.Start(); err != nil {
			panic(err)
		}

		// start API endpoint and add the queue
		// the queue is initated in the core and n workers takes request from it.

		// HTTP
		server := api.NewServer(c, logger)
		go func() {
			err = server.ListenAndServe(":" + httpPort)
			if err != nil {
				panic(err)
			}
		}()
		// GRPC
		grpcServer := api.NewCoreGrpcServer(c, logger)
		go func() {
			err = grpcServer.ListenAndServe(":" + grpcPort)
			if err != nil {
				panic(err)
			}
		}()

		signalChan := make(chan os.Signal, 1)

		signal.Notify(
			signalChan,
			syscall.SIGHUP,  // kill -SIGHUP XXXX
			syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
			syscall.SIGQUIT, // kill -SIGQUIT XXXX
		)

		<-signalChan
		cmd.Println("os.Interrupt - shutting down...")

		go func() {
			<-signalChan
			cmd.Println("os.Kill - terminating...")
			os.Exit(1)
		}()

		// manually cancel context if not using httpServer.RegisterOnShutdown(cancel)

		defer os.Exit(0)
	},
}

// Version returns the version of SWP-X
var Version = &cobra.Command{
	Use:   "version",
	Short: "print the version number of SWPX",
	Long:  `all software has versions. This is SwitchpollerX's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SwitchPoller X - Poller and Commander - ", core.VERSION)
	},
}

type eventRequest struct {
	ID    uuid.UUID
	Key   string
	Value string
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

			p := func(i int) {
				resp, err := swpx.Poll(cmd.Context(), &pb.Request{
					Settings: &pb.Request_Settings{
						ProviderPlugin:         []string{"default_provider"},
						ResourcePlugin:         "vrp",
						RecreateIndex:          false,
						DisableDistributedLock: false,
						Timeout:                "90s",
						CacheTtl:               "0s",
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

				wg.Done()

			}

			if concurrent {
				go p(i)
			} else {
				p(i)
			}
		}

		wg.Wait()
		cmd.Printf("completed all requests (%d) in: %s \n", stop, time.Since(startTime))
	},
}
