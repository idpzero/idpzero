package config

import (
	"bytes"
	"io"
	"os"

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

func LoadFromFile(doc *IDPConfiguration, path string) error {

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return Parse(doc, file)
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

// SetKey adds a new key to the configuration. If the key already exists, it will be replaced if replaceExisting is true.
// Returns true if the key was replaced, false if it was added.
func SetKey(doc *IDPConfiguration, key Key, replaceExisting bool) bool {

	for i, k := range doc.Server.Keys {
		if k.ID == key.ID {
			if replaceExisting {
				doc.Server.Keys[i] = key
				return true
			}
			break
		}
	}

	// insert at the beginning so it gets picked up as priority
	doc.Server.Keys = append([]Key{key}, doc.Server.Keys...)
	return false
}

// RemoveKey removes a key from the configuration if it exists. Returns true if removed, false if not found.
func RemoveKey(cfg *IDPConfiguration, kid string) bool {
	for i, key := range cfg.Server.Keys {
		if key.ID == kid {
			cfg.Server.Keys = append(cfg.Server.Keys[:i], cfg.Server.Keys[i+1:]...)
			return true
		}
	}

	return false
}

func KeyExists(doc *IDPConfiguration, id string) bool {

	for _, k := range doc.Server.Keys {
		if k.ID == id {
			return true
		}
	}

	return false
}
