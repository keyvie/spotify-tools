package main

import (
	"log"
	"os"
	"context"

	"github.com/keyvie/spotify-tools/cmd"
	"github.com/keyvie/spotify-tools/internal/app"
	"github.com/keyvie/spotify-tools/internal/config"

	"github.com/urfave/cli/v3"
)

func main() {
	if os.Getenv("CONFIG_PATH") != "" {
		config.Path = os.Getenv("CONFIG_PATH")
	}

	container := &app.Container{}
	RootCommand := &cli.Command{
		Name:  "spotify-tools",
		Usage: "My application",

		Commands: []*cli.Command{
			cmd.BuildAccountCommand(container),
			cmd.BuildAutomationCommand(container),
			cmd.BuildSecretCommand(container),
			cmd.BuildVersionCommand(container),
		},
	}

	if err := RootCommand.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
