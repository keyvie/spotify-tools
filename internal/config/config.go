package config

const DefaultPath = "config/config.yaml"
var Path = DefaultPath

type Config struct {
	Secrets     Secrets `yaml:"secrets"`
	Accounts    map[string]Account `yaml:"accounts"`
	Automations map[string]Automation `yaml:"automations"`
}

type Secrets struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectURI  string `yaml:"redirect_uri"`
}

type Account struct {
	ID           string `yaml:"id"`
	DisplayName  string `yaml:"display_name"`
	AccessToken  string `yaml:"access_token"`
	RefreshToken string `yaml:"refresh_token"`
	Expiry       int64  `yaml:"expiry"`
}

type Automation struct {
	Type     string         `yaml:"type"`
	Account  string         `yaml:"account"`
	Enabled  bool   		`yaml:"enabled"`
	Settings map[string]any `yaml:"settings"`
}