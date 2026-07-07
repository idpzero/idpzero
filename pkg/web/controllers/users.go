package controllers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/idpzero/idpzero/pkg/web/models"
	"github.com/idpzero/idpzero/pkg/web/views/pages"
)

func users(config func() *configuration.ServerConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		idpConfig := config()

		im := models.NewUsersModel(idpConfig.Users)

		view := pages.UsersView(im)

		templ.Handler(view).ServeHTTP(w, r)
	})
}
