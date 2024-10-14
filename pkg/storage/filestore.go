package storage

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log/slog"
	"sync"
	"time"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/idpzero/idpzero/pkg/config"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/oidc/v3/pkg/op"
)

type StorageWithConfig interface {
	op.Storage
	Issuer() string
	ServerPort() int
	Secret() [32]byte
}

var _ StorageWithConfig = &Storage{}

type Storage struct {
	logger     *slog.Logger
	configFile string
	lock       sync.Mutex
	config     *config.Document

	codes        map[string]string
	authRequests map[string]*authReq
}

func NewStorage(logger *slog.Logger, configFilePath string) (StorageWithConfig, error) {

	store := &Storage{
		logger:       logger,
		lock:         sync.Mutex{},
		config:       &config.Document{},
		configFile:   configFilePath,
		codes:        make(map[string]string),
		authRequests: make(map[string]*authReq),
	}

	if store.configFile == "" {
		return nil, ErrConfigNotSupplied
	}

	logger.Info(fmt.Sprintf("Loading config from '%s'", store.configFile))
	err := config.ParseConfiguration(store.config, store.configFile)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (s *Storage) Issuer() string {
	if s.config.Server.Issuer == "" {
		return "idpzero"
	}
	return s.config.Server.Issuer
}

func (s *Storage) Secret() [32]byte {
	if s.config.Server.KeyPhrase != "" {
		return sha256.Sum256([]byte(s.config.Server.KeyPhrase))
	} else {
		s.logger.Warn("No keyphrase found in config file. Generating random secret.")
		buf := make([]byte, 32)
		rand.Read(buf)
		var array32 [32]byte
		copy(array32[:], buf)
		return array32
	}
}

func (s *Storage) ServerPort() int {
	if s.config.Server.Port == 0 {
		return 4379
	}

	return s.config.Server.Port
}

func (s *Storage) AuthRequestByCode(ctx context.Context, code string) (op.AuthRequest, error) {
	requestID, ok := func() (string, bool) {
		s.lock.Lock()
		defer s.lock.Unlock()
		requestID, ok := s.codes[code]
		return requestID, ok
	}()
	if !ok {
		return nil, fmt.Errorf("code invalid or expired")
	}
	return s.AuthRequestByID(ctx, requestID)
}

func (s *Storage) AuthRequestByID(ctx context.Context, id string) (op.AuthRequest, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	request, ok := s.authRequests[id]
	if !ok {
		return nil, fmt.Errorf("request not found")
	}
	return request, nil
}

func (s *Storage) AuthorizeClientIDSecret(ctx context.Context, clientID string, clientSecret string) error {

	for _, client := range s.config.Clients {
		if client.ClientId == clientID && client.Secret == clientSecret {
			return nil
		}
	}

	return fmt.Errorf("client not found")
}

// CreateAccessAndRefreshTokens implements op.Storage.
func (s *Storage) CreateAccessAndRefreshTokens(ctx context.Context, request op.TokenRequest, currentRefreshToken string) (accessTokenID string, newRefreshTokenID string, expiration time.Time, err error) {
	panic("unimplemented")
}

// CreateAccessToken implements op.Storage.
func (s *Storage) CreateAccessToken(context.Context, op.TokenRequest) (accessTokenID string, expiration time.Time, err error) {
	panic("unimplemented")
}

// CreateAuthRequest implements op.Storage.
func (s *Storage) CreateAuthRequest(context.Context, *oidc.AuthRequest, string) (op.AuthRequest, error) {
	panic("unimplemented")
}

// DeleteAuthRequest implements op.Storage.
func (s *Storage) DeleteAuthRequest(context.Context, string) error {
	panic("unimplemented")
}

// GetClientByClientID implements op.Storage.
func (s *Storage) GetClientByClientID(ctx context.Context, clientID string) (op.Client, error) {
	panic("unimplemented")
}

// GetKeyByIDAndClientID implements op.Storage.
func (s *Storage) GetKeyByIDAndClientID(ctx context.Context, keyID string, clientID string) (*jose.JSONWebKey, error) {
	panic("unimplemented")
}

// GetPrivateClaimsFromScopes implements op.Storage.
func (s *Storage) GetPrivateClaimsFromScopes(ctx context.Context, userID string, clientID string, scopes []string) (map[string]any, error) {
	panic("unimplemented")
}

// GetRefreshTokenInfo implements op.Storage.
func (s *Storage) GetRefreshTokenInfo(ctx context.Context, clientID string, token string) (userID string, tokenID string, err error) {
	panic("unimplemented")
}

// Health implements op.Storage.
func (s *Storage) Health(context.Context) error {
	return nil
}

// KeySet implements op.Storage.
func (s *Storage) KeySet(context.Context) ([]op.Key, error) {
	panic("unimplemented")
}

// RevokeToken implements op.Storage.
func (s *Storage) RevokeToken(ctx context.Context, tokenOrTokenID string, userID string, clientID string) *oidc.Error {
	panic("unimplemented")
}

// SaveAuthCode implements op.Storage.
func (s *Storage) SaveAuthCode(context.Context, string, string) error {
	panic("unimplemented")
}

// SetIntrospectionFromToken implements op.Storage.
func (s *Storage) SetIntrospectionFromToken(ctx context.Context, userinfo *oidc.IntrospectionResponse, tokenID string, subject string, clientID string) error {
	panic("unimplemented")
}

// SetUserinfoFromScopes implements op.Storage.
func (s *Storage) SetUserinfoFromScopes(ctx context.Context, userinfo *oidc.UserInfo, userID string, clientID string, scopes []string) error {
	panic("unimplemented")
}

// SetUserinfoFromToken implements op.Storage.
func (s *Storage) SetUserinfoFromToken(ctx context.Context, userinfo *oidc.UserInfo, tokenID string, subject string, origin string) error {
	panic("unimplemented")
}

// SignatureAlgorithms implements op.Storage.
func (s *Storage) SignatureAlgorithms(context.Context) ([]jose.SignatureAlgorithm, error) {
	panic("unimplemented")
}

// SigningKey implements op.Storage.
func (s *Storage) SigningKey(context.Context) (op.SigningKey, error) {
	panic("unimplemented")
}

// TerminateSession implements op.Storage.
func (s *Storage) TerminateSession(ctx context.Context, userID string, clientID string) error {
	panic("unimplemented")
}

// TokenRequestByRefreshToken implements op.Storage.
func (s *Storage) TokenRequestByRefreshToken(ctx context.Context, refreshTokenID string) (op.RefreshTokenRequest, error) {
	panic("unimplemented")
}

// ValidateJWTProfileScopes implements op.Storage.
func (s *Storage) ValidateJWTProfileScopes(ctx context.Context, userID string, scopes []string) ([]string, error) {
	panic("unimplemented")
}
