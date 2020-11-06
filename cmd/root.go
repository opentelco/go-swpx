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
var logger hclog.Logger

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
		logger = hclog.New(&hclog.LoggerOptions{
			Name:   APP_NAME,
			Level:  hclog.Debug,
			Color: hclog.AutoColor,
			IncludeLocation: true,
		})
		
		c,err := core.New(logger)
		if err != nil {
			panic(err)
		}
		if err := c.Start(); err != nil {
			panic(err)
		}

		// start API endpoint and add the queue
		// the queue is initated in the core and n workers takes request from it.

		server := api.NewServer(core.RequestQueue)
		//server := api.NewGRPCServer(core.RequestQueue)

		err = server.ListenAndServe(":" + port)
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
