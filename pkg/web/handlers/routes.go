package handlers

import (
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/idpzero/idpzero/pkg/store/query"
	"github.com/idpzero/idpzero/pkg/web/assets"
	"github.com/zitadel/oidc/v3/pkg/op"
)

func mustFs(fs fs.FS, err error) fs.FS {
	if err != nil {
		panic(err)
	}
	return fs
}

func Routes(router *chi.Mux, config func() *configuration.ServerConfig, query *query.Queries, provider op.OpenIDProvider) {

	router.Handle("/static/*", http.FileServer(http.FS(assets.Static)))

	// manage the favicons
	favhanlder := http.FileServer(http.FS(mustFs(fs.Sub(assets.Static, "static/favicon"))))
	router.Handle("/android-chrome-192x192.png", favhanlder)
	router.Handle("/android-chrome-512x512.png", favhanlder)
	router.Handle("/apple-touch-icon.png", favhanlder)
	router.Handle("/favicon-16x16.png", favhanlder)
	router.Handle("/favicon-32x32.png", favhanlder)
	router.Handle("/favicon.ico", favhanlder)
	router.Handle("/site.webmanifest", favhanlder)

	router.Get("/", index(config))
	router.Get("/login", userlogin(config, query))
	router.Post("/login", userloginSubmit(config, query, op.AuthCallbackURL(provider)))

}
