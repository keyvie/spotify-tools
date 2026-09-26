package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/keyvie/spotify-tools/internal/app"
	"github.com/keyvie/spotify-tools/internal/config"
	"github.com/keyvie/spotify-tools/internal/spotify"

	"github.com/urfave/cli/v3"
)

func BuildAccountCommand(container *app.Container) *cli.Command {
	return &cli.Command{
		Name:  "account",
		Usage: "Manage accounts that are connected",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all connected accounts",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("Listing accounts...")
					return nil
				},
			},
			{
				Name:  "add",
				Usage: "Add a new account",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					cfg, err := config.Load()
					if err != nil {
						return err
					}
					secrets := cfg.Secrets

					// Error if missing stuff you know
					if secrets.ClientID == "" {
						return fmt.Errorf("ClientID is not set in the configuration")
					}
					if secrets.ClientSecret == "" {
						return fmt.Errorf("ClientSecret is not set in the configuration")
					}
					if secrets.RedirectURI == "" {
						return fmt.Errorf("RedirectURI is not set in the configuration")
					}

					// Build the OAuth connection thingy you know
					oauth, err := spotify.NewOAuth(
						secrets.ClientID,
						secrets.ClientSecret,
						secrets.RedirectURI,
					)
					if err != nil {
						return err
					}
					url := oauth.GetAuthURL()
					fmt.Println("Open and Authorize yourself:", url)
					fmt.Println("")
					fmt.Println("Enter the code you received after authorizing")
					fmt.Print("> ")
					var code string
					fmt.Scanln(&code)
					fmt.Println("")

					tokenData, err := oauth.ExchangeCodeForToken(code)
					if err != nil {
						return err
					}
					if tokenData.TokenType != "Bearer" {
						return errors.New("unexpected token type")
					}
					client := spotify.NewClient(tokenData.AccessToken)
					ctx, cancel := spotify.MakeContext()
					defer cancel()

					userData, err := client.CurrentUser(ctx)
					if err != nil {
						return err
					}

					if err := config.Lock(); err != nil {
						return err
					}
					defer config.Unlock()
					cfg, err = config.Load()
					if err != nil {
						return err
					}
					cfg.Accounts[userData.ID] = config.ConfigAccount{
						DisplayName:  userData.DisplayName,
						AccessToken:  tokenData.AccessToken,
						RefreshToken: tokenData.RefreshToken,
						Expiry:       tokenData.Expiry.Unix(),
					}
					config.Save(cfg)
					fmt.Printf("Added %s (%s) successfully\n", userData.ID, userData.DisplayName)
					return nil
				},
			},
			{
				Name:  "remove",
				Usage: "Remove a connected account",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					fmt.Println("Removing account...")
					return nil
				},
			},
		},
	}
}