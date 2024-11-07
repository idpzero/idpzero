package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/idpzero/idpzero/pkg/dbg"
	"github.com/idpzero/idpzero/pkg/idp"
	"github.com/idpzero/idpzero/pkg/web/handlers"
	"github.com/savioxavier/termlink"
)

type Server struct {
	server *http.Server
	waiter sync.WaitGroup
	lock   sync.RWMutex
	logger *slog.Logger
	config *configuration.ServerConfig
}

func NewServer(logger *slog.Logger, config *configuration.ConfigurationManager, storage *idp.Storage) (*Server, error) {

	// Use chi as this is what OIDC is using internally, so keep it conistent
	router := chi.NewRouter()

	c, err := config.LoadServer()

	if err != nil {
		return nil, err
	}

	server := &Server{
		waiter: sync.WaitGroup{},
		lock:   sync.RWMutex{},
		logger: logger,
		config: c,
		server: &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", c.Server.Port), Handler: router},
	}

	// update for changes to be loaded
	config.OnServerChanged(server.setConfig)

	router.Use(middleware.RequestID)
	router.Use(setProviderFromRequest) // set the issuer based on the request URL
	router.Use(middleware.Recoverer)

	options := idp.ProviderOptions{
		Storage: storage,
	}

	provider, err := idp.NewProvider(logger, options)
	if err != nil {
		return nil, err
	}

	// we need to add a route to the root because we  are mounting
	// the provider on the root, we cant double map the '/'
	rtr := provider.Handler.(*chi.Mux)
	handlers.Routes(rtr, func() *configuration.ServerConfig {
		return server.config
	})

	//rtr.Get("/", )

	router.Mount("/", provider)

	return server, nil
}

func (s *Server) setConfig(config *configuration.ServerConfig) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.config = config
}

func (s *Server) Run(ctx context.Context) error {

	// server context so we can control the shutdown order
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	go func() {
		<-ctx.Done()
		fmt.Println()
		fmt.Println("Initiating shutting down...")
		if err := s.server.Shutdown(context.Background()); err != nil {
			dbg.Logger.Error("Error shutting down server", "error", err.Error())
		}
		serverStopCtx()
	}()

	hosted := fmt.Sprintf("http://localhost:%d", s.config.Server.Port)

	fmt.Println(
		"Identity Provider started at",
		color.CyanString(termlink.Link(hosted, hosted)),
	)

	// Run the server
	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	<-serverCtx.Done() // wait for shutdown to complete
	color.Green("Shutdown complete. Bye!")

	return nil
}
