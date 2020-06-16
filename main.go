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
