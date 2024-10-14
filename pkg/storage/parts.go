package storage

import (
	"time"

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

var _ op.AuthRequest = authReq{}

type authReq struct {
}

// Done implements op.AuthRequest.
func (a authReq) Done() bool {
	panic("unimplemented")
}

// GetACR implements op.AuthRequest.
func (a authReq) GetACR() string {
	panic("unimplemented")
}

// GetAMR implements op.AuthRequest.
func (a authReq) GetAMR() []string {
	panic("unimplemented")
}

// GetAudience implements op.AuthRequest.
func (a authReq) GetAudience() []string {
	panic("unimplemented")
}

// GetAuthTime implements op.AuthRequest.
func (a authReq) GetAuthTime() time.Time {
	panic("unimplemented")
}

// GetClientID implements op.AuthRequest.
func (a authReq) GetClientID() string {
	panic("unimplemented")
}

// GetCodeChallenge implements op.AuthRequest.
func (a authReq) GetCodeChallenge() *oidc.CodeChallenge {
	panic("unimplemented")
}

// GetID implements op.AuthRequest.
func (a authReq) GetID() string {
	panic("unimplemented")
}

// GetNonce implements op.AuthRequest.
func (a authReq) GetNonce() string {
	panic("unimplemented")
}

// GetRedirectURI implements op.AuthRequest.
func (a authReq) GetRedirectURI() string {
	panic("unimplemented")
}

// GetResponseMode implements op.AuthRequest.
func (a authReq) GetResponseMode() oidc.ResponseMode {
	panic("unimplemented")
}

// GetResponseType implements op.AuthRequest.
func (a authReq) GetResponseType() oidc.ResponseType {
	panic("unimplemented")
}

// GetScopes implements op.AuthRequest.
func (a authReq) GetScopes() []string {
	panic("unimplemented")
}

// GetState implements op.AuthRequest.
func (a authReq) GetState() string {
	panic("unimplemented")
}

// GetSubject implements op.AuthRequest.
func (a authReq) GetSubject() string {
	panic("unimplemented")
}

// type AuthRequest struct {
// 	ID            string
// 	CreationDate  time.Time
// 	ApplicationID string
// 	CallbackURI   string
// 	TransferState string
// 	Prompt        []string
// 	UiLocales     []language.Tag
// 	LoginHint     string
// 	MaxAuthAge    *time.Duration
// 	UserID        string
// 	Scopes        []string
// 	ResponseType  oidc.ResponseType
// 	ResponseMode  oidc.ResponseMode
// 	Nonce         string
// 	CodeChallenge *OIDCCodeChallenge

// 	done     bool
// 	authTime time.Time
// }

// // LogValue allows you to define which fields will be logged.
// // Implements the [slog.LogValuer]
// func (a *AuthRequest) LogValue() slog.Value {
// 	return slog.GroupValue(
// 		slog.String("id", a.ID),
// 		slog.Time("creation_date", a.CreationDate),
// 		slog.Any("scopes", a.Scopes),
// 		slog.String("response_type", string(a.ResponseType)),
// 		slog.String("app_id", a.ApplicationID),
// 		slog.String("callback_uri", a.CallbackURI),
// 	)
// }

// func (a *AuthRequest) GetID() string {
// 	return a.ID
// }

// func (a *AuthRequest) GetACR() string {
// 	return "" // we won't handle acr in this example
// }

// func (a *AuthRequest) GetAMR() []string {
// 	// this example only uses password for authentication
// 	if a.done {
// 		return []string{"pwd"}
// 	}
// 	return nil
// }

// func (a *AuthRequest) GetAudience() []string {
// 	return []string{a.ApplicationID} // this example will always just use the client_id as audience
// }

// func (a *AuthRequest) GetAuthTime() time.Time {
// 	return a.authTime
// }

// func (a *AuthRequest) GetClientID() string {
// 	return a.ApplicationID
// }

// func (a *AuthRequest) GetCodeChallenge() *oidc.CodeChallenge {
// 	return CodeChallengeToOIDC(a.CodeChallenge)
// }

// func (a *AuthRequest) GetNonce() string {
// 	return a.Nonce
// }

// func (a *AuthRequest) GetRedirectURI() string {
// 	return a.CallbackURI
// }

// func (a *AuthRequest) GetResponseType() oidc.ResponseType {
// 	return a.ResponseType
// }

// func (a *AuthRequest) GetResponseMode() oidc.ResponseMode {
// 	return a.ResponseMode
// }

// func (a *AuthRequest) GetScopes() []string {
// 	return a.Scopes
// }

// func (a *AuthRequest) GetState() string {
// 	return a.TransferState
// }

// func (a *AuthRequest) GetSubject() string {
// 	return a.UserID
// }

// func (a *AuthRequest) Done() bool {
// 	return a.done
// }

// type OIDCCodeChallenge struct {
// 	Challenge string
// 	Method    string
// }

// func CodeChallengeToOIDC(challenge *OIDCCodeChallenge) *oidc.CodeChallenge {
// 	if challenge == nil {
// 		return nil
// 	}
// 	challengeMethod := oidc.CodeChallengeMethodPlain
// 	if challenge.Method == "S256" {
// 		challengeMethod = oidc.CodeChallengeMethodS256
// 	}
// 	return &oidc.CodeChallenge{
// 		Challenge: challenge.Challenge,
// 		Method:    challengeMethod,
// 	}
// }
