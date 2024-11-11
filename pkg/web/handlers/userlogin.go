package handlers

import (
	"fmt"
	"net/http"

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
			im := models.UserLoginModel{}
			im.Error = "Select user to continue"
			im.AuthRequestID = r.FormValue("req")

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

		im := models.UserLoginModel{}
		im.AuthRequestID = r.FormValue("req")

		idx := views.UserLogin(im)

		templ.Handler(views.PanelView((idx))).ServeHTTP(w, r)
	})
}
