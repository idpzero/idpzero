package server

import "net/http"

func register(server *Server) {

	server.mux.HandleFunc("/", handleLanding())
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
