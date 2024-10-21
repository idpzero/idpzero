package configuration

import (
	"crypto/rsa"

	jose "github.com/go-jose/go-jose/v4"
)

type SigningKey struct {
	id        string
	algorithm jose.SignatureAlgorithm
	key       *rsa.PrivateKey
}

func (s *SigningKey) SignatureAlgorithm() jose.SignatureAlgorithm {
	return s.algorithm
}

func (s *SigningKey) Key() any {
	return s.key
}

func (s *SigningKey) ID() string {
	return s.id
}

type PublicKey struct {
	SigningKey
}

func (s *PublicKey) ID() string {
	return s.id
}

func (s *PublicKey) Algorithm() jose.SignatureAlgorithm {
	return s.algorithm
}

func (s *PublicKey) Use() string {
	return "sig"
}

func (s *PublicKey) Key() any {
	return &s.key.PublicKey
}
