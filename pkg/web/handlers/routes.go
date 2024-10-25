package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/idpzero/idpzero/pkg/web"
)

func Routes(router *chi.Mux) {
	router.Handle("/assets/*", http.FileServer(http.FS(web.Assets)))

	router.Get("/", index())
}
