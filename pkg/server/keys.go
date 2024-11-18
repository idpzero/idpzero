package server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"
	"unicode"

	jose "github.com/go-jose/go-jose/v4"
	"github.com/idpzero/idpzero/pkg/store/query"
	"github.com/zitadel/oidc/v3/pkg/op"
)

const (
	KeyUseSig = "sig"
)

var (
	ErrUnsupportedKeyUse         = errors.New("unsupported key use")
	ErrContainsInvalidCharacters = errors.New("invalid characters found, expecting letters, numbers, - or _")
)

var (
	_ op.Key        = &opPublicKey{}  // make sure my type implements the interface
	_ op.SigningKey = &opPrivateKey{} // make sure my type implements the interface
)

type opPublicKey struct {
	key query.Key
}

func (s *opPublicKey) ID() string {
	return s.key.ID
}

func (s *opPublicKey) Algorithm() jose.SignatureAlgorithm {
	return jose.SignatureAlgorithm(s.key.Alg)
}

func (s *opPublicKey) Use() string {
	return s.key.Usage
}

func (s *opPublicKey) Key() any {

	parsed, err := parseRSAPublicKey(s.key.Alg, []byte(s.key.PublicKey))
	if err != nil {
		return err
	}
	return parsed

}

type opPrivateKey struct {
	key query.Key
}

func (s *opPrivateKey) SignatureAlgorithm() jose.SignatureAlgorithm {
	return jose.SignatureAlgorithm(s.key.Alg)
}

func (s *opPrivateKey) Key() any {

	priv, _, err := parseRSAKey(s.key.Alg, []byte(s.key.PrivateKey), []byte(s.key.PublicKey))

	if err != nil {
		return err
	}

	return priv
}

func (s *opPrivateKey) ID() string {
	return s.key.ID
}

func parseRSAPublicKey(alg string, data []byte) (*rsa.PublicKey, error) {

	if alg != "RS256" {
		return nil, errors.New("unsupported algorithm - expecting RS256")
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pubkey, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return nil, err
	}

	return pubkey.(*rsa.PublicKey), nil
}

func parseRSAKey(alg string, privData []byte, pubData []byte) (*rsa.PrivateKey, *rsa.PublicKey, error) {

	if alg != "RS256" {
		return nil, nil, errors.New("unsupported algorithm - expecting RS256")
	}

	block, _ := pem.Decode(privData)
	if block == nil {
		return nil, nil, errors.New("failed to parse PEM block containing the key")
	}

	privkey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, err
	}

	block, _ = pem.Decode(pubData)
	if block == nil {
		return nil, nil, errors.New("failed to parse PEM block containing the key")
	}

	pubkey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, nil, err
	}

	return privkey, pubkey.(*rsa.PublicKey), nil
}

// NewRSAKey creates a new RSA key pair with the given id and use based on the OIDC specifcation. The
// id should contain only letters, numbers, - or _
func newRSAKey(id string, use string) (query.CreateKeyParams, error) {

	if !validKeyCharacters(id) {
		return query.CreateKeyParams{}, ErrContainsInvalidCharacters
	}

	if use != KeyUseSig {
		return query.CreateKeyParams{}, ErrUnsupportedKeyUse
	}

	privkey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return query.CreateKeyParams{}, err
	}

	privkey_bytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkey_bytes,
		},
	)
	privKeyPem := string(privkey_pem)

	pubkey_bytes, err := x509.MarshalPKIXPublicKey(&privkey.PublicKey)
	if err != nil {
		return query.CreateKeyParams{}, err
	}
	pubkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkey_bytes,
		},
	)

	pubkey := string(pubkey_pem)

	key := query.CreateKeyParams{}
	key.ID = id
	key.Alg = "RS256"
	key.Usage = use
	key.PublicKey = pubkey
	key.PrivateKey = privKeyPem
	key.CreatedAt = time.Now().Unix()

	return key, nil
}

func validKeyCharacters(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '-' && r != '_' {
			return false
		}
	}
	return true
}
