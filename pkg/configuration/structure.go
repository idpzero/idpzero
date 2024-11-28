package configuration

import (
	"time"
)

// ServerConfig is a struct that holds the server configuration and is generally stored in source control for shared use.
type ServerConfig struct {
	Server  HostConfig      `yaml:"server"`
	Clients []*ClientConfig `yaml:"clients"`
	Users   []*User         `yaml:"users"`
}

type HostConfig struct {
	Port      int    `yaml:"port"`
	KeyPhrase string `yaml:"keyphrase"`
}

type ClientConfig struct {
	Name                           string              `yaml:"name"`
	ClientID                       string              `yaml:"client_id"`
	AccessTokenType                string              `yaml:"access_token_type"`           // bearer or jwt.
	ApplicationType                string              `yaml:"application_type"`            // web, native, or service.
	AuthMethod                     string              `yaml:"auth_method"`                 // client_secret_basic,client_secret_post,none,private_key_jwt
	ClockSkew                      time.Duration       `yaml:"clock_skew,omitempty"`        // time in duration format
	IDTokenLifetime                time.Duration       `yaml:"id_token_lifetime,omitempty"` // time in duration format
	IDTokenUserinfoClaimsAssertion bool                `yaml:"id_token_userinfo_claims_assertion,omitempty"`
	GrantTypes                     []string            `yaml:"grant_types"` // authorization_code,implicit,password,client_credentials,refresh_token etc
	RedirectURIs                   []string            `yaml:"redirect_uris"`
	PostLogoutRedirectURIs         []string            `yaml:"post_logout_redirect_uris,omitempty"`
	ResponseTypes                  []string            `yaml:"response_types"`
	ClientSecret                   string              `yaml:"-"` // ignore when marshalling
	CustomScopes                   map[string][]string `yaml:"custom_scopes,omitempty"`
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
