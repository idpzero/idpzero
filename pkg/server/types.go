package server

import (
	"time"

	"github.com/fatih/color"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type Client struct {
	config configuration.ClientConfig

	// parsed fields
	appType         op.ApplicationType
	accessTokenType op.AccessTokenType
	authMehtod      *oidc.AuthMethod
	grantTypes      []oidc.GrantType
	responseTypes   []oidc.ResponseType
}

var _ op.Client = &Client{}

func NewClient(config configuration.ClientConfig) *Client {
	c := &Client{
		config:        config,
		grantTypes:    []oidc.GrantType{},
		responseTypes: []oidc.ResponseType{},
	}

	for _, gt := range config.GrantTypes {
		c.grantTypes = append(c.grantTypes, oidc.GrantType(gt))
	}

	for _, gt := range config.ResponseTypes {
		c.responseTypes = append(c.responseTypes, oidc.ResponseType(gt))
	}

	c.authMehtod = (*oidc.AuthMethod)(&config.AuthMethod)

	att, err := op.AccessTokenTypeString(config.AccessTokenType)

	if err != nil {
		color.Red("Error parsing client AccessTokenType: %s - defaulting to jwt", config.AccessTokenType)
		c.accessTokenType = op.AccessTokenTypeJWT
	} else {
		c.accessTokenType = att
	}

	at, err := op.ApplicationTypeString(config.ApplicationType)

	if err != nil {
		color.Red("Error parsing client ApplicationType: %s - defaulting to web", config.ApplicationType)
		c.appType = op.ApplicationTypeWeb
	} else {
		c.appType = at
	}

	return c
}

// AccessTokenType implements op.Client.
func (c *Client) AccessTokenType() op.AccessTokenType {
	return c.accessTokenType
}

// ApplicationType implements op.Client.
func (c *Client) ApplicationType() op.ApplicationType {
	return c.appType
}

// AuthMethod implements op.Client.
func (c *Client) AuthMethod() oidc.AuthMethod {
	// will always be validated first before this is called.
	return *c.authMehtod
}

// ClockSkew implements op.Client.
func (c *Client) ClockSkew() time.Duration {
	return c.config.ClockSkew
}

// DevMode implements op.Client.
func (c *Client) DevMode() bool {
	return true
}

// GetID implements op.Client.
func (c *Client) GetID() string {
	return c.config.ClientID
}

// GrantTypes implements op.Client.
func (c *Client) GrantTypes() []oidc.GrantType {
	return c.grantTypes
}

// IDTokenLifetime implements op.Client.
func (c *Client) IDTokenLifetime() time.Duration {
	return c.config.IDTokenLifetime
}

// IDTokenUserinfoClaimsAssertion implements op.Client.
func (c *Client) IDTokenUserinfoClaimsAssertion() bool {

	// we need to invert the value, so the default (false) will work
	// as expected and NOT omit user info claims
	//
	//  If set to true, the zitadel framework will omit these scopes from the ID token:
	//          oidc.ScopeProfile,
	// 			oidc.ScopeEmail,
	// 			oidc.ScopeAddress,
	// 			oidc.ScopePhone:
	return !c.config.IDTokenOmitUserInfoClaims
}

// IsScopeAllowed implements op.Client.
func (c *Client) IsScopeAllowed(scope string) bool {
	return true
}

// LoginURL implements op.Client.
func (c *Client) LoginURL(id string) string {
	return "/login?req=" + id
}

// PostLogoutRedirectURIs implements op.Client.
func (c *Client) PostLogoutRedirectURIs() []string {
	return c.config.PostLogoutRedirectURIs
}

// RedirectURIs implements op.Client.
func (c *Client) RedirectURIs() []string {
	return c.config.RedirectURIs
}

// ResponseTypes implements op.Client.
func (c *Client) ResponseTypes() []oidc.ResponseType {
	return c.responseTypes
}

// RestrictAdditionalAccessTokenScopes implements op.Client.
func (c *Client) RestrictAdditionalAccessTokenScopes() func(scopes []string) []string {
	return func(scopes []string) []string {
		return scopes
	}
}

// RestrictAdditionalIdTokenScopes implements op.Client.
func (c *Client) RestrictAdditionalIdTokenScopes() func(scopes []string) []string {
	return func(scopes []string) []string {
		return scopes
	}
}
