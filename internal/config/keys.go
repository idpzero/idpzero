package config

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"unicode"
)

const (
	KeyUseSig = "sig"
)

var (
	ErrUnsupportedKeyUse         = errors.New("unsupported key use")
	ErrContainsInvalidCharacters = errors.New("invalid characters found, expecting letters, numbers, - or _")
)

// NewRSAKey creates a new RSA key pair with the given id and use based on the OIDC specifcation. The
// id should contain only letters, numbers, - or _
func NewRSAKey(id string, use string) (*Key, error) {

	if !validKeyCharacters(id) {
		return nil, ErrContainsInvalidCharacters
	}

	if use != KeyUseSig {
		return nil, ErrUnsupportedKeyUse
	}

	privkey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return nil, err
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
		return nil, err
	}
	pubkey_pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkey_bytes,
		},
	)

	pubkey := string(pubkey_pem)

	key := &Key{}
	key.ID = id
	key.Algorithm = "RS256"
	key.Use = use
	key.Data = map[string]string{
		"private": privKeyPem,
		"public":  pubkey,
	}

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
