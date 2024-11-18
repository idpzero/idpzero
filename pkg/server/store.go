package server

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/google/uuid"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/idpzero/idpzero/pkg/store/query"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

var _ op.Storage = &Storage{}

type Storage struct {
	logger    *slog.Logger
	lock      sync.Mutex
	configMgr *configuration.ConfigurationManager
	server    *configuration.ServerConfig
	query     *query.Queries
	users     *users
	clients   map[string]configuration.ClientConfig
}

func NewStorage(logger *slog.Logger, configMgr *configuration.ConfigurationManager, query *query.Queries, users *users) (*Storage, error) {

	store := &Storage{
		logger:    logger,
		configMgr: configMgr,
		lock:      sync.Mutex{},
		query:     query,
		users:     users,
		clients:   make(map[string]configuration.ClientConfig),
	}

	svrconf, err := configMgr.LoadServer()

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
	s.lock.Lock()
	defer s.lock.Unlock()

	s.server = config
	s.clients = make(map[string]configuration.ClientConfig)

	for _, client := range config.Clients {
		s.clients[client.ClientID] = client
	}
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
	for _, client := range s.server.Clients {
		if client.ClientID == clientID {
			if client.ClientSecret == clientSecret {
				return nil
			} else {
				return fmt.Errorf("incorrect client secret provided")
			}
		}
	}

	return fmt.Errorf("client not found")
}

// CreateAccessAndRefreshTokens implements op.Storage.
func (s *Storage) CreateAccessAndRefreshTokens(ctx context.Context, request op.TokenRequest, currentRefreshToken string) (accessTokenID string, newRefreshTokenID string, expiration time.Time, err error) {
	panic("***********************") // up to here!
	panic("unimplemented CreateAccessAndRefreshTokens")
}

// CreateAccessToken implements op.Storage.
func (s *Storage) CreateAccessToken(ctx context.Context, request op.TokenRequest) (accessTokenID string, expiration time.Time, err error) {

	var applicationID string
	var requestID sql.NullString
	switch req := request.(type) {
	case *query.AuthRequest:
		applicationID = req.ApplicationID
		requestID.String = req.ID
		requestID.Valid = true
	case op.TokenExchangeRequest:
		applicationID = req.GetClientID()
	}

	fmt.Println(applicationID, requestID, request.GetScopes(), request.GetSubject(), request.GetAudience())

	token, err := s.query.CreateToken(ctx, query.CreateTokenParams{
		ID:             uuid.NewString(),
		ApplicationID:  applicationID,
		AuthRequestID:  requestID,
		Scopes:         strings.Join(request.GetScopes(), " "),
		RefreshTokenID: "",
		Subject:        request.GetSubject(),
		Audience:       strings.Join(request.GetAudience(), " "),
		Expiration:     time.Now().Add(time.Hour).Unix(),
		CreatedAt:      time.Now().Unix(),
	})

	if err != nil {
		return "", time.Time{}, err
	}

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

	s.logger.Debug("created auth request", slog.String("id", req.ID))

	// finally, return the request (which implements the AuthRequest interface of the OP
	return req, nil
}

// DeleteAuthRequest implements op.Storage.
func (s *Storage) DeleteAuthRequest(ctx context.Context, id string) error {
	// TODO: Do we need to do anything here? Add soft delete?
	return nil
}

// GetClientByClientID implements op.Storage.
func (s *Storage) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {

	for _, c := range s.server.Clients {
		if c.ClientID == clientID {
			return NewClient(c), nil
		}
	}

	return nil, fmt.Errorf("no client registered with ID '%s'", clientID)
}

// GetKeyByIDAndClientID implements op.Storage.
func (s *Storage) GetKeyByIDAndClientID(ctx context.Context, keyID string, clientID string) (*jose.JSONWebKey, error) {
	panic("unimplemented GetKeyByIDAndClientID")
}

// GetPrivateClaimsFromScopes implements op.Storage.
func (s *Storage) GetPrivateClaimsFromScopes(ctx context.Context, userID string, clientID string, scopes []string) (map[string]any, error) {
	panic("unimplemented GetPrivateClaimsFromScopes")
}

// GetRefreshTokenInfo implements op.Storage.
func (s *Storage) GetRefreshTokenInfo(ctx context.Context, clientID string, token string) (userID string, tokenID string, err error) {
	panic("unimplemented GetRefreshTokenInfo")
}

// Health implements op.Storage.
func (s *Storage) Health(context.Context) error {
	if s.server == nil {
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
func (s *Storage) RevokeToken(ctx context.Context, tokenOrTokenID string, userID string, clientID string) *oidc.Error {
	panic("unimplemented RevokeToken")
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
func (s *Storage) SetIntrospectionFromToken(ctx context.Context, userinfo *oidc.IntrospectionResponse, tokenID string, subject string, clientID string) error {
	panic("unimplemented SetIntrospectionFromToken")
}

func (s *Storage) SetUserinfoFromScopes(ctx context.Context, userinfo *oidc.UserInfo, userID string, clientID string, scopes []string) error {
	return nil // not required
}

func (s *Storage) SetUserinfoFromRequest(ctx context.Context, userinfo *oidc.UserInfo, token op.IDTokenRequest, scopes []string) error {
	return s.populateUserInfo(ctx, userinfo, token.GetSubject(), token.GetClientID(), scopes)
}

// SetUserinfoFromToken implements op.Storage.
func (s *Storage) SetUserinfoFromToken(ctx context.Context, userinfo *oidc.UserInfo, tokenID string, subject string, origin string) error {

	token, err := s.query.GetTokenByID(ctx, tokenID)

	if err != nil {
		return err
	}

	return s.populateUserInfo(ctx, userinfo, subject, token.ApplicationID, strings.Split(token.Scopes, " "))
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
func (s *Storage) TerminateSession(ctx context.Context, userID string, clientID string) error {
	panic("unimplemented TerminateSession")
}

// TokenRequestByRefreshToken implements op.Storage.
func (s *Storage) TokenRequestByRefreshToken(ctx context.Context, refreshTokenID string) (op.RefreshTokenRequest, error) {
	panic("unimplemented TokenRequestByRefreshToken")
}

// ValidateJWTProfileScopes implements op.Storage.
func (s *Storage) ValidateJWTProfileScopes(ctx context.Context, userID string, scopes []string) ([]string, error) {
	panic("unimplemented ValidateJWTProfileScopes")
}

func (s *Storage) populateUserInfo(_ context.Context, userInfo *oidc.UserInfo, userID string, clientID string, scopes []string) (err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	user, ok := s.users.GetByID(userID)
	if !ok {
		return fmt.Errorf("user not found")
	}

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
		default:
			{
				client, ok := s.clients[clientID]
				if ok {
					if claims, ok := client.CustomScopes[scope]; ok {
						for _, claim := range claims {
							if value, ok := user.Claims.Custom[claim]; ok {
								userInfo.AppendClaims(claim, value)
							}
						}
					}
				}
			}
		}
	}
	return nil
}
