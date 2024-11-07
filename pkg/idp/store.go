package idp

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

var _ op.Storage = &Storage{}

type Storage struct {
	logger    *slog.Logger
	lock      sync.Mutex
	configMgr *configuration.ConfigurationManager
	server    *configuration.ServerConfig
	keys      *configuration.KeysConfiguration
}

func NewStorage(logger *slog.Logger, configMgr *configuration.ConfigurationManager) (*Storage, error) {

	store := &Storage{
		logger:    logger,
		configMgr: configMgr,
		lock:      sync.Mutex{},
	}

	keys, err := configMgr.LoadKeys()

	if err != nil {
		return nil, err
	}

	store.keys = keys

	svrconf, err := configMgr.LoadServer()

	if err != nil {
		return nil, err
	}

	store.server = svrconf

	// setup watching for changes
	configMgr.OnKeysChanged(store.setKeys)
	configMgr.OnServerChanged(store.setConfig)

	return store, nil
}

// allow updating externally / on demand
func (s *Storage) setConfig(config *configuration.ServerConfig) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.server = config
}

func (s *Storage) setKeys(keys *configuration.KeysConfiguration) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.keys = keys
}

func (s *Storage) AuthRequestByCode(ctx context.Context, code string) (op.AuthRequest, error) {
	panic("unimplemented AuthRequestByCode")
	// requestID, ok := func() (string, bool) {
	// 	s.lock.Lock()
	// 	defer s.lock.Unlock()
	// 	requestID, ok := s.codes[code]
	// 	return requestID, ok
	// }()
	// if !ok {
	// 	return nil, fmt.Errorf("code invalid or expired")
	// }
	// return s.AuthRequestByID(ctx, requestID)
}

func (s *Storage) AuthRequestByID(ctx context.Context, id string) (op.AuthRequest, error) {
	panic("unimplemented AuthRequestByID")
	// s.lock.Lock()
	// defer s.lock.Unlock()
	// request, ok := s.authRequests[id]
	// if !ok {
	// 	return nil, fmt.Errorf("request not found")
	// }
	// return request, nil
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
	panic("unimplemented CreateAccessAndRefreshTokens")
}

// CreateAccessToken implements op.Storage.
func (s *Storage) CreateAccessToken(context.Context, op.TokenRequest) (accessTokenID string, expiration time.Time, err error) {
	panic("unimplemented CreateAccessToken")
}

// CreateAuthRequest implements op.Storage.
func (s *Storage) CreateAuthRequest(context.Context, *oidc.AuthRequest, string) (op.AuthRequest, error) {
	panic("****************************")
	panic("unimplemented CreateAuthRequest")
}

// DeleteAuthRequest implements op.Storage.
func (s *Storage) DeleteAuthRequest(context.Context, string) error {
	panic("unimplemented DeleteAuthRequest")
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
func (s *Storage) KeySet(context.Context) ([]op.Key, error) {
	keys := make([]op.Key, 0)

	for _, key := range s.keys.Keys {
		if key.Use == "sig" {
			keys = append(keys, &opPublicKey{key: key})
		}
	}

	return keys, nil
}

// RevokeToken implements op.Storage.
func (s *Storage) RevokeToken(ctx context.Context, tokenOrTokenID string, userID string, clientID string) *oidc.Error {
	panic("unimplemented RevokeToken")
}

// SaveAuthCode implements op.Storage.
func (s *Storage) SaveAuthCode(context.Context, string, string) error {
	panic("unimplemented SaveAuthCode")
}

// SetIntrospectionFromToken implements op.Storage.
func (s *Storage) SetIntrospectionFromToken(ctx context.Context, userinfo *oidc.IntrospectionResponse, tokenID string, subject string, clientID string) error {
	panic("unimplemented SetIntrospectionFromToken")
}

// SetUserinfoFromScopes implements op.Storage.
func (s *Storage) SetUserinfoFromScopes(ctx context.Context, userinfo *oidc.UserInfo, userID string, clientID string, scopes []string) error {
	panic("unimplemented SetUserinfoFromScopes")
}

// SetUserinfoFromToken implements op.Storage.
func (s *Storage) SetUserinfoFromToken(ctx context.Context, userinfo *oidc.UserInfo, tokenID string, subject string, origin string) error {
	panic("unimplemented SetUserinfoFromToken")
}

// SignatureAlgorithms implements op.Storage.
func (s *Storage) SignatureAlgorithms(context.Context) ([]jose.SignatureAlgorithm, error) {
	algs := make([]jose.SignatureAlgorithm, 0)

	for _, key := range s.keys.Keys {
		if key.Use == "sig" {
			sa := jose.SignatureAlgorithm(key.Algorithm)
			algs = append(algs, sa)
		}
	}

	return algs, nil
}

// SigningKey implements op.Storage.
func (s *Storage) SigningKey(context.Context) (op.SigningKey, error) {

	for _, key := range s.keys.Keys {
		if key.Use == "sig" {
			return &opPrivateKey{key: key}, nil
		}
	}

	return nil, fmt.Errorf("no signing key found")
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
