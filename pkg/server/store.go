package server

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/google/uuid"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/idpzero/idpzero/pkg/dbg"
	"github.com/idpzero/idpzero/pkg/store/query"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

var _ op.Storage = &Storage{}

type Storage struct {
	logger       *slog.Logger
	lock         sync.Mutex
	configMgr    *configuration.ConfigurationManager
	server       atomic.Pointer[configuration.ServerConfig] // replaced whole on reload, never modified
	db           *sql.DB
	query        *query.Queries
	users        *users
	clients      *clients
	tokenManager *tokenManager
	//clients   map[string]configuration.ClientConfig
}

func NewStorage(logger *slog.Logger, configMgr *configuration.ConfigurationManager, db *sql.DB, queries *query.Queries, users *users, clients *clients) (*Storage, error) {

	store := &Storage{
		logger:       logger,
		configMgr:    configMgr,
		lock:         sync.Mutex{},
		db:           db,
		query:        queries,
		users:        users,
		clients:      clients,
		tokenManager: newTokenManager(queries, clients),
	}

	svrconf, err := configMgr.LoadConfiguration()

	if err != nil {
		return nil, err
	}

	store.setConfig(svrconf)

	// setup watching for changes
	configMgr.OnServerChanged(store.setConfig)

	return store, nil
}

// allow updating externally / on demand
func (s *Storage) setConfig(config *configuration.ServerConfig) {
	s.server.Store(config)
}

func (s *Storage) AuthRequestByCode(ctx context.Context, code string) (op.AuthRequest, error) {

	s.lock.Lock()
	defer s.lock.Unlock()

	return s.query.GetAuthRequestByAuthCode(ctx, sql.NullString{
		String: code,
		Valid:  true,
	})
}

func (s *Storage) AuthRequestByID(ctx context.Context, id string) (op.AuthRequest, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	return s.query.GetAuthRequestByID(ctx, id)
}

func (s *Storage) AuthorizeClientIDSecret(ctx context.Context, clientID string, clientSecret string) error {

	// validate the client with client secret. plain text is used
	for _, client := range s.server.Load().Clients {
		if client.ClientID == clientID {
			// Public clients authenticate via PKCE and have no secret; never allow
			// them to be authenticated through the client_secret path (which would
			// otherwise succeed on an empty secret match).
			if !client.RequiresClientSecret() {
				return fmt.Errorf("client '%s' is a public client and does not use a client secret", clientID)
			}
			if client.ClientSecret == clientSecret {
				return nil
			} else {
				return fmt.Errorf("incorrect client secret provided")
			}
		}
	}

	return fmt.Errorf("client not found")
}

func (s *Storage) CreateAccessAndRefreshTokens(ctx context.Context, request op.TokenRequest, currentRefreshToken string) (accessTokenID string, newRefreshTokenID string, expiration time.Time, err error) {

	dbg.Logger.Debug("Processing create and refresh token for request", slog.String("type", fmt.Sprintf("%T", request)), slog.String("currentRefreshToken", currentRefreshToken))

	// generate tokens via token exchange flow if request is relevant
	if teReq, ok := request.(op.TokenExchangeRequest); ok {
		tx, err := s.db.Begin()

		if err != nil {
			return "", "", time.Time{}, err
		}

		defer tx.Rollback()

		accessToken, refreshToken, err := s.tokenManager.ExchangeRefreshToken(ctx, tx, ExchangeRefreshTokenArgs{
			TokenRequest: teReq,
		})

		if err != nil {
			return "", "", time.Time{}, err
		}

		tx.Commit()
		return accessToken.ID, refreshToken.ID, time.Unix(accessToken.Expiration, 0), nil
	}

	// if currentRefreshToken is empty (Code Flow) we will have to create a new refresh token
	if currentRefreshToken == "" {

		tx, err := s.db.Begin()

		if err != nil {
			return "", "", time.Time{}, err
		}

		defer tx.Rollback()

		accessToken, refreshToken, err := s.tokenManager.IssueTokens(ctx, tx, IssueTokensArgs{
			TokenRequest:      request,
			IssueRefreshToken: true,
		})

		if err != nil {
			return "", "", time.Time{}, err
		}

		tx.Commit()
		return accessToken.ID, refreshToken.ID, time.Unix(accessToken.Expiration, 0), nil
	}

	tx, err := s.db.Begin()

	if err != nil {
		return "", "", time.Time{}, err
	}

	defer func() { _ = tx.Rollback() }()

	accessToken, refreshToken, err := s.tokenManager.IssueFromRequestToken(ctx, tx, request, currentRefreshToken)

	if err != nil {
		return "", "", time.Time{}, err
	}

	if err := tx.Commit(); err != nil {
		return "", "", time.Time{}, err
	}

	return accessToken.ID, refreshToken.ID, time.Unix(accessToken.Expiration, 0), nil
}

