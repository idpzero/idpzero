package idp

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

var _ op.Storage = &Storage{}

type Storage struct {
	logger *slog.Logger
	lock   sync.Mutex
	config *IDPConfiguration
}

func NewStorage(logger *slog.Logger) (*Storage, error) {

	store := &Storage{
		logger: logger,
		lock:   sync.Mutex{},
	}

	return store, nil
}

// allow updating externally / on demand
func (s *Storage) SetConfig(config *IDPConfiguration) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.config = config
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
	panic("unimplemented AuthorizeClientIDSecret")
	// for _, client := range s.config.Clients {
	// 	if client.ClientId == clientID && client.Secret == clientSecret {
	// 		return nil
	// 	}
	// }

	// return fmt.Errorf("client not found")
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
	panic("unimplemented CreateAuthRequest")
}

// DeleteAuthRequest implements op.Storage.
func (s *Storage) DeleteAuthRequest(context.Context, string) error {
	panic("unimplemented DeleteAuthRequest")
}

// GetClientByClientID implements op.Storage.
func (s *Storage) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {
	panic("unimplemented GetClientByClientID")
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
	if s.config == nil {
		return fmt.Errorf("no config loaded")
	}

	return nil
}

// KeySet implements op.Storage.
func (s *Storage) KeySet(context.Context) ([]op.Key, error) {
	keys := make([]op.Key, 0, len(s.config.Server.SigningKeys))

	for _, key := range s.config.Server.SigningKeys {
		if key.Use == "sig" {
			keys = append(keys, &opKey{key: key})
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

	for _, key := range s.config.Server.SigningKeys {
		if key.Use == "sig" {
			sa := jose.SignatureAlgorithm(key.Algorithm)
			algs = append(algs, sa)
		}
	}

	return algs, nil
}

// SigningKey implements op.Storage.
func (s *Storage) SigningKey(context.Context) (op.SigningKey, error) {
	panic("unimplemented SigningKey")
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
