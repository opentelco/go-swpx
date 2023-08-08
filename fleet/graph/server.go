package graph

import (
	"net"
	"net/http"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/notificationpb"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/hashicorp/go-hclog"
)

func NewResolver(device devicepb.DeviceServiceServer, notification notificationpb.NotificationServiceServer) Config {
	return Config{
		Resolvers: &Resolver{
			devices:       device,
			notifications: notification,
		},
	}
}

func Serve(listner net.Listener, config Config, logger hclog.Logger) error {

	srv := handler.NewDefaultServer(NewExecutableSchema(config))
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"http://localhost:8080", "https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// srv.SetRecoverFunc(func(ctx context.Context, err interface{}) (userMessage error) {
	// 	// send this panic somewhere
	// 	log.Print(err)
	// 	debug.PrintStack()
	// 	return errors.New("user message on panic")
	// })

	r.Handle("/", playground.Handler("Fleet", "/gql"))
	r.Handle("/gql", srv)

	logger.Info("serving graphql server", "addr", listner.Addr())
	return http.Serve(listner, srv)
}
