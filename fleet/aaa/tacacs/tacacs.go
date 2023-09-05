package tacacs

import (
	"context"
	"net"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/nwaples/tacplus"
)

func New() {
	l, err := net.Listen("tcp", ":49")
	if err != nil {
		panic(err)
	}
	secret := []byte("secret")

	handler := &tacplus.ServerConnHandler{
		Handler: &handler{
			logger: hclog.New(&hclog.LoggerOptions{
				Name:  "tacacs",
				Level: hclog.Debug,
				Color: hclog.AutoColor,
			}),
		},
		ConnConfig: tacplus.ConnConfig{
			Mux:    true,
			Secret: secret,
		},
	}

	srv := tacplus.Server{
		ServeConn: func(nc net.Conn) {
			handler.Serve(nc)

		},
	}
	err = srv.Serve(l)
	if err != nil {
		panic(err)
	}

}

type handler struct {
	logger hclog.Logger
}

func (h *handler) HandleAuthenStart(ctx context.Context, a *tacplus.AuthenStart, s *tacplus.ServerSession) *tacplus.AuthenReply {
	h.logger.Debug("HandleAuthenStart", "AuthenStart", a)
	return nil
}
func (h *handler) HandleAuthorRequest(ctx context.Context, a *tacplus.AuthorRequest, s *tacplus.ServerSession) *tacplus.AuthorResponse {
	h.logger.Debug("HandleAuthorRequest", "AuthorRequest", a)
	h.logger.Debug("arguments", strings.Join(a.Arg, "\n"))
	return &tacplus.AuthorResponse{
		Status: tacplus.AuthenStatusPass,
	}
}
func (h *handler) HandleAcctRequest(ctx context.Context, a *tacplus.AcctRequest, s *tacplus.ServerSession) *tacplus.AcctReply {
	h.logger.Debug("HandleAcctRequest", "AuthenService", a.AuthenService, "AuthenMethod")
	h.logger.Debug("arguments", strings.Join(a.Arg, "\n"))

	return nil
}
