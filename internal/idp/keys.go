package idp

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/google/uuid"
)

func NewRSAKey() (*Key, error) {

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
	key.ID = uuid.New().String()
	key.Algorithm = "RS256"
	key.Data = map[string]string{
		"private": privKeyPem,
		"public":  pubkey,
	}

	return key, nil
}

func parseRSAPublicKey(key Key) (*rsa.PublicKey, error) {

	if key.Algorithm != "RS256" {
		return nil, errors.New("unsupported algorithm - expecting RS256")
	}

	pubkey_pem := key.Data["public"]
	block, _ := pem.Decode([]byte(pubkey_pem))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pubkey, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return nil, err
	}

	return pubkey.(*rsa.PublicKey), nil
}

func parseRSAKey(key Key) (*rsa.PrivateKey, *rsa.PublicKey, error) {

	if key.Algorithm != "RS256" {
		return nil, nil, errors.New("unsupported algorithm - expecting RS256")
	}

	privkey_pem := key.Data["private"]
	block, _ := pem.Decode([]byte(privkey_pem))
	if block == nil {
		return nil, nil, errors.New("failed to parse PEM block containing the key")
	}

	privkey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, nil, err
	}

	pubkey_pem := key.Data["public"]
	block, _ = pem.Decode([]byte(pubkey_pem))
	if block == nil {
		return nil, nil, errors.New("failed to parse PEM block containing the key")
	}

	pubkey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, nil, err
	}

	return privkey, pubkey.(*rsa.PublicKey), nil
}
