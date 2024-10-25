package server

import (
	"fmt"
	"net/http"

	"github.com/zitadel/oidc/v3/pkg/op"
)

func setProviderFromRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		scheme := "http"
		if r.TLS != nil {
			scheme = "https"
		}

		ctx := op.ContextWithIssuer(r.Context(), fmt.Sprintf("%s://%s", scheme, r.Host))

		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
