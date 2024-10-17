package server

import (
	"net/http"

	"github.com/zitadel/oidc/v3/pkg/op"
)

func register(server *Server, provider op.OpenIDProvider) {

	server.mux.Handle("/", provider)
	server.mux.HandleFunc("/foo", handleLanding())
	server.mux.HandleFunc(pathLoggedOut, handleLoggedOut())
}

func handleLoggedOut() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Logged out"))
	}
}

func handleLanding() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Landing"))
	}
}
