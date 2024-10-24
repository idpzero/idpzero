package configuration

import (
	"bytes"
	"io"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Port      int    `yaml:"port"`
	KeyPhrase string `yaml:"keyphrase"`
	Keys      []Key  `yaml:"keys"`
}

type Key struct {
	ID        string            `yaml:"id"`
	Algorithm string            `yaml:"algorithm"`
	Use       string            `yaml:"use"`
	Data      map[string]string `yaml:"data"`
}

type ClientConfig struct {
	Name        string `yaml:"issuer"`
	ClientId    string `yaml:"port"`
	Secret      string `yaml:"keyphrase"`
	RedirectUri string `yaml:"redirectUri"`
}

type IDPConfiguration struct {
	Server  ServerConfig   `yaml:"server"`
	Clients []ClientConfig `yaml:"clients"`
}

func Parse(reader io.Reader) (*IDPConfiguration, error) {
	doc := &IDPConfiguration{}
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)

	if err != nil {
		return nil, err
	}

	if yaml.Unmarshal(buf.Bytes(), doc); err != nil {
		return nil, err
	}

	return doc, nil
}

// func Save(doc *IDPConfiguration, path string) error {

// 	data, err := yaml.Marshal(*doc)

// 	if err != nil {
// 		return err
// 	}

// 	return os.WriteFile(path, data, 0644)
// }
