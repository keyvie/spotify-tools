package cmd

import (
	"context"

	"github.com/keyvie/spotify-tools/internal/app"
	"github.com/keyvie/spotify-tools/internal/config"
	"github.com/urfave/cli/v3"
)


func BuildDaemonCommand(container *app.Container) *cli.Command {
	return &cli.Command{
		Name:  "daemon",
		Usage: "Run the background code for the cli",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			_, err := config.Load()
			if err != nil {
				return err
			}
			return nil
		},
	}
}