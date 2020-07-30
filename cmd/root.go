package cmd

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"os"
	"time"

	"git.liero.se/opentelco/go-swpx/api"
	"git.liero.se/opentelco/go-swpx/core"
	"git.liero.se/opentelco/go-swpx/core/requestcache"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var port string

func init() {
	Root.AddCommand(Version)
	Root.AddCommand(Start)
	Root.AddCommand(Test)

	Start.Flags().StringVarP(&port, "port", "p", "1337", "the port to use")
}

const APP_NAME = "go-swpx"

var Root = &cobra.Command{Use: APP_NAME}

// Start starts Switchpoller daemon/server
var Start = &cobra.Command{
	Use:   "start",
	Short: "start the swpx daemon",
	Long:  `switchpoller x. the long description of the application`,
	Run: func(cmd *cobra.Command, args []string) {
		c := core.CreateCore()
		c.Start()

		// start API endpoint and add the queue
		// the queue is initated in the core and n workers takes request from it.
		server := api.New(core.RequestQueue)
		err := server.ListenAndServe(":" + port)
		if err != nil {
			panic(err)
		}
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
var Test = &cobra.Command{
	Use:   "test",
	Short: "run some testing",
	Long:  `test is a command used under development to test libraries and other`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := hclog.New(&hclog.LoggerOptions{
			Name:   APP_NAME,
			Output: os.Stdout,
			Level:  hclog.Debug,
		})

		orchester := requestcache.New()
		id, _ := uuid.NewUUID()

		go func() {
			for {
				time.Sleep(2 * time.Second)
				x, err := orchester.Pop(id)
				if err != nil {
					logger.Error(err.Error())
					continue
				}
				x <- &eventRequest{ID: id, Key: "some string", Value: "some value"}
			}
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		println("handling request:", id.String())

		responseChan := orchester.Put(ctx, id)
		go func() {
			for {
				select {
				case resp := <-responseChan:
					if x, ok := resp.(*eventRequest); ok {
						logger.Info("id:", x.ID.String())
						logger.Info("key:", x.Key)
						logger.Info("value:", x.Value)
					} else {
						logger.Warn("the returned type could not be converted")
					}
					logger.Info("size of cache:", orchester.GetSize())
					return
				case <-ctx.Done():
					println("time out waiting for response")
					orchester.Delete(id)
					logger.Info("Size of cache:", orchester.GetSize())
				}
			}
		}()
		time.Sleep(30 * time.Second)
	},
}
