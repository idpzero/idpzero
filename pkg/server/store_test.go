package server

import (
	"context"
	"testing"

	"github.com/idpzero/idpzero/pkg/configuration"
)

func newSecretTestStorage() *Storage {
	return &Storage{
		server: &configuration.ServerConfig{
			Clients: []*configuration.ClientConfig{
				{ClientID: "web", AuthMethod: "client_secret_basic", ClientSecret: "topsecret"},
				{ClientID: "spa", AuthMethod: "none"},
			},
		},
	}
}

func TestAuthorizeClientIDSecret(t *testing.T) {
	s := newSecretTestStorage()
	ctx := context.Background()

	t.Run("confidential correct secret", func(t *testing.T) {
		if err := s.AuthorizeClientIDSecret(ctx, "web", "topsecret"); err != nil {
			t.Fatalf("expected success, got %v", err)
		}
	})

	t.Run("confidential wrong secret", func(t *testing.T) {
		if err := s.AuthorizeClientIDSecret(ctx, "web", "wrong"); err == nil {
			t.Fatal("expected error for wrong secret, got nil")
		}
	})

	t.Run("public client rejected even with empty secret", func(t *testing.T) {
		// Guards against an empty-secret match succeeding for a public client.
		if err := s.AuthorizeClientIDSecret(ctx, "spa", ""); err == nil {
			t.Fatal("expected public client to be rejected on secret auth, got nil")
		}
	})

	t.Run("unknown client", func(t *testing.T) {
		if err := s.AuthorizeClientIDSecret(ctx, "nope", "x"); err == nil {
			t.Fatal("expected error for unknown client, got nil")
		}
	})
}
