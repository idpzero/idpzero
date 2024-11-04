package configuration

import "time"

// ServerConfig is a struct that holds the server configuration and is generally stored in source control for shared use.
type ServerConfig struct {
	Server  HostConfig     `yaml:"server"`
	Clients []ClientConfig `yaml:"clients"`
}

type HostConfig struct {
	Port      int    `yaml:"port"`
	KeyPhrase string `yaml:"keyphrase"`
}

type ClientConfig struct {
	ID                             string        `yaml:"id"`
	AccessTokenType                string        `yaml:"access_token_type"` // bearer or jwt.
	ApplicationType                string        `yaml:"application_type"`  // web, native, or service.
	AuthMethod                     string        `yaml:"auth_method"`       // client_secret_basic,client_secret_post,none,private_key_jwt
	ClockSkew                      time.Duration `yaml:"clock_skew"`        // time in duration format
	IDTokenLifetime                time.Duration `yaml:"id_token_lifetime"` // time in duration format
	IDTokenUserinfoClaimsAssertion bool          `yaml:"id_token_userinfo_claims_assertion"`
	GrantTypes                     []string      `yaml:"grant_types"` // authorization_code,implicit,password,client_credentials,refresh_token etc
	RedirectURIs                   []string      `yaml:"redirect_uris"`
	PostLogoutRedirectURIs         []string      `yaml:"post_logout_redirect_uris"`
	ResponseTypes                  []string      `yaml:"response_types"`
}

// KeysConfiguration is a struct that holds the keys configuration and is stored against the local user account so that it
// can be used to sign and verify tokens, and survive restarts (and not be committed to source control).
type KeysConfiguration struct {
	Keys []Key `yaml:"keys"`
}

type Key struct {
	ID        string            `yaml:"id"`
	Algorithm string            `yaml:"algorithm"`
	Use       string            `yaml:"use"`
	Data      map[string]string `yaml:"data"`
}
