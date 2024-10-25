package idp

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/idpzero/idpzero/pkg/configuration"
)

func parseRSAPublicKey(key configuration.Key) (*rsa.PublicKey, error) {

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

func parseRSAKey(key configuration.Key) (*rsa.PrivateKey, *rsa.PublicKey, error) {

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
