package controllers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/idpzero/idpzero/pkg/web/views/pages"
)

// loggedOut is where end_session redirects when the request names no
// post_logout_redirect_uri.
func loggedOut() http.Handler {
	return templ.Handler(pages.LoggedOutView())
}
