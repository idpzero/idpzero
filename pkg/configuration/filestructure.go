package configuration

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
