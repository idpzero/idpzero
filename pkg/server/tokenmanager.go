package server

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/idpzero/idpzero/pkg/store/query"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type tokenManager struct {
	queries *query.Queries
	clients *clients
}

func newTokenManager(queries *query.Queries, clients *clients) *tokenManager {
	return &tokenManager{
		queries: queries,
		clients: clients,
	}
}

type CreateTokenArgs struct {
	ApplicationID  string
	RequestID      *string
	Scopes         []string
	Audience       []string
	Subject        string
	RefreshTokenID *string
}

func (t *tokenManager) CreateToken(ctx context.Context, tx *sql.Tx, args CreateTokenArgs) (*query.Token, error) {
	qtx := t.queries.WithTx(tx)

	// clientConfig, ok := t.clients.GetByID(args.ApplicationID)

	// if !ok {
	// 	return nil, errors.New("client not found")
	// }

	var tokenLifetime = time.Hour

	var refToken = sql.NullString{}
	var authReqId = sql.NullString{}

	if args.RequestID != nil {
		authReqId.String = *args.RequestID
		authReqId.Valid = true
	}

	if args.RefreshTokenID != nil {
		refToken.String = *args.RefreshTokenID
		refToken.Valid = true
	}

	return qtx.CreateToken(ctx, query.CreateTokenParams{
		ID:             uuid.NewString(),
		ApplicationID:  args.ApplicationID,
		AuthRequestID:  authReqId,
		Scopes:         strings.Join(args.Scopes, " "),
		RefreshTokenID: refToken,
		Subject:        args.Subject,
		Audience:       strings.Join(args.Audience, " "),
		Expiration:     time.Now().Add(tokenLifetime).Unix(),
		CreatedAt:      time.Now().Unix(),
	})
}

type CreateRefreshTokenArgs struct {
	RefreshTokenID string
	ApplicationID  string
	RequestID      *string
	Scopes         []string
	Audience       []string
	AMR            []string
	Subject        string
	AuthTime       time.Time
}

func (t *tokenManager) CreateRefreshToken(ctx context.Context, tx *sql.Tx, args CreateRefreshTokenArgs) (*query.RefreshToken, error) {
	qtx := t.queries.WithTx(tx)

	// clientConfig, ok := t.clients.GetByID(args.ApplicationID)

	// if !ok {
	// 	return nil, errors.New("client not found")
	// }

	var amr = sql.NullString{}

	if len(args.AMR) > 0 {
		amr.String = strings.Join(args.AMR, " ")
		amr.Valid = true
	}

	return qtx.CreateRefreshToken(ctx, query.CreateRefreshTokenParams{
		ID:            args.RefreshTokenID,
		AuthTime:      args.AuthTime.Unix(),
		Amr:           amr,
		Audience:      strings.Join(args.Audience, " "),
		Scopes:        strings.Join(args.Scopes, " "),
		Subject:       args.Subject,
		Expiration:    time.Now().Add(time.Hour).Unix(),
		ApplicationID: args.ApplicationID,
		CreatedAt:     time.Now().Unix(),
	})
}

type ExchangeRefreshTokenArgs struct {
	TokenRequest op.TokenExchangeRequest
}

