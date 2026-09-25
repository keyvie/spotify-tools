package cmd

import (
	"context"
	"fmt"

	"github.com/keyvie/spotify-tools/internal/app"

	"github.com/urfave/cli/v3"
)

func BuildVersionCommand(container *app.Container) *cli.Command {
	return &cli.Command{
		Name:  "version",
		Usage: "Show version",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("Version 1.0.0")
			return nil
		},
	}
}