// CreateAccessToken implements op.Storage.
func (s *Storage) CreateAccessToken(ctx context.Context, request op.TokenRequest) (accessTokenID string, expiration time.Time, err error) {

	tx, err := s.db.Begin()

	if err != nil {
		return "", time.Time{}, err
	}

	defer tx.Rollback()

	token, _, err := s.tokenManager.IssueTokens(ctx, tx, IssueTokensArgs{
		TokenRequest:      request,
		IssueRefreshToken: false,
	})

	if err != nil {
		return "", time.Time{}, err
	}

	tx.Commit()

	return token.ID, time.Unix(token.Expiration, 0), nil
}

// CreateAuthRequest implements op.Storage.
func (s *Storage) CreateAuthRequest(ctx context.Context, authReq *oidc.AuthRequest, userID string) (op.AuthRequest, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(authReq.Prompt) == 1 && authReq.Prompt[0] == "none" {
		// With prompt=none, there is no way for the user to log in
		// so return error right away.
		return nil, oidc.ErrLoginRequired()
	}

	var maxAge int64 = 0
	if authReq.MaxAge != nil {
		maxAge = int64(*authReq.MaxAge)
	}

	req, err := s.query.CreateAuthRequest(ctx, query.CreateAuthRequestParams{
		ID:                  uuid.NewString(),
		ApplicationID:       authReq.ClientID,
		Scopes:              authReq.Scopes.String(),
		UserID:              userID,
		RedirectUri:         authReq.RedirectURI,
		State:               authReq.State,
		Nonce:               authReq.Nonce,
		Prompt:              authReq.Prompt.String(),
		MaxAuthAgeSeconds:   maxAge,
		LoginHint:           authReq.LoginHint,
		ResponseType:        string(authReq.ResponseType),
		ResponseMode:        string(authReq.ResponseMode),
		CodeChallenge:       authReq.CodeChallenge,
		CodeChallengeMethod: string(authReq.CodeChallengeMethod),
		Complete:            false,
		CreatedAt:           time.Now().Unix(),
		AuthenticatedAt:     0,
	})

	if err != nil {
		return nil, err
	}

	return req, nil
}

// DeleteAuthRequest implements op.Storage.
func (s *Storage) DeleteAuthRequest(ctx context.Context, id string) error {
	// TODO: Do we need to do anything here? Add soft delete?
	return nil
}

// GetClientByClientID implements op.Storage.
func (s *Storage) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {

	for _, c := range s.server.Load().Clients {
		if c.ClientID == clientID {
			return NewClient(*c), nil
		}
	}

	return nil, fmt.Errorf("no client registered with ID '%s'", clientID)
}

// GetKeyByIDAndClientID implements op.Storage.
//
// Clients have no registered public keys in the configuration, so a JWT
// assertion (private_key_jwt client authentication) cannot be verified.
func (s *Storage) GetKeyByIDAndClientID(ctx context.Context, keyID string, clientID string) (*jose.JSONWebKey, error) {
	return nil, fmt.Errorf("client '%s' has no registered keys: private_key_jwt is not supported", clientID)
}

