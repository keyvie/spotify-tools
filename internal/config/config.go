package config

const DefaultPath = "config/config.yaml"
var Path = DefaultPath

type Config struct {
	Secrets  ConfigSecrets `yaml:"secrets"`
	Accounts map[string]ConfigAccount `yaml:"accounts"`
}

type ConfigSecrets struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectURI  string `yaml:"redirect_uri"`
}

type ConfigAccount struct {
	DisplayName  string `yaml:"display_name"`
	AccessToken  string `yaml:"access_token"`
	RefreshToken string `yaml:"refresh_token"`
	Expiry       int64  `yaml:"expiry"`
}
