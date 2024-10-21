package server

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Issuer    string `yaml:"issuer"`
	Port      int    `yaml:"port"`
	KeyPhrase string `yaml:"keyphrase"`
}

type ClientConfig struct {
	Name        string `yaml:"issuer"`
	ClientId    string `yaml:"port"`
	Secret      string `yaml:"keyphrase"`
	RedirectUri string `yaml:"redirectUri"`
}

type Document struct {
	Server  ServerConfig   `yaml:"server"`
	Clients []ClientConfig `yaml:"clients"`
}

// func (sc *ServerConfig) Key() [32]byte {
// 	if sc.KeyPhrase != "" {
// 		return sha256.Sum256([]byte(sc.KeyPhrase))
// 	} else {
// 		buf := make([]byte, 32)
// 		rand.Read(buf)
// 		var array32 [32]byte
// 		copy(array32[:], buf)
// 		return array32
// 	}
// }

// func (sc *ServerConfig) IssuerOrDefault() string {

// 	if sc.Issuer == "" {
// 		return "https://idpzero"
// 	}

// 	return sc.Issuer
// }

func Parse(doc *Document, file string) error {

	data, err := os.ReadFile(file)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, doc)
}
