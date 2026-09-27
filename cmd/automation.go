package cmd

import (
	"context"
	"fmt"
	"errors"

	"github.com/keyvie/spotify-tools/internal/app"
	"github.com/keyvie/spotify-tools/internal/automations"
	"github.com/keyvie/spotify-tools/internal/config"
	"github.com/keyvie/spotify-tools/internal/spotify"

	"github.com/urfave/cli/v3"
	"github.com/manifoldco/promptui"
)

func SelectAPlaylist(ctx context.Context, account *config.Account) (*spotify.Playlist, error) {
	client := spotify.NewClient(account.AccessToken)

	currentUser, err := client.CurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch current user: %w", err)
	}

	playlists, err := client.CurrentUsersPlaylists(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch playlists: %w", err)
	}

	if len(playlists) == 0 {
		return nil, errors.New("no playlists found")
	}

	var items []spotify.Playlist
	for _, p := range playlists {
		if p.Owner.ID == currentUser.ID || p.Collaborative {
			items = append(items, p)
		}
	}

	if len(items) == 0 {
		return nil, errors.New("no modifiable playlists found")
	}

	prompt := promptui.Select{
		Label: "Select a playlist",
		Items: items,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   "> {{ .Name }}",
			Inactive: "  {{ .Name }}",
			Selected: "{{ \"v\" | green }} {{ .Name }}",
		},
	}
	idx, _, err := prompt.Run()
	if err != nil {
		return nil, fmt.Errorf("Prompt failed: %w", err)
	}

	return &items[idx], nil
}

func BuildAutomationCommand(container *app.Container) *cli.Command {
	return &cli.Command{
		Name:  "automation",
		Usage: "Manage automation tasks",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all automation tasks",
				Flags: []cli.Flag{app.AccountFlag},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					account, err := app.GetAccount(cmd.String("account"))
					if err != nil {
						return err
					}
					cfg, err := config.Load()
					if err != nil {
						return err
					}

					for name, automation := range cfg.Automations {
						if automation.AccountID == account.ID {
							fmt.Printf("Name: %s, Type: %s, Enabled: %t\n", name, automation.Type, automation.Enabled)
						}
					}
				
					return nil
				},
			},
			{
				Name:  "add",
				Usage: "Add a new automation task",
				Flags: []cli.Flag{app.AccountFlag},
				ArgsUsage: "<name>",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "name"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					account, err := app.GetAccount(cmd.String("account"))
					if err != nil {
						return err
					}

					name := cmd.StringArg("name")
					if name == "" {
						return errors.New("automation name is required")
					}

					cfg, err := config.Load()
					if err != nil {
						return err
					}
					if _, exists := cfg.Automations[name]; exists {
						return fmt.Errorf("automation with name %s already exists", name)
					}

					automationTypes := []string{"sorter"}
					prompt := promptui.Select{
						Label: "Select automation type",
						Items: automationTypes,
					}
					_, automationType, err := prompt.Run()
					if err != nil {
						return fmt.Errorf("Prompt failed: %w", err)
					}

					switch automationType {
						case "sorter":
							playlist, err := SelectAPlaylist(ctx, account)
							if err != nil {
								return err
							}

							sortOrders := []string{
								automations.OrderDateAddedNewestFirst,
								automations.OrderDateAddedOldestFirst,
							}
							prompt := promptui.Select{
								Label: "Select an order to sort the playlist in",
								Items: sortOrders,
							}
							_, sortOrder, err := prompt.Run()
							if err != nil {
								return fmt.Errorf("Prompt failed: %w", err)
							}

							if err := config.Lock(); err != nil {
								return err
							}
							defer config.Unlock()
							cfg, err = config.Load()
							if err != nil {
								return err
							}

							cfg.Automations[name] = config.Automation{
								Type:    	"sorter",
								AccountID: 	account.ID,
								Enabled: 	true,
								Interval: 	3600,
								LastRunAt: 	0,
								Options: 	map[string]any{
									"playlist_id": playlist.ID,
									"order": sortOrder,
								},
							}
							if err := config.Save(cfg); err != nil {
								return err
							}
							break
						
						default:
							return fmt.Errorf("Unknown automation type: %s", automationType)
					}

					fmt.Println("Automation task added")
					return nil
				},
			},
			{
				Name:  "run",
				Usage: "Run an automation task immediately",
				ArgsUsage: "<name>",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "name"},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					name := cmd.StringArg("name")
					if name == "" {
						return errors.New("automation name is required")
					}

					automation, err := automations.Get(name)
					if err != nil {
						return fmt.Errorf("failed to get automation %s: %w", name, err)
					}

					_, err = app.GetAccount(automation.Account.ID)
					if err != nil {
						return err
					}

					if err := automation.Run(true); err != nil {
						return fmt.Errorf("failed to run automation %s: %w", name, err)
					}
					fmt.Println("Automation task executed")
					return nil
				},
			},
		},
	}
}
