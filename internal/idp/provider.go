package idp

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/zitadel/oidc/v3/pkg/op"
	"golang.org/x/text/language"
)

const (
	pathLoggedOut = "/logged-out"
)

type ProviderOptions struct {
	Key          [32]byte
	Storage      op.Storage
	ExtraOptions []op.Option
}

func NewProvider(logger *slog.Logger, options ProviderOptions) (*op.Provider, error) {
	config := &op.Config{
		CryptoKey: options.Key,

		// will be used if the end_session endpoint is called without a post_logout_redirect_uri
		DefaultLogoutRedirectURI: pathLoggedOut,

		// enables code_challenge_method S256 for PKCE (and therefore PKCE in general)
		CodeMethodS256: true,

		// enables additional client_id/client_secret authentication by form post (not only HTTP Basic Auth)
		AuthMethodPost: true,

		// enables additional authentication by using private_key_jwt
		AuthMethodPrivateKeyJWT: true,

		// enables refresh_token grant use
		GrantTypeRefreshToken: true,

		// enables use of the `request` Object parameter
		RequestObjectSupported: true,

		// this example has only static texts (in English), so we'll set the here accordingly
		SupportedUILocales: []language.Tag{language.English},

		DeviceAuthorization: op.DeviceAuthorizationConfig{
			Lifetime:     5 * time.Minute,
			PollInterval: 5 * time.Second,
			UserFormPath: "/device",
			UserCode:     op.UserCodeBase20,
		},
	}

	// options for the provider.
	po := append([]op.Option{
		//we must explicitly allow the use of the http issuer
		op.WithAllowInsecure(),
		// as an example on how to customize an endpoint this will change the authorization_endpoint from /authorize to /auth
		//op.WithCustomAuthEndpoint(op.NewEndpoint("auth")),
		// Pass our logger to the OP
		op.WithLogger(logger.WithGroup("op")),
	}, options.ExtraOptions...)

	// use the value from the context
	issuerFunc := func(insecure bool) (op.IssuerFromRequest, error) {
		var x op.IssuerFromRequest = func(r *http.Request) string {
			return op.IssuerFromContext(r.Context())
		}
		return x, nil
	}

	handler, err := op.NewProvider(config, options.Storage, issuerFunc, po...)

	if err != nil {
		return nil, err
	}
	return handler, nil
}
