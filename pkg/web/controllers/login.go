package controllers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/idpzero/idpzero/pkg/store/query"
	"github.com/idpzero/idpzero/pkg/web/models"
	"github.com/idpzero/idpzero/pkg/web/views/pages"
)

func login(config func() *configuration.ServerConfig, queries *query.Queries) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
			return
		}

		conf := config()

		im := models.UserLoginModel{}
		im.AuthRequestID = r.FormValue("req")

		a, err := queries.GetAuthRequestByID(r.Context(), im.AuthRequestID)

		if err != nil {
			if err == sql.ErrNoRows {
				errorView(w, r, err_request_not_found, err)
			} else {
				errorView(w, r, err_unknown_error, err)
			}

			return
		}

		if a.AuthenticatedAt != 0 {
			errorView(w, r, err_already_authenticated, err)
			return
		}

		populateScenarios(&im, conf)

		view := pages.LoginView(im)

		templ.Handler(view).ServeHTTP(w, r)
	})
}

func loginSubmit(config func() *configuration.ServerConfig, queries *query.Queries, callback func(context.Context, string) string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
			return
		}

		username := r.FormValue("username")
		req := r.FormValue("req")

		if username == "" {
			conf := config()

			im := models.UserLoginModel{}
			im.Error = "Select a user to continue"
			im.AuthRequestID = req
			populateScenarios(&im, conf)

			view := pages.LoginView(im)

			templ.Handler(view).ServeHTTP(w, r)
		} else {

			auth := userAuth(config(), username)
			count, err := queries.UpdateAuthRequestUser(r.Context(), query.UpdateAuthRequestUserParams{
				UserID:          username,
				AuthenticatedAt: time.Now().Unix(),
				Acr:             auth.ACR,
				Amr:             strings.Join(auth.AMR, " "),
				ID:              req,
			})

			if err != nil {
				errorView(w, r, err_unknown_error, err)
				return
			}

			if count == 0 {
				errorView(w, r, err_request_not_found, nil)
				return
			}

			http.Redirect(w, r, callback(r.Context(), req), http.StatusFound)
		}
	})
}

func populateScenarios(model *models.UserLoginModel, config *configuration.ServerConfig) {

	model.Users = make([]models.UserOption, 0)

	for _, scenario := range config.Users {
		model.Users = append(model.Users, models.UserOption{
			ID:          scenario.Subject,
			DisplayName: scenario.LoginDisplay,
		})
	}

}

// userAuth is the configured auth of the user with subject, or none.
func userAuth(config *configuration.ServerConfig, subject string) configuration.UserAuth {
	for _, u := range config.Users {
		if u.Subject == subject {
			return u.Auth
		}
	}
	return configuration.UserAuth{}
}