// GetPrivateClaimsFromScopes implements op.Storage.
//
// These are the claims of JWT access tokens beyond the registered ones; the
// library has already removed the profile, email, phone and address scopes,
// so only custom claims are added.
func (s *Storage) GetPrivateClaimsFromScopes(ctx context.Context, userID string, clientID string, scopes []string) (map[string]any, error) {
	user, ok := s.users.GetByID(userID)
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	client, _ := s.clients.GetByID(clientID)
	return customClaims(user, client, scopes), nil
}

// GetRefreshTokenInfo implements op.Storage.
//
// Refresh tokens are handed out as their ID.
func (s *Storage) GetRefreshTokenInfo(ctx context.Context, clientID string, token string) (userID string, tokenID string, err error) {
	rt, err := s.query.GetRefreshTokenByID(ctx, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", op.ErrInvalidRefreshToken
		}
		return "", "", err
	}
	if rt.ApplicationID != clientID {
		return "", "", op.ErrInvalidRefreshToken
	}
	return rt.Subject, rt.ID, nil
}

// Health implements op.Storage.
func (s *Storage) Health(context.Context) error {
	if s.server.Load() == nil {
		return fmt.Errorf("no config loaded")
	}

	return nil
}

// KeySet implements op.Storage.
func (s *Storage) KeySet(ctx context.Context) ([]op.Key, error) {
	keys := make([]op.Key, 0)

	signingKeys, err := s.query.GetKeysByUse(ctx, KeyUseSig)

	if err != nil {
		return nil, err
	}

	for _, key := range signingKeys {
		keys = append(keys, &opPublicKey{key: *key})
	}

	return keys, nil
}

// RevokeToken implements op.Storage.
//
// tokenOrTokenID is an access token ID or a refresh token ID. Revoking a
// refresh token also revokes the access tokens issued with it. A token
// that does not exist is already revoked (RFC 7009 section 2.2).
func (s *Storage) RevokeToken(ctx context.Context, tokenOrTokenID string, userID string, clientID string) *oidc.Error {
	if token, err := s.query.GetTokenByID(ctx, tokenOrTokenID); err == nil {
		if token.ApplicationID != clientID {
			return oidc.ErrInvalidClient().WithDescription("token was not issued for this client")
		}
		if _, err := s.query.DeleteTokenByID(ctx, token.ID); err != nil {
			return oidc.ErrServerError().WithParent(err)
		}
		return nil
	} else if !errors.Is(err, sql.ErrNoRows) {
		return oidc.ErrServerError().WithParent(err)
	}

	rt, err := s.query.GetRefreshTokenByID(ctx, tokenOrTokenID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	if err != nil {
		return oidc.ErrServerError().WithParent(err)
	}
	if rt.ApplicationID != clientID {
		return oidc.ErrInvalidClient().WithDescription("token was not issued for this client")
	}

	tx, err := s.db.Begin()
	if err != nil {
		return oidc.ErrServerError().WithParent(err)
	}
	defer func() { _ = tx.Rollback() }()
	if err := s.tokenManager.DeleteRefreshToken(ctx, tx, rt.ID); err != nil {
		return oidc.ErrServerError().WithParent(err)
	}
	if err := tx.Commit(); err != nil {
		return oidc.ErrServerError().WithParent(err)
	}
	return nil
}

// SaveAuthCode implements op.Storage.
func (s *Storage) SaveAuthCode(ctx context.Context, id string, code string) error {
	count, err := s.query.UpdateAuthCode(ctx, query.UpdateAuthCodeParams{
		ID: id,
		AuthCode: sql.NullString{
			String: code,
			Valid:  true,
		},
	})

	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("auth request not found")
	}

	return nil
}

