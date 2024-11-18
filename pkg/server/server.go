package server

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/idpzero/idpzero/pkg/console"
	"github.com/idpzero/idpzero/pkg/dbg"
	"github.com/idpzero/idpzero/pkg/store/query"
	"github.com/idpzero/idpzero/pkg/web/controllers"
	"github.com/savioxavier/termlink"
)

type Server struct {
	server  *http.Server
	waiter  sync.WaitGroup
	lock    sync.RWMutex
	logger  *slog.Logger
	config  *configuration.ServerConfig
	users   *users
	queries *query.Queries
}

func NewServer(logger *slog.Logger, config *configuration.ConfigurationManager, queries *query.Queries) (*Server, error) {

	// Use chi as this is what OIDC is using internally, so keep it conistent
	router := chi.NewRouter()

	c, err := config.LoadServer()

	if err != nil {
		return nil, err
	}

	server := &Server{
		waiter:  sync.WaitGroup{},
		lock:    sync.RWMutex{},
		logger:  logger,
		server:  &http.Server{Addr: fmt.Sprintf("0.0.0.0:%d", c.Server.Port), Handler: router},
		users:   newUsers(),
		queries: queries,
	}

	// set the config using the initial data provided.
	server.setConfig(c)

	// update for changes to be loaded
	config.OnServerChanged(server.setConfig)

	router.Use(middleware.RequestID)
	router.Use(setProviderFromRequest) // set the issuer based on the request URL
	router.Use(middleware.Recoverer)

	// create the storage provider
	storage, err := NewStorage(logger, config, queries, server.users)

	if err != nil {
		return nil, err
	}

	options := ProviderOptions{
		Storage: storage,
	}

	provider, err := NewProvider(logger, options)
	if err != nil {
		return nil, err
	}

	// we need to add a route to the root because we  are mounting
	// the provider on the root, we cant double map the '/'
	rtr := provider.Handler.(*chi.Mux)
	controllers.Routes(rtr, func() *configuration.ServerConfig {
		return server.config
	}, queries, provider)

	//rtr.Get("/", )

	router.Mount("/", provider)

	return server, nil
}

func (s *Server) setConfig(config *configuration.ServerConfig) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.config = config
	s.users.Update(config.Users)
}

const signingKeyID = "signing"

func (s *Server) onStart(ctx context.Context) error {

	key, err := s.queries.GetKeyByID(ctx, signingKeyID)

	if err != nil {
		if err == sql.ErrNoRows {
			console.PrintCheck(console.IconQuestion, "No signing key found. Generating a new one.")
			keyArgs, err := newRSAKey(signingKeyID, KeyUseSig)

			if err != nil {
				return err
			}

			key, err = s.queries.CreateKey(ctx, keyArgs)

			if err != nil {
				return err
			}

			console.PrintCheck(console.IconCheck, "Signing key generated successfully.")

		} else {
			return err
		}
	} else {
		if key.Usage == KeyUseSig {
			console.PrintCheck(console.IconCheck, "Existing signing key found OK")
		} else {
			console.PrintCheck(console.IconCross, "Expecing a signing key, but found a different use key with ID '%s'", key.ID)
			return fmt.Errorf("key with ID '%s' is not a signing key", key.ID)
		}
	}

	return nil
}

func (s *Server) Run(ctx context.Context) error {

	// setup local state as needed.
	if err := s.onStart(ctx); err != nil {
		return err
	}

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
