package idp

import (
	"fmt"
	"strings"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

var (
	_ op.Key        = &opPublicKey{}  // make sure my type implements the interface
	_ op.SigningKey = &opPrivateKey{} // make sure my type implements the interface
)

type opPublicKey struct {
	key configuration.Key
}

func (s *opPublicKey) ID() string {
	return s.key.ID
}

func (s *opPublicKey) Algorithm() jose.SignatureAlgorithm {
	return jose.SignatureAlgorithm(s.key.Algorithm)
}

func (s *opPublicKey) Use() string {
	return s.key.Use
}

func (s *opPublicKey) Key() any {

	parsed, err := parseRSAPublicKey(s.key)
	if err != nil {
		return err
	}
	return parsed

}

type opPrivateKey struct {
	key configuration.Key
}

func (s *opPrivateKey) SignatureAlgorithm() jose.SignatureAlgorithm {
	return jose.SignatureAlgorithm(s.key.Algorithm)
}

func (s *opPrivateKey) Key() any {

	priv, _, err := parseRSAKey(s.key)

	if err != nil {
		return err
	}

	return priv
}

func (s *opPrivateKey) ID() string {
	return s.key.ID
}

type Client struct {
	config           configuration.ClientConfig
	validationErrors []error

	// parsed fields
	appType         op.ApplicationType
	accessTokenType op.AccessTokenType
	authMehtod      *oidc.AuthMethod
	grantTypes      []oidc.GrantType
	clockSkew       time.Duration
	responseTypes   []oidc.ResponseType
	idTokenLifetime time.Duration
}

var _ op.Client = &Client{}

func NewClient(config configuration.ClientConfig) (*Client, []error) {
	c := &Client{
		config:           config,
		validationErrors: []error{},
		grantTypes:       []oidc.GrantType{},
		responseTypes:    []oidc.ResponseType{},
	}

	// parse and validate the config, so we know if its valid and can
	// prevent starting the server if its not.
	c.parseAndValidate()

	return c, c.validationErrors
}

func (c *Client) parseAndValidate() {
	c.validationErrors = []error{}

	ot, err := op.ApplicationTypeString(c.config.ApplicationType)
	if err != nil {
		c.validationErrors = append(c.validationErrors, err)
	} else {
		c.appType = ot
	}

	at, err := op.AccessTokenTypeString(c.config.AccessTokenType)
	if err != nil {
		c.validationErrors = append(c.validationErrors, err)
	} else {
		c.accessTokenType = at
	}

	parsedAuthMethod := oidc.AuthMethod(c.config.AuthMethod)

	allAuthMethods := []string{}
	for _, am := range oidc.AllAuthMethods {
		allAuthMethods = append(allAuthMethods, string(am))
		if parsedAuthMethod == am {
			c.authMehtod = &am
			break
		}
	}

	if c.authMehtod == nil {
		c.validationErrors = append(c.validationErrors, fmt.Errorf("'auth_method' not valid, expecting one of %s", strings.Join(allAuthMethods, ",")))
	}

	for _, gt := range c.config.GrantTypes {
		parsedGrantType := oidc.GrantType(gt)
		allGrantTypes := []string{}
		valid := false
		for _, gt := range oidc.AllGrantTypes {
			allGrantTypes = append(allGrantTypes, string(gt))
			if parsedGrantType == gt {
				valid = true
				c.grantTypes = append(c.grantTypes, gt)
				break
			}
		}

		if !valid {
			c.validationErrors = append(c.validationErrors, fmt.Errorf("'grant_types' not valid, expecting one of more of %s", strings.Join(allGrantTypes, ",")))
		}
	}

	for _, gt := range c.config.ResponseTypes {
		parsedResponseType := oidc.ResponseType(gt)
		if parsedResponseType == oidc.ResponseTypeCode ||
			parsedResponseType == oidc.ResponseTypeIDToken ||
			parsedResponseType == oidc.ResponseTypeIDTokenOnly {
			c.responseTypes = append(c.responseTypes, parsedResponseType)
		} else {
			c.validationErrors = append(c.validationErrors, fmt.Errorf("'response_types' not valid, expecting one of more of %s, %s, %s", oidc.ResponseTypeCode, oidc.ResponseTypeIDToken, oidc.ResponseTypeIDTokenOnly))
		}
	}

	// parse clock skew or use the default
	defaultSkew, _ := time.ParseDuration("60s")
	clockSkew, err := parseDurationOrDefault(c.config.ClockSkew, defaultSkew)
	if err != nil {
		c.validationErrors = append(c.validationErrors, fmt.Errorf("'clock_skew' not valid, expecting a duration or ommit to use default (%s)", defaultSkew))
	}
	c.clockSkew = clockSkew

	// parse id token lifetime or use the default
	defaultIDTokenLifetime, _ := time.ParseDuration("1h")
	idTokenLifetime, err := parseDurationOrDefault(c.config.IDTokenLifetime, defaultIDTokenLifetime)
	if err != nil {
		c.validationErrors = append(c.validationErrors, fmt.Errorf("'id_token_lifetime' not valid, expecting a duration or ommit to use default (%s)", defaultIDTokenLifetime))
	}
	c.idTokenLifetime = idTokenLifetime

}

func (c *Client) IsValid() bool {
	return len(c.validationErrors) == 0
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
	return c.clockSkew
}

// DevMode implements op.Client.
func (c *Client) DevMode() bool {
	return true
}

// GetID implements op.Client.
func (c *Client) GetID() string {
	return c.config.ID
}

// GrantTypes implements op.Client.
func (c *Client) GrantTypes() []oidc.GrantType {
	return c.grantTypes
}

// IDTokenLifetime implements op.Client.
func (c *Client) IDTokenLifetime() time.Duration {
	return c.idTokenLifetime
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
	return "/login/user?authRequestID=" + id
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

func parseDurationOrDefault(s *string, def time.Duration) (time.Duration, error) {
	if s == nil {
		return def, nil
	}

	t, err := time.ParseDuration(*s)
	if err != nil {
		return def, err
	}

	return t, nil
}