// SetIntrospectionFromToken implements op.Storage.
//
// An error makes the token inactive. clientID is the resource server that
// asks: it must be in the token's audience.
func (s *Storage) SetIntrospectionFromToken(ctx context.Context, introspection *oidc.IntrospectionResponse, tokenID string, subject string, clientID string) error {
	token, err := s.activeToken(ctx, tokenID)
	if err != nil {
		return err
	}
	if !slices.Contains(strings.Fields(token.Audience), clientID) {
		return fmt.Errorf("token is not valid for this client")
	}

	userInfo := new(oidc.UserInfo)
	if err := s.populateUserInfo(ctx, userInfo, subject, token.ApplicationID, strings.Fields(token.Scopes)); err != nil {
		return err
	}
	introspection.SetUserInfo(userInfo)
	introspection.Scope = strings.Fields(token.Scopes)
	introspection.ClientID = token.ApplicationID
	introspection.Audience = strings.Fields(token.Audience)
	introspection.Expiration = oidc.FromTime(time.Unix(token.Expiration, 0))
	introspection.IssuedAt = oidc.FromTime(time.Unix(token.CreatedAt, 0))
	introspection.TokenType = oidc.BearerToken
	return nil
}

// activeToken returns the access token tokenID unless it is missing or expired.
func (s *Storage) activeToken(ctx context.Context, tokenID string) (*query.Token, error) {
	token, err := s.query.GetTokenByID(ctx, tokenID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("token is invalid or has been revoked")
		}
		return nil, err
	}
	if time.Now().Unix() >= token.Expiration {
		return nil, fmt.Errorf("token has expired")
	}
	return token, nil
}

func (s *Storage) SetUserinfoFromScopes(ctx context.Context, userinfo *oidc.UserInfo, userID string, clientID string, scopes []string) error {
	return nil // not required
}

func (s *Storage) SetUserinfoFromRequest(ctx context.Context, userinfo *oidc.UserInfo, token op.IDTokenRequest, scopes []string) error {
	return s.populateUserInfo(ctx, userinfo, token.GetSubject(), token.GetClientID(), scopes)
}

// SetUserinfoFromToken implements op.Storage.
func (s *Storage) SetUserinfoFromToken(ctx context.Context, userinfo *oidc.UserInfo, tokenID string, subject string, origin string) error {

	token, err := s.activeToken(ctx, tokenID)

	if err != nil {
		return err
	}

	return s.populateUserInfo(ctx, userinfo, subject, token.ApplicationID, strings.Fields(token.Scopes))
}

// SignatureAlgorithms implements op.Storage.
func (s *Storage) SignatureAlgorithms(ctx context.Context) ([]jose.SignatureAlgorithm, error) {
	algs := make([]jose.SignatureAlgorithm, 0)

	signingKeys, err := s.query.GetKeysByUse(ctx, KeyUseSig)

	if err != nil {
		return nil, err
	}

	for _, key := range signingKeys {
		sa := jose.SignatureAlgorithm(key.Alg)
		algs = append(algs, sa)
	}

	return algs, nil
}

// SigningKey implements op.Storage.
func (s *Storage) SigningKey(ctx context.Context) (op.SigningKey, error) {

	key, err := s.query.GetKeyByID(ctx, signingKeyID)

	if err != nil {
		return nil, err
	}

	return &opPrivateKey{key: *key}, nil
}

// TerminateSession implements op.Storage.
//
// idpzero keeps no browser session: every authorization request shows the
// login page. There is nothing to terminate, so RP-initiated logout only
// redirects.
func (s *Storage) TerminateSession(ctx context.Context, userID string, clientID string) error {
	return nil
}

