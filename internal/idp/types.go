package idp

import (
	jose "github.com/go-jose/go-jose/v4"
	"github.com/idpzero/idpzero/internal/config"
	"github.com/zitadel/oidc/v3/pkg/op"
)

var (
	_ op.Key = &opPublicKey{} // make sure my type implements the interface
)

type opPublicKey struct {
	key config.Key
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
