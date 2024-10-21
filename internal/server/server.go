package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/idpzero/idpzero/internal/idp"
)

type Server struct {
	server *http.Server
	waiter sync.WaitGroup
	lock   sync.RWMutex
	logger *slog.Logger
	config idp.IDPConfiguration
}

func NewServer(logger *slog.Logger, config idp.IDPConfiguration, storage *idp.Storage) (*Server, error) {

	// Use chi as this is what OIDC is using internally, so keep it conistent
	router := chi.NewRouter()

	server := &Server{
		waiter: sync.WaitGroup{},
		lock:   sync.RWMutex{},
		logger: logger,
		config: config,
		server: &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", config.Server.Port), Handler: router},
	}

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	options := idp.ProviderOptions{
		Storage: storage,
		Issuer:  config.Server.Issuer,
	}

	provider, err := idp.NewProvider(logger, options)
	if err != nil {
		return nil, err
	}

	router.Mount("/", provider)
	//router.Handle("/", provider)

	return server, nil
}

func (s *Server) UpdateConfig(config idp.IDPConfiguration) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.config = config
}

func (s *Server) Run(ctx context.Context) error {

	// server context so we can control the shutdown order
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	go func() {
		<-ctx.Done()
		s.server.Shutdown(context.Background())
		serverStopCtx()
	}()

	// Run the server
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	<-serverCtx.Done() // wait for shutdown to complete

	return nil
}
