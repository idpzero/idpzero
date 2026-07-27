package configuration

import (
	"strings"
	"time"
)

// Supported client authentication methods (mirrors the oidc.AuthMethod* values).
const (
	AuthMethodBasic         = "client_secret_basic"
	AuthMethodPost          = "client_secret_post"
	AuthMethodNone          = "none"
	AuthMethodPrivateKeyJWT = "private_key_jwt"
)

// ServerConfig is a struct that holds the server configuration and is generally stored in source control for shared use.
type ServerConfig struct {
	Server  HostConfig      `yaml:"server"`
	SCIM    *SCIMConfig     `yaml:"scim,omitempty"`
	Clients []*ClientConfig `yaml:"clients"`
	Users   []*User         `yaml:"users"`
}

type HostConfig struct {
	Port      int    `yaml:"port"`
	KeyPhrase string `yaml:"keyphrase"`
}

// SCIMConfig configures outbound SCIM 2.0 provisioning, where idpzero acts as a
// SCIM client and pushes its configured users to an external SCIM service.
type SCIMConfig struct {
	Endpoint    string `yaml:"endpoint"`               // base URL of the target SCIM service, e.g. https://host/scim/v2
	BearerToken string `yaml:"bearer_token,omitempty"` // optional; the IDPZERO_SCIM_TOKEN env var takes precedence
}

type ClientConfig struct {
	Name                      string              `yaml:"name"`
	ClientID                  string              `yaml:"client_id"`
	AccessTokenType           string              `yaml:"access_token_type"`                       // bearer or jwt.
	ApplicationType           string              `yaml:"application_type"`                        // web, native, or service.
	AuthMethod                string              `yaml:"auth_method"`                             // client_secret_basic,client_secret_post,none,private_key_jwt
	ClockSkew                 time.Duration       `yaml:"clock_skew,omitempty"`                    // time in duration format
	IDTokenLifetime           time.Duration       `yaml:"id_token_lifetime,omitempty"`             // time in duration format
	IDTokenOmitUserInfoClaims bool                `yaml:"id_token_omit_userinfo_claims,omitempty"` // true to omit profile, email, phone, address scopes from ID tokens
	GrantTypes                []string            `yaml:"grant_types"`                             // authorization_code,implicit,password,client_credentials,refresh_token etc
	RedirectURIs              []string            `yaml:"redirect_uris"`
	PostLogoutRedirectURIs    []string            `yaml:"post_logout_redirect_uris,omitempty"`
	ResponseTypes             []string            `yaml:"response_types"`
	ClientSecret              string              `yaml:"-"` // ignore when marshalling
	CustomScopes              map[string][]string `yaml:"custom_scopes,omitempty"`
}

// ResolvedAuthMethod returns the effective authentication method for the client.
// An empty value defaults to "none", which enables a public client using PKCE
// and therefore requires no client secret. This makes it safe to commit client
// configuration to source control without storing any secrets.
func (c *ClientConfig) ResolvedAuthMethod() string {
	if strings.TrimSpace(c.AuthMethod) == "" {
		return AuthMethodNone
	}
	return c.AuthMethod
}

// IsPublic reports whether the client is a public client (auth_method: none),
// which authenticates the authorization code exchange using PKCE rather than a
// client secret.
func (c *ClientConfig) IsPublic() bool {
	return c.ResolvedAuthMethod() == AuthMethodNone
}

// RequiresClientSecret reports whether the configured authentication method
// relies on a shared client secret. Only client_secret_basic and
// client_secret_post use a secret; none (PKCE) and private_key_jwt do not.
func (c *ClientConfig) RequiresClientSecret() bool {
	switch c.ResolvedAuthMethod() {
	case AuthMethodBasic, AuthMethodPost:
		return true
	default:
		return false
	}
}

type User struct {
	Subject      string     `yaml:"subject"`
	LoginDisplay string     `yaml:"login_display"`
	Claims       UserClaims `yaml:"claims"`
}

type UserClaims struct {
	Email             *string                `yaml:"email,omitempty"`
	EmailVerified     *bool                  `yaml:"email_verified,omitempty"`
	Phone             *string                `yaml:"phone,omitempty"`
	PhoneVerified     *bool                  `yaml:"phone_verified,omitempty"`
	Name              *string                `yaml:"name,omitempty"`
	PreferredUsername *string                `yaml:"preferred_username,omitempty"`
	Nickname          *string                `yaml:"nickname,omitempty"`
	GivenName         *string                `yaml:"given_name,omitempty"`
	MiddleName        *string                `yaml:"middle_name,omitempty"`
	FamilyName        *string                `yaml:"family_name,omitempty"`
	UpdatedAt         *time.Time             `yaml:"updated_at,omitempty"`
	Custom            map[string]interface{} `yaml:"custom,omitempty"`
}
