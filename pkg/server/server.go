package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type Server struct {
	wg       *sync.WaitGroup
	mux      *http.ServeMux
	server   *http.Server
	logger   *slog.Logger
	provider op.OpenIDProvider
}

func NewServer(logger *slog.Logger, conf *configuration.Document, store op.Storage) (*Server, error) {

	if conf == nil {
		return nil, fmt.Errorf("config is nil")
	}

	if conf.Server.Port == 0 {
		return nil, fmt.Errorf("field 'Server.Port' is not set in configuration")
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.Port),
		Handler: mux,
	}

	p, err := setupProvider(store, conf.Server.Issuer, conf.Server.Key(), logger)

	if err != nil {
		return nil, err
	}

	svr := &Server{
		logger:   logger,
		server:   server,
		mux:      mux,
		wg:       &sync.WaitGroup{},
		provider: p,
	}

	// setup the routes
	register(svr, p)

	return svr, nil
}

func (s *Server) Start() {
	go runAndWait(s)
	s.logger.Info(fmt.Sprintf("Server listening on '%s'", s.server.Addr))
}

func (s *Server) Shutdown(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	s.wg.Wait() // wait for the shutdown to complete
	if err != nil {
		return err
	}

	s.logger.Debug("Server Shutdown OK")

	return nil
}

func runAndWait(s *Server) error {
	s.wg.Add(1)
	err := s.server.ListenAndServe()
	s.wg.Done()

	if err != http.ErrServerClosed {
		return err
	}

	return nil
}
