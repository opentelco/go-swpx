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

package main

import (
	"fmt"
	"os"

	"git.liero.se/opentelco/go-swpx/cmd"
)

//go:generate protoc -I /usr/local/include/ -I $GOPATH/src/ -I ./proto/src --go-grpc_out=$GOPATH/src/ --go_out=$GOPATH/src/ ./proto/src/healtcheck.proto
//go:generate protoc -I /usr/local/include/ -I $GOPATH/src/ -I ./proto/src --go-grpc_out=$GOPATH/src/ --go_out=$GOPATH/src/ ./proto/src/resource.proto
//go:generate protoc -I /usr/local/include/ -I $GOPATH/src/ -I ./proto/src --go-grpc_out=$GOPATH/src/ --go_out=$GOPATH/src/ ./proto/src/provider.proto
//go:generate protoc -I /usr/local/include/ -I $GOPATH/src/ -I ./proto/src --go-grpc_out=$GOPATH/src/ --go_out=$GOPATH/src/ ./proto/src/dnc.proto
//go:generate protoc -I /usr/local/include/ -I $GOPATH/src/ -I ./proto/src --go-grpc_out=$GOPATH/src/ --go_out=$GOPATH/src/ ./proto/src/core.proto
//go:generate protoc -I /usr/local/include/ -I $GOPATH/src/ -I ./proto/src --go-grpc_out=$GOPATH/src/ --go_out=$GOPATH/src/ ./proto/src/network_element.proto
//go:generate protoc -I /usr/local/include/ -I $GOPATH/src/ -I ./proto/src --go-grpc_out=$GOPATH/src/ --go_out=$GOPATH/src/ ./proto/src/traffic_policy.proto
//go:generate protoc -I /usr/local/include/ -I $GOPATH/src/ -I ./proto/src --go-grpc_out=$GOPATH/src/ --go_out=$GOPATH/src/ ./proto/src/analysis.proto
//go:generate protoc -I /usr/local/include/ -I $GOPATH/src/ -I ./proto/src --go-grpc_out=$GOPATH/src/ --go_out=$GOPATH/src/ ./proto/src/fleet.proto
//go:generate protoc -I /usr/local/include/ -I $GOPATH/src/ -I ./proto/src --go-grpc_out=$GOPATH/src/ --go_out=$GOPATH/src/ ./proto/src/fleet_device.proto
//go:generate protoc -I /usr/local/include/ -I $GOPATH/src/ -I ./proto/src --go-grpc_out=$GOPATH/src/ --go_out=$GOPATH/src/ ./proto/src/fleet_configuration.proto
//go:generate protoc -I /usr/local/include/ -I $GOPATH/src/ -I ./proto/src --go-grpc_out=$GOPATH/src/ --go_out=$GOPATH/src/ ./proto/src/fleet_template.proto
//go:generate go run ./proto/tagparser/tagparser.go ./proto/go

func main() {
	if err := cmd.Root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
