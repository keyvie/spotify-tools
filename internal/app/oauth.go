package app

import (
	"errors"
	"time"

	"github.com/keyvie/spotify-tools/internal/spotify"
	"github.com/keyvie/spotify-tools/internal/config"
)

func RefreshAccountToken(id string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	account, exists := cfg.Accounts[id]
	if !exists {
		return errors.New("Account not found")
	}
	
	if time.Now().Unix() < account.Expiry {
		return nil
	}

	// Error if missing stuff you know
	if cfg.Secrets.ClientID == "" {
		return errors.New("ClientID is not set in the configuration")
	}
	if cfg.Secrets.ClientSecret == "" {
		return errors.New("ClientSecret is not set in the configuration")
	}
	if cfg.Secrets.RedirectURI == "" {
		return errors.New("RedirectURI is not set in the configuration")
	}

	oauth, err := spotify.NewOAuth(
		cfg.Secrets.ClientID,
		cfg.Secrets.ClientSecret,
		cfg.Secrets.RedirectURI,
	)
	if err != nil {
		return err
	}
	
	tokenData, err := oauth.RefreshToken(account.RefreshToken)
	if err != nil {
		return err
	}
	account.AccessToken = tokenData.AccessToken
	if tokenData.RefreshToken != "" {
		account.RefreshToken = tokenData.RefreshToken
	}
	account.Expiry = tokenData.Expiry.Unix()
	
	if err := config.Lock(); err != nil {
		return err
	}
	defer config.Unlock()
	cfg, err = config.Load()
	if err != nil {
		return err
	}
	cfg.Accounts[id] = account
	config.Save(cfg)
	return nil
}