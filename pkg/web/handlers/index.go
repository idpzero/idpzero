package handlers

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/a-h/templ"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/idpzero/idpzero/pkg/web/models"
	"github.com/idpzero/idpzero/pkg/web/views"
)

func index(config func() *configuration.IDPConfiguration) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		idpConfig := config()

		hosted := fmt.Sprintf("http://localhost:%d", idpConfig.Server.Port)

		discovery, _ := url.JoinPath(hosted, "/.well-known/openid-configuration")
		keys, _ := url.JoinPath(hosted, "/keys")

		im := models.IndexModel{}
		im.Urls = []models.UrlInfo{
			models.UrlInfo{
				Description: "OpenID Connect Discovery Endpoint",
				Url:         discovery,
			},
			models.UrlInfo{
				Description: "JSON Web Key Set (JWKS) Endpoint",
				Url:         keys,
			},
		}

		idx := views.Index(im)

		templ.Handler(views.View((idx))).ServeHTTP(w, r)
	})
}
