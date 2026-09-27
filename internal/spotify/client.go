package spotify

import (
	"context"
	"time"

	"fmt"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	http *resty.Client
}

func NewClient(accessToken string) *Client {
	http := resty.New().
		SetAuthToken(accessToken).
		SetBaseURL("https://api.spotify.com/v1").
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second)

	return &Client{http: http}
}

type User struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

func (c *Client) CurrentUser(ctx context.Context) (*User, error) {
	var user User
	resp, err := c.http.R().SetContext(ctx).SetResult(&user).Get("/me")
	if resp.IsError() {
		return nil, fmt.Errorf("spotify error %d: %s", resp.StatusCode(), resp.String())
	}
	return &user, err
}

type Playlist struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Owner struct {
		ID string `json:"id"`
	} `json:"owner"`
	Collaborative bool `json:"collaborative"`
}

func (c *Client) CurrentUsersPlaylists(ctx context.Context) ([]Playlist, error) {
	var all []Playlist
	url := "/me/playlists?limit=50"

	for url != "" {
		var page struct {
			Items []Playlist `json:"items"`
			Next  string     `json:"next"`
		}
		resp, err := c.http.R().SetContext(ctx).SetResult(&page).Get(url)
		if err != nil {
			return nil, err
		}
		if resp.IsError() {
			return nil, fmt.Errorf("spotify error %d: %s", resp.StatusCode(), resp.String())
		}

		all = append(all, page.Items...)
		url = page.Next
	}

	return all, nil
}

type PlaylistItem struct {
	AddedAt string `json:"added_at"`
	Item    struct {
		URI string `json:"uri"`
	} `json:"item"`
}

func (c *Client) GetPlaylistItems(ctx context.Context, playlistID string) ([]PlaylistItem, error) {
	var all []PlaylistItem
	url := "/playlists/" + playlistID + "/items?limit=50"

	for url != "" {
		var page struct {
			Items []PlaylistItem `json:"items"`
			Next  string         `json:"next"`
		}
		resp, err := c.http.R().SetContext(ctx).SetResult(&page).Get(url)
		if err != nil {
			return nil, err
		}
		if resp.IsError() {
			return nil, fmt.Errorf("spotify error %d: %s", resp.StatusCode(), resp.String())
		}

		all = append(all, page.Items...)
		url = page.Next
	}

	return all, nil
}

type reorderPayload struct {
	RangeStart   int `json:"range_start"`
	RangeLength  int `json:"range_length"`
	InsertBefore int `json:"insert_before"`
}

func (c *Client) ReorderPlaylistItems(ctx context.Context, playlistID string, rangeStart, rangeLength, insertBefore int) error {
	url := "/playlists/" + playlistID + "/items"

	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(reorderPayload{
			RangeStart:   rangeStart,
			RangeLength:  rangeLength,
			InsertBefore: insertBefore,
		}).
		Put(url)

	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("spotify error %d: %s", resp.StatusCode(), resp.String())
	}
	return nil
}