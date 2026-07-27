package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/idpzero/idpzero/pkg/scim"
	"github.com/idpzero/idpzero/pkg/web/models"
	"github.com/idpzero/idpzero/pkg/web/views/pages"
)

// scimTokenEnv is an optional environment variable holding the outbound SCIM
// bearer token. It takes precedence over the value in the configuration file so
// real secrets need not be committed to source control.
const scimTokenEnv = "IDPZERO_SCIM_TOKEN"

func scimStatusView(config func() *configuration.ServerConfig, statusGet func() *scim.Report) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conf := config()

		im := models.NewSCIMModel(conf.SCIM, scimTokenConfigured(conf), len(conf.Users), statusGet())

		view := pages.SCIMView(im)

		templ.Handler(view).ServeHTTP(w, r)
	})
}

func scimSync(config func() *configuration.ServerConfig, statusSet func(*scim.Report)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conf := config()

		if conf.SCIM == nil || conf.SCIM.Endpoint == "" {
			// Nothing to sync against; return to the page which renders the warning.
			http.Redirect(w, r, "/scim", http.StatusSeeOther)
			return
		}

		client := scim.New(conf.SCIM.Endpoint, resolveSCIMToken(conf))
		report := scim.Sync(r.Context(), client, conf.SCIM.Endpoint, conf.Users)
		report.RanAt = time.Now()
		statusSet(&report)

		// Post/Redirect/Get so a refresh doesn't re-fire the sync.
		http.Redirect(w, r, "/scim", http.StatusSeeOther)
	})
}

func scimReset(statusSet func(*scim.Report)) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Clear the last-run status so the page returns to its "never run" state.
		// This only purges idpzero's local sync history; it does not remove any
		// users already provisioned on the target.
		statusSet(nil)

		http.Redirect(w, r, "/scim", http.StatusSeeOther)
	})
}

// resolveSCIMToken returns the outbound bearer token, preferring the environment
// variable over the configured value.
func resolveSCIMToken(conf *configuration.ServerConfig) string {
	if v := os.Getenv(scimTokenEnv); v != "" {
		return v
	}
	if conf.SCIM != nil {
		return conf.SCIM.BearerToken
	}
	return ""
}

func scimTokenConfigured(conf *configuration.ServerConfig) bool {
	return resolveSCIMToken(conf) != ""
}
