package graph

import (
	"net"
	"net/http"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/hashicorp/go-hclog"
)

func NewResolver(device devicepb.DeviceServiceServer) Config {
	return Config{
		Resolvers: &Resolver{
			devices: device,
		},
	}
}

func Serve(listner net.Listener, config Config, logger hclog.Logger) error {

	srv := handler.NewDefaultServer(NewExecutableSchema(config))

	// srv.SetRecoverFunc(func(ctx context.Context, err interface{}) (userMessage error) {
	// 	// send this panic somewhere
	// 	log.Print(err)
	// 	debug.PrintStack()
	// 	return errors.New("user message on panic")
	// })

	http.Handle("/", playground.Handler("Fleet", "/gql"))
	http.Handle("/gql", srv)

	logger.Info("serving graphql server", "addr", listner.Addr())
	return http.Serve(listner, srv)
}
