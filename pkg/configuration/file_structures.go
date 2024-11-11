package configuration

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func init() {
	// set the tag name so the right field is shown in the error message
	validation.ErrorTag = "yaml"
}

// ServerConfig is a struct that holds the server configuration and is generally stored in source control for shared use.
type ServerConfig struct {
	Server         HostConfig               `yaml:"server"`
	Clients        []ClientConfig           `yaml:"clients"`
	ScenarioGroups map[string]ScenarioGroup `yaml:"scenario_groups"`
}

func (h ServerConfig) Validate() error {
	return validation.ValidateStruct(&h,
		validation.Field(&h.Server),
		validation.Field(&h.Clients),
	)
}

type HostConfig struct {
	Port      int    `yaml:"port"`
	KeyPhrase string `yaml:"keyphrase"`
}

func (h HostConfig) Validate() error {
	return validation.ValidateStruct(&h,
		validation.Field(&h.Port, validation.Required, validation.Min(1), validation.Max(65535)),
		validation.Field(&h.KeyPhrase, validation.Required),
	)
}

type ClientConfig struct {
	ClientID                       string        `yaml:"client_id"`
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
	ClientSecret                   string        `yaml:"client_secret,omitempty"`
}

type ScenarioGroup struct {
	Order     int        `yaml:"order"`
	Display   string     `yaml:"display"`
	Scenarios []Scenario `yaml:"scenarios"`
}

type Scenario struct {
	ID      string                 `yaml:"id"`
	Display string                 `yaml:"display"`
	Comment string                 `yaml:"comment"`
	Claims  map[string]interface{} `yaml:"claims"`
}

// func (h ClientConfig) Validate() error {

// 	return validation.ValidateStruct(&h,
// 		validation.Field(&h.ID, validation.Required),
// 		validation.Field(&h.AccessTokenType, validation.Required, validation.In(all(op.AccessTokenTypeStrings())...)),
// 		validation.Field(&h.ApplicationType, validation.Required, validation.In(all(op.ApplicationTypeStrings())...)),
// 		validation.Field(&h.AuthMethod, validation.Required, validation.In(all(oidc.AllAuthMethods)...)),
// 		validation.Field(&h.GrantTypes, validation.Required, validation.Each(validation.In(all(oidc.AllGrantTypes)...))),
// 		validation.Field(&h.ResponseTypes, validation.Required, validation.Each(validation.In(string(oidc.ResponseTypeCode), string(oidc.ResponseTypeIDToken), string(oidc.ResponseTypeIDTokenOnly)))),
// 	)
// }

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
