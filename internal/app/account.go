package app

import (
	"errors"
	"fmt"

	"github.com/keyvie/spotify-tools/internal/config"
	"github.com/urfave/cli/v3"
)

var AccountFlag = &cli.StringFlag{
	Name:    "account",
	Aliases: []string{"acc"},
	Usage:   "Account to use for the command",
}

func GetAccount(id string) (*config.Account, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	if len(cfg.Accounts) == 0 {
		return nil, errors.New("no accounts configured")
	}

	if id == "" {
		if len(cfg.Accounts) == 1 {
			for _, data := range cfg.Accounts {
				account := data
				return &account, nil
			}
		}

		return nil, errors.New("multiple accounts configured, specify one with --account")
	}

	if account, exists := cfg.Accounts[id]; exists {
		return &account, nil
	}

	return nil, fmt.Errorf("account with ID %s does not exist", id)
}