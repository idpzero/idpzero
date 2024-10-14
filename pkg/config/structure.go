package config

import (
	"os"
	"path/filepath"

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

func ParseConfiguration(doc *Document, directory string) error {

	toRead := filepath.Join(directory, configFilename)

	data, err := os.ReadFile(toRead)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, doc)
}
