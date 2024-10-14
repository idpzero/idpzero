package storage

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

type ConfigDocument struct {
	Server  ServerConfig   `yaml:"server"`
	Clients []ClientConfig `yaml:"clients"`
}

func parse(doc *ConfigDocument, path string) error {
	data, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, doc)
}
