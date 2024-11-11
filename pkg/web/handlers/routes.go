package handlers

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/idpzero/idpzero/pkg/web"
)

func mustFs(fs fs.FS, err error) fs.FS {
	if err != nil {
		panic(err)
	}
	return fs
}

func Routes(router *chi.Mux, config func() *configuration.ServerConfig) {

	router.Handle("/assets/*", http.FileServer(http.FS(web.Assets)))

	// manage the favicons
	favhanlder := http.FileServer(http.FS(mustFs(fs.Sub(web.Assets, "assets/favicon"))))
	router.Handle("/android-chrome-192x192.png", favhanlder)
	router.Handle("/android-chrome-512x512.png", favhanlder)
	router.Handle("/apple-touch-icon.png", favhanlder)
	router.Handle("/favicon-16x16.png", favhanlder)
	router.Handle("/favicon-32x32.png", favhanlder)
	router.Handle("/favicon.ico", favhanlder)
	router.Handle("/site.webmanifest", favhanlder)

	router.Get("/", index(config))
	router.Get("/login", userlogin(config))
	router.Post("/login", userloginSubmit(config))

}
