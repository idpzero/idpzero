package idp

import (
	"bytes"
	"io"
	"os"

	"gopkg.in/yaml.v2"
)

type ServerConfig struct {
	Issuer      string `yaml:"issuer"`
	Port        int    `yaml:"port"`
	KeyPhrase   string `yaml:"keyphrase"`
	SigningKeys []Key  `yaml:"signingKeys"`
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

func Parse(doc *IDPConfiguration, reader io.Reader) error {

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(buf.Bytes(), doc)
}

func Save(doc *IDPConfiguration, path string) error {

	data, err := yaml.Marshal(*doc)

	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
