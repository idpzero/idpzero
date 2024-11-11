package query

import (
	"strings"
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

var _ op.AuthRequest = (*AuthRequest)(nil)

func (a *AuthRequest) GetID() string {
	return a.ID
}

func (a *AuthRequest) GetACR() string {
	return "" // we won't handle acr in this example
}

func (a *AuthRequest) GetAMR() []string {
	// this example only uses password for authentication
	if a.Complete {
		return []string{"pwd"}
	}
	return nil
}

func (a *AuthRequest) GetAudience() []string {
	return []string{a.ApplicationID} // this example will always just use the client_id as audience
}

func (a *AuthRequest) GetAuthTime() time.Time {
	return time.Unix(a.AuthenticatedAt, 0)
}

func (a *AuthRequest) GetClientID() string {
	return a.ApplicationID
}

func (a *AuthRequest) GetCodeChallenge() *oidc.CodeChallenge {
	return &oidc.CodeChallenge{
		Challenge: a.CodeChallenge,
		Method:    oidc.CodeChallengeMethod(a.CodeChallengeMethod),
	}
}

func (a *AuthRequest) GetNonce() string {
	return a.Nonce
}

func (a *AuthRequest) GetRedirectURI() string {
	return a.RedirectUri
}

func (a *AuthRequest) GetResponseType() oidc.ResponseType {
	return oidc.ResponseType(a.ResponseType)
}

func (a *AuthRequest) GetResponseMode() oidc.ResponseMode {
	return oidc.ResponseMode(a.ResponseMode)
}

func (a *AuthRequest) GetScopes() []string {
	return strings.Split(a.Scopes, " ")
}

func (a *AuthRequest) GetState() string {
	return a.State
}

func (a *AuthRequest) GetSubject() string {
	return a.UserID
}

func (a *AuthRequest) Done() bool {
	return a.Complete
}
