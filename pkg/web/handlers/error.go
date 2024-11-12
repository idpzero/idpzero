package handlers

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/idpzero/idpzero/pkg/dbg"
	"github.com/idpzero/idpzero/pkg/web/models"
	"github.com/idpzero/idpzero/pkg/web/views/pages"
)

const (
	err_already_authenticated = "err_already_authenticated"
	err_request_not_found     = "err_request_not_found"
	err_unknown_error         = "err_unknown_error"
)

func errorView(w http.ResponseWriter, r *http.Request, code string, err error) {

	message := "Request is invalid"

	switch code {
	case err_already_authenticated:
		{
			message = "The request identifier has already been authenticated."
		}
	case err_request_not_found:
		{
			message = "Request could not be updated as the identifier was not found."
		}
	case err_unknown_error:
		{
			message = "Unknown error, check logs for more information"
		}
	}

	args := []interface{}{slog.String("code", code), slog.String("message", message)}

	if err != nil {
		args = append(args, slog.String("error", err.Error()))
	}

	dbg.Logger.Error("Error occured", args...)

	ev := pages.ErrorView(models.ErrorModel{
		Code:    code,
		Message: message,
	})

	templ.Handler(ev).ServeHTTP(w, r)
}
