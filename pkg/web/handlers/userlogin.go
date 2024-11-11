package handlers

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/a-h/templ"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/idpzero/idpzero/pkg/web/models"
	"github.com/idpzero/idpzero/pkg/web/views"
)

func userloginSubmit(config func() *configuration.ServerConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
			return
		}

		username := r.FormValue("username")

		if username == "" {

			conf := config()

			im := models.UserLoginModel{}
			im.Error = "Select a user to continue"
			im.AuthRequestID = r.FormValue("req")
			populateScenarios(&im, conf)

			idx := views.UserLogin(im)

			templ.Handler(views.PanelView((idx))).ServeHTTP(w, r)
		} else {
			// todo - call back into the OIDC flow.
			http.Redirect(w, r, "/", http.StatusFound)
		}
	})
}

func userlogin(config func() *configuration.ServerConfig) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseForm()
		if err != nil {
			http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
			return
		}

		conf := config()

		im := models.UserLoginModel{}
		im.AuthRequestID = r.FormValue("req")
		populateScenarios(&im, conf)

		idx := views.UserLogin(im)

		templ.Handler(views.PanelView((idx))).ServeHTTP(w, r)
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
		fmt.Println(group.Display)
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
