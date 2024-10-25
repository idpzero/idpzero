package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/idpzero/idpzero/pkg/web/views"
)

func index() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		idx := views.Index()

		templ.Handler(views.View((idx))).ServeHTTP(w, r)
	})
}
