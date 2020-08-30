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

//go:generate protoc -I /usr/local/include/ -I $GOPATH/src/ -I $GOPATH/src/git.liero.se/opentelco/go-swpx/proto/ --go_out=plugins=grpc:$GOPATH/src/ $GOPATH/src/git.liero.se/opentelco/go-swpx/proto/health/healtcheck.proto
//go:generate protoc -I /usr/local/include/ -I $GOPATH/src/ -I $GOPATH/src/git.liero.se/opentelco/go-swpx/proto/ --go_out=plugins=grpc:$GOPATH/src/ $GOPATH/src/git.liero.se/opentelco/go-swpx/proto/resource/resource.proto
//go:generate protoc -I /usr/local/include/ -I $GOPATH/src/ -I $GOPATH/src/git.liero.se/opentelco/go-swpx/proto/ --go_out=plugins=grpc:$GOPATH/src/ $GOPATH/src/git.liero.se/opentelco/go-swpx/proto/provider/provider.proto
//go:generate protoc -I /usr/local/include/ -I $GOPATH/src/ -I $GOPATH/src/git.liero.se/opentelco/go-swpx/proto/ --go_out=plugins=grpc:$GOPATH/src/ $GOPATH/src/git.liero.se/opentelco/go-swpx/proto/dnc/dnc.proto
//go:generate protoc -I /usr/local/include/ -I $GOPATH/src/ -I $GOPATH/src/git.liero.se/opentelco/go-swpx/proto/ --go_out=plugins=grpc:$GOPATH/src/ $GOPATH/src/git.liero.se/opentelco/go-swpx/proto/networkelement/element.proto

func main() {
	if err := cmd.Root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
