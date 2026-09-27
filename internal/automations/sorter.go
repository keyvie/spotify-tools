package automations

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/keyvie/spotify-tools/internal/spotify"
)

const (
	OrderDateAddedNewestFirst = "Date Added (Newest First)"
	OrderDateAddedOldestFirst = "Date Added (Oldest First)"
)

func RunSorter(a *Automation) error {
	playlistID := a.Options["playlist_id"].(string)
	sortOrder := a.Options["order"].(string)
	account := a.Account

	client := spotify.NewClient(account.AccessToken)
	ctx := context.Background()

	tracks, err := client.GetPlaylistItems(ctx, playlistID)
	if err != nil {
		return err
	}

	target := make([]spotify.PlaylistItem, len(tracks))
	copy(target, tracks)

	switch sortOrder {
	case OrderDateAddedNewestFirst:
		sort.SliceStable(target, func(i, j int) bool {
			return target[i].AddedAt > target[j].AddedAt
		})
	case OrderDateAddedOldestFirst:
		sort.SliceStable(target, func(i, j int) bool {
			return target[i].AddedAt < target[j].AddedAt
		})
	default:
		return fmt.Errorf("unknown sort order: %s", sortOrder)
	}

	current := make([]string, len(tracks))
	for i, t := range tracks {
		current[i] = t.Item.URI
	}

	moves := 0

	for i := 0; i < len(target); i++ {
		uri := target[i].Item.URI
		if current[i] == uri {
			continue
		}

		j := -1
		for k := i; k < len(current); k++ {
			if current[k] == uri {
				j = k
				break
			}
		}
		if j == -1 {
			return fmt.Errorf("track %s not found", uri)
		}

		if err := client.ReorderPlaylistItems(ctx, playlistID, j, 1, i); err != nil {
			return err
		}
		moves++

		time.Sleep(500 * time.Millisecond)

		current = append(current[:j], current[j+1:]...)
		current = append(current[:i], append([]string{uri}, current[i:]...)...)
	}

	fmt.Printf("Sorted playlist %s (%d moves)\n", playlistID, moves)
	return nil
}