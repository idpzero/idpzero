package idp

import (
	jose "github.com/go-jose/go-jose/v4"
	"github.com/zitadel/oidc/v3/pkg/op"
)

var (
	_ op.Key = &opKey{} // make sure my type implements the interface
)

type opKey struct {
	key Key
}

func (s *opKey) ID() string {
	return s.key.ID
}

func (s *opKey) Algorithm() jose.SignatureAlgorithm {
	return jose.SignatureAlgorithm(s.key.Algorithm)
}

func (s *opKey) Use() string {
	return s.key.Use
}

func (s *opKey) Key() any {

	if s.key.Use == "sig" {
		parsed, err := parseRSAPublicKey(s.key)
		if err != nil {
			return err
		}
		return parsed
	}

	return nil
}
