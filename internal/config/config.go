package config

const DefaultPath = "config/config.yaml"
var Path = DefaultPath

type Config struct {
	Secrets  ConfigSecrets `yaml:"secrets"`
	Accounts []ConfigAccount `yaml:"accounts"`
}

type ConfigSecrets struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectURI  string `yaml:"redirect_uri"`
}

type ConfigAccount struct {
	Name         string `yaml:"name"`
	Username     string `yaml:"username"`
	Token        string `yaml:"token"`
	RefreshToken string `yaml:"refresh_token"`
	ExpiresAt    int64  `yaml:"expires_at"`
}
