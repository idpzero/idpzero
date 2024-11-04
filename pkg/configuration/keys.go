package configuration

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

// SetKey adds a new key to the configuration. If the key already exists, it will be replaced if replaceExisting is true.
// Returns true if the key was replaced, false if it was added.
func SetKey(doc *KeysConfiguration, key Key, replaceExisting bool) bool {

	for i, k := range doc.Keys {
		if k.ID == key.ID {
			if replaceExisting {
				doc.Keys[i] = key
				return true
			}
			break
		}
	}

	// insert at the beginning so it gets picked up as priority
	doc.Keys = append([]Key{key}, doc.Keys...)
	return false
}

// RemoveKey removes a key from the configuration if it exists. Returns true if removed, false if not found.
func RemoveKey(cfg *KeysConfiguration, kid string) bool {
	for i, key := range cfg.Keys {
		if key.ID == kid {
			cfg.Keys = append(cfg.Keys[:i], cfg.Keys[i+1:]...)
			return true
		}
	}

	return false
}

func KeyExists(doc *KeysConfiguration, id string) bool {

	for _, k := range doc.Keys {
		if k.ID == id {
			return true
		}
	}

	return false
}

func validKeyCharacters(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '-' && r != '_' {
			return false
		}
	}
	return true
}
