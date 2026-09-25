package cmd

import (
	"context"
	"fmt"

	"github.com/keyvie/spotify-tools/internal/app"
	"github.com/keyvie/spotify-tools/internal/config"
	"github.com/urfave/cli/v3"
)

func BuildSecretCommand(container *app.Container) *cli.Command {
	return &cli.Command{
		Name:  "secret",
		Usage: "Set the secrets values for the application",
		Commands: []*cli.Command{
			{
				Name:  "set",
				Usage: "Set a secret value",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					var typeChoice string
					fmt.Println("What do you want to set?")
					fmt.Println("1. Client ID")
					fmt.Println("2. Client Secret")
					fmt.Println("3. Redirect URI")
					fmt.Print("> ")
					fmt.Scanln(&typeChoice)
					fmt.Println("- - - - - - - - - - - - - - -")
					var value string
					fmt.Print("> ")
					fmt.Scanln(&value)
					fmt.Println("- - - - - - - - - - - - - - -")

					if err := config.Lock(); err != nil {
						return err
					}
					defer config.Unlock()
					cfg, err := config.Load()
					if err != nil {
						return err
					}

					switch typeChoice {
					case "1":
						cfg.Secrets.ClientID = value
					case "2":
						cfg.Secrets.ClientSecret = value
					case "3":
						cfg.Secrets.RedirectURI = value
					default:
						fmt.Println("Invalid option")
					}
					err = config.Save(cfg)
					if err != nil {
						return fmt.Errorf("Error saving config: %w", err)
					}
					fmt.Println("Secret updated successfully")
					return nil
				},
			},
		},
	}
}