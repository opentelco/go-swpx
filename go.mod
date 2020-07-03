module git.liero.se/opentelco/go-swpx

go 1.14

require (
	git.liero.se/opentelco/go-dnc v0.0.0-20200624104445-9445e837da36
	github.com/go-chi/chi v4.1.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/gorilla/context v1.1.1
	github.com/hashicorp/go-hclog v0.0.0-20180709165350-ff2cf002a8dd
	github.com/hashicorp/go-plugin v1.3.0
	github.com/hashicorp/go-version v1.2.1
	github.com/hashicorp/yamux v0.0.0-20180604194846-3520598351bb // indirect
	github.com/jhump/protoreflect v1.6.0 // indirect
	github.com/mitchellh/go-testing-interface v1.0.0 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/oklog/run v1.0.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/segmentio/ksuid v1.0.2
	github.com/spf13/cobra v1.0.0
	golang.org/x/net v0.0.0-20190620200207-3b0461eec859 // indirect
	google.golang.org/grpc v1.27.1
	google.golang.org/protobuf v1.23.0

)

replace git.liero.se/opentelco/go-dnc v0.0.0 => ../go-dnc
