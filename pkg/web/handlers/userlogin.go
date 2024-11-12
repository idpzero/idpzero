package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/a-h/templ"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/idpzero/idpzero/pkg/store/query"
	"github.com/idpzero/idpzero/pkg/web/models"
	"github.com/idpzero/idpzero/pkg/web/views/pages"
)

func userlogin(config func() *configuration.ServerConfig, queries *query.Queries) http.HandlerFunc {
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

func userloginSubmit(config func() *configuration.ServerConfig, queries *query.Queries, callback func(context.Context, string) string) http.HandlerFunc {
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

			count, err := queries.UpdateAuthRequestUser(r.Context(), query.UpdateAuthRequestUserParams{
				UserID:          username,
				AuthenticatedAt: time.Now().Unix(),
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

	model.Users = make([]models.OptionGroup, 0)

	groups := []configuration.ScenarioGroup{}

	for _, v := range config.ScenarioGroups {
		groups = append(groups, v)
	}

	// Sort the groups by the field order
	sort.Slice(groups, func(i, j int) bool {
		return groups[i].Order < groups[j].Order
	})

	for _, group := range groups {
		og := models.OptionGroup{
			DisplayName: group.Display,
			Options:     make([]models.Option, 0),
		}

		for _, scenario := range group.Scenarios {
			og.Options = append(og.Options, models.Option{
				ID:          scenario.ID,
				DisplayName: scenario.Display,
			})
		}

		model.Users = append(model.Users, og)
	}

}
