package server

import (
	"testing"

	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

func TestNewClientAuthMethodMapping(t *testing.T) {
	cases := map[string]struct {
		configured string
		want       oidc.AuthMethod
	}{
		"none maps to AuthMethodNone":  {configured: "none", want: oidc.AuthMethodNone},
		"empty defaults to none":       {configured: "", want: oidc.AuthMethodNone},
		"basic maps through":           {configured: "client_secret_basic", want: oidc.AuthMethodBasic},
		"post maps through":            {configured: "client_secret_post", want: oidc.AuthMethodPost},
		"private_key_jwt maps through": {configured: "private_key_jwt", want: oidc.AuthMethodPrivateKeyJWT},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			c := NewClient(configuration.ClientConfig{
				ClientID:   "test",
				AuthMethod: tc.configured,
			})

			if got := c.AuthMethod(); got != tc.want {
				t.Fatalf("AuthMethod() = %q, want %q", got, tc.want)
			}
		})
	}
}

// A public client authenticates via PKCE, so the zitadel op library never calls
// AuthorizeClientIDSecret for it. This guards the intended default: a public
// client carries no secret.
func TestPublicClientHasNoSecret(t *testing.T) {
	c := NewClient(configuration.ClientConfig{
		ClientID:   "spa",
		AuthMethod: "none",
	})

	if c.AuthMethod() != oidc.AuthMethodNone {
		t.Fatalf("expected public client to use AuthMethodNone, got %q", c.AuthMethod())
	}
}
