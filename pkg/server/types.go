package server

import (
	"time"

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

	return c
}

// func (c *Client) parseAndValidate() {
// 	c.validations = []*validation.ChecklistItem{}

// 	ot, err := op.ApplicationTypeString(c.config.ApplicationType)
// 	if err != nil {
// 		ci := validation.NewChecklistItem(false, "application_type").
// 			WithValue(c.config.ApplicationType).
// 			WithError(err)

// 		c.validations = append(c.validations, ci)
// 	} else {
// 		c.appType = ot
// 	}

// 	at, err := op.AccessTokenTypeString(c.config.AccessTokenType)
// 	if err != nil {
// 		ci := validation.NewChecklistItem(false, "access_token_type").
// 			WithValue(c.config.AccessTokenType).
// 			WithError(err)
// 		c.validations = append(c.validations, ci)
// 	} else {
// 		c.accessTokenType = at
// 	}

// 	parsedAuthMethod := oidc.AuthMethod(c.config.AuthMethod)

// 	allAuthMethods := []string{}
// 	for _, am := range oidc.AllAuthMethods {
// 		allAuthMethods = append(allAuthMethods, string(am))
// 		if parsedAuthMethod == am {
// 			c.authMehtod = &am
// 			break
// 		}
// 	}

// 	if c.authMehtod == nil {
// 		ci := validation.NewChecklistItem(false, "auth_method").
// 			WithValue(c.config.AuthMethod).
// 			WithError(fmt.Errorf("'auth_method' not valid")).
// 			WithOptions(allAuthMethods)

// 		c.validations = append(c.validations, ci)
// 	}

// 	for _, gt := range c.config.GrantTypes {
// 		parsedGrantType := oidc.GrantType(gt)
// 		allGrantTypes := []string{}
// 		valid := false
// 		for _, gt := range oidc.AllGrantTypes {
// 			allGrantTypes = append(allGrantTypes, string(gt))
// 			if parsedGrantType == gt {
// 				valid = true
// 				c.grantTypes = append(c.grantTypes, gt)
// 				break
// 			}
// 		}

// 		if !valid {
// 			ci := validation.NewChecklistItem(false, "grant_types").
// 				WithValue(gt).
// 				WithError(fmt.Errorf("'grant_types' not valid")).
// 				WithOptions(allGrantTypes)

// 			c.validations = append(c.validations, ci)
// 		}
// 	}

// 	for _, gt := range c.config.ResponseTypes {
// 		parsedResponseType := oidc.ResponseType(gt)
// 		if parsedResponseType == oidc.ResponseTypeCode ||
// 			parsedResponseType == oidc.ResponseTypeIDToken ||
// 			parsedResponseType == oidc.ResponseTypeIDTokenOnly {
// 			c.responseTypes = append(c.responseTypes, parsedResponseType)
// 		} else {
// 			ci := validation.NewChecklistItem(false, "response_types").
// 				WithValue(gt).
// 				WithError(fmt.Errorf("'response_types' not valid")).
// 				WithOptions([]string{
// 					string(oidc.ResponseTypeCode),
// 					string(oidc.ResponseTypeIDToken),
// 					string(oidc.ResponseTypeIDTokenOnly),
// 				})

// 			c.validations = append(c.validations, ci)
// 		}
// 	}

// }

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
	return c.config.IDTokenUserinfoClaimsAssertion
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