// TokenRequestByRefreshToken implements op.Storage.
func (s *Storage) TokenRequestByRefreshToken(ctx context.Context, refreshTokenID string) (op.RefreshTokenRequest, error) {
	rt, err := s.query.GetRefreshTokenByID(ctx, refreshTokenID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, op.ErrInvalidRefreshToken
		}
		return nil, err
	}
	if time.Now().Unix() >= rt.Expiration {
		return nil, fmt.Errorf("refresh token has expired")
	}
	return rt, nil
}

// ValidateJWTProfileScopes implements op.Storage.
//
// It is only used by the JWT authorization grant (RFC 7523), which the
// provider does not enable.
func (s *Storage) ValidateJWTProfileScopes(ctx context.Context, userID string, scopes []string) ([]string, error) {
	return nil, fmt.Errorf("the JWT authorization grant is not supported")
}

func (s *Storage) populateUserInfo(_ context.Context, userInfo *oidc.UserInfo, userID string, clientID string, scopes []string) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	user, ok := s.users.GetByID(userID)
	if !ok {
		return fmt.Errorf("user not found")
	}

	client, _ := s.clients.GetByID(clientID)

	// loop through the scopes and assign the claims based on those supported
	// within each of the well known scope types, and then custom scope configuration
	for _, scope := range scopes {
		switch scope {
		case oidc.ScopeOpenID:
			userInfo.Subject = user.Subject
		case oidc.ScopeEmail:
			if user.Claims.Email != nil {
				userInfo.Email = *user.Claims.Email
			}
			if user.Claims.EmailVerified != nil {
				userInfo.EmailVerified = oidc.Bool(*user.Claims.EmailVerified)
			}
		case oidc.ScopeProfile:
			if user.Claims.PreferredUsername != nil {
				userInfo.PreferredUsername = *user.Claims.PreferredUsername
			}

			if user.Claims.GivenName != nil {
				userInfo.GivenName = *user.Claims.GivenName
			}

			if user.Claims.FamilyName != nil {
				userInfo.FamilyName = *user.Claims.FamilyName
			}

			if user.Claims.MiddleName != nil {
				userInfo.MiddleName = *user.Claims.MiddleName
			}

			if user.Claims.Nickname != nil {
				userInfo.Nickname = *user.Claims.Nickname
			}

			if user.Claims.Name != nil {
				userInfo.Name = *user.Claims.Name
			} else if userInfo.GivenName != "" && userInfo.FamilyName != "" {
				userInfo.Name = fmt.Sprintf("%s %s", userInfo.GivenName, userInfo.FamilyName)
			}

			if user.Claims.UpdatedAt != nil {
				userInfo.UpdatedAt = oidc.Time(user.Claims.UpdatedAt.Unix())
			}

		case oidc.ScopePhone:
			if user.Claims.Phone != nil {
				userInfo.PhoneNumber = *user.Claims.Phone
			}
			if user.Claims.PhoneVerified != nil {
				userInfo.PhoneNumberVerified = *user.Claims.PhoneVerified
			}
		}
	}

	for claim, value := range customClaims(user, client, scopes) {
		userInfo.AppendClaims(claim, value)
	}
	return nil
}

// customClaims returns the user's custom claims released for scopes: those
// that no custom scope of the client lists come with openid, and those a
// custom scope lists only with that scope. It is nil when there are none.
func customClaims(user configuration.User, client configuration.ClientConfig, scopes []string) map[string]any {
	var claims map[string]any
	release := func(claim string) {
		if value, ok := user.Claims.Custom[claim]; ok {
			if claims == nil {
				claims = make(map[string]any)
			}
			claims[claim] = value
		}
	}

	scoped := make(map[string]bool)
	for _, listed := range client.CustomScopes {
		for _, claim := range listed {
			scoped[claim] = true
		}
	}

	for _, scope := range scopes {
		if scope == oidc.ScopeOpenID {
			for claim := range user.Claims.Custom {
				if !scoped[claim] {
					release(claim)
				}
			}
			continue
		}
		for _, claim := range client.CustomScopes[scope] {
			release(claim)
		}
	}
	return claims
}