func (t *tokenManager) ExchangeRefreshToken(ctx context.Context, tx *sql.Tx, args ExchangeRefreshTokenArgs) (*query.Token, *query.RefreshToken, error) {

	applicationID := args.TokenRequest.GetClientID()
	authTime := args.TokenRequest.GetAuthTime()

	refreshTokenID := uuid.NewString()

	// create the new token
	accessToken, err := t.CreateToken(ctx, tx, CreateTokenArgs{
		ApplicationID:  applicationID,
		RequestID:      nil,
		Scopes:         args.TokenRequest.GetScopes(),
		Audience:       args.TokenRequest.GetAudience(),
		Subject:        args.TokenRequest.GetSubject(),
		RefreshTokenID: &refreshTokenID,
	})

	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := t.CreateRefreshToken(ctx, tx, CreateRefreshTokenArgs{
		RefreshTokenID: refreshTokenID,
		ApplicationID:  applicationID,
		RequestID:      nil,
		Scopes:         args.TokenRequest.GetScopes(),
		Audience:       args.TokenRequest.GetAudience(),
		AMR:            args.TokenRequest.GetAMR(),
		Subject:        args.TokenRequest.GetSubject(),
		AuthTime:       authTime,
	})

	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (t *tokenManager) IssueFromRequestToken(ctx context.Context, tx *sql.Tx, request op.TokenRequest, currentRefreshToken string) (*query.Token, *query.RefreshToken, error) {

	qtx := t.queries.WithTx(tx)

	crt, err := qtx.GetRefreshTokenByID(ctx, currentRefreshToken)

	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := qtx.CreateRefreshToken(ctx, query.CreateRefreshTokenParams{
		ID:            uuid.NewString(),
		AuthTime:      crt.AuthTime,
		Amr:           crt.Amr,
		Audience:      crt.Audience,
		Scopes:        crt.Scopes,
		Subject:       crt.Subject,
		Expiration:    time.Now().Add(time.Hour).Unix(),
		ApplicationID: crt.ApplicationID,
		CreatedAt:     time.Now().Unix(),
	})

	if err != nil {
		return nil, nil, err
	}

	accessToken, err := t.CreateToken(ctx, tx, CreateTokenArgs{
		ApplicationID:  crt.ApplicationID,
		RequestID:      nil,
		Scopes:         request.GetScopes(),
		Audience:       request.GetAudience(),
		Subject:        request.GetSubject(),
		RefreshTokenID: &refreshToken.ID,
	})

	if err != nil {
		return nil, nil, err
	}

	// TODO - delete the old refresh token and all access tokens based on it
	// deletes the refresh token and all access tokens which were issued based on this refresh token
	// delete(s.refreshTokens, currentRefreshToken)
	// for _, token := range s.tokens {
	// 	if token.RefreshTokenID == currentRefreshToken {
	// 		delete(s.tokens, token.ID)
	// 		break
	// 	}
	// }

	return accessToken, refreshToken, err

}

type IssueTokensArgs struct {
	TokenRequest      op.TokenRequest
	IssueRefreshToken bool
}

func (t *tokenManager) IssueTokens(ctx context.Context, tx *sql.Tx, args IssueTokensArgs) (*query.Token, *query.RefreshToken, error) {

	var refreshTokenID *string = nil

	if args.IssueRefreshToken {
		str := uuid.NewString()
		refreshTokenID = &str
	}

	var requestID *string = nil
	var applicationID string
	var authTime time.Time
	var amr []string = []string{}

	switch req := args.TokenRequest.(type) {
	case *query.AuthRequest:
		applicationID = req.ApplicationID
		authTime = req.GetAuthTime()
		requestID = &req.ID
		amr = req.GetAMR()
	case op.TokenExchangeRequest:
		applicationID = req.GetClientID()
		amr = req.GetAMR()
	case *query.RefreshToken:
		applicationID = req.ApplicationID
		authTime = req.GetAuthTime()
		amr = req.GetAMR()
	}

	// create the new token
	accessToken, err := t.CreateToken(ctx, tx, CreateTokenArgs{
		ApplicationID:  applicationID,
		RequestID:      requestID,
		Scopes:         args.TokenRequest.GetScopes(),
		Audience:       args.TokenRequest.GetAudience(),
		Subject:        args.TokenRequest.GetSubject(),
		RefreshTokenID: refreshTokenID,
	})

	if err != nil {
		return nil, nil, err
	}

	var refreshToken *query.RefreshToken = nil

	if args.IssueRefreshToken {
		refreshToken, err = t.CreateRefreshToken(ctx, tx, CreateRefreshTokenArgs{
			RefreshTokenID: *refreshTokenID,
			ApplicationID:  applicationID,
			RequestID:      requestID,
			Scopes:         args.TokenRequest.GetScopes(),
			Audience:       args.TokenRequest.GetAudience(),
			AMR:            amr,
			Subject:        args.TokenRequest.GetSubject(),
			AuthTime:       authTime,
		})

		if err != nil {
			return nil, nil, err
		}
	}

	return accessToken, refreshToken, err

}
