package configuration

import "testing"

func TestResolvedAuthMethodDefaultsToNone(t *testing.T) {
	cases := map[string]struct {
		configured string
		want       string
	}{
		"empty defaults to none":    {configured: "", want: AuthMethodNone},
		"whitespace defaults":       {configured: "   ", want: AuthMethodNone},
		"explicit none":             {configured: AuthMethodNone, want: AuthMethodNone},
		"basic preserved":           {configured: AuthMethodBasic, want: AuthMethodBasic},
		"post preserved":            {configured: AuthMethodPost, want: AuthMethodPost},
		"private_key_jwt preserved": {configured: AuthMethodPrivateKeyJWT, want: AuthMethodPrivateKeyJWT},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			c := ClientConfig{AuthMethod: tc.configured}
			if got := c.ResolvedAuthMethod(); got != tc.want {
				t.Fatalf("ResolvedAuthMethod() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestIsPublic(t *testing.T) {
	cases := map[string]struct {
		method string
		want   bool
	}{
		"none is public":        {method: AuthMethodNone, want: true},
		"empty is public":       {method: "", want: true},
		"basic is confidential": {method: AuthMethodBasic, want: false},
		"post is confidential":  {method: AuthMethodPost, want: false},
		"jwt is not public":     {method: AuthMethodPrivateKeyJWT, want: false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			c := ClientConfig{AuthMethod: tc.method}
			if got := c.IsPublic(); got != tc.want {
				t.Fatalf("IsPublic() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRequiresClientSecret(t *testing.T) {
	cases := map[string]struct {
		method string
		want   bool
	}{
		"basic requires secret": {method: AuthMethodBasic, want: true},
		"post requires secret":  {method: AuthMethodPost, want: true},
		"none no secret":        {method: AuthMethodNone, want: false},
		"empty no secret":       {method: "", want: false},
		"jwt no secret":         {method: AuthMethodPrivateKeyJWT, want: false},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			c := ClientConfig{AuthMethod: tc.method}
			if got := c.RequiresClientSecret(); got != tc.want {
				t.Fatalf("RequiresClientSecret() = %v, want %v", got, tc.want)
			}
		})
	}
}
