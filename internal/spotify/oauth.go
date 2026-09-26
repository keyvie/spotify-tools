package spotify

import (
	"context"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

type OAuth struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

func NewOAuth(clientID, clientSecret, redirectURI string) (*OAuth, error) {
	return &OAuth{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}, nil
}

func (oauth *OAuth) GetAuthURL() string {
	scopes := []string{
		"user-read-playback-position",
		"user-top-read",
		"user-read-recently-played",

		"playlist-read-private",
		"playlist-read-collaborative",
		"playlist-modify-private",
		"playlist-modify-public",

		"user-follow-read",

		"user-library-read",
		"user-library-modify",
	}

	params := url.Values{}
	params.Add("client_id", oauth.ClientID)
	params.Add("response_type", "code")
	params.Add("redirect_uri", oauth.RedirectURI)
	params.Add("scope", strings.Join(scopes, " "))
	params.Add("show_dialog", "true")

	return "https://accounts.spotify.com/authorize?" + params.Encode()
}

func (oauth *OAuth) ExchangeCodeForToken(code string) (*oauth2.Token, error) {
	config := &oauth2.Config{
		ClientID:     oauth.ClientID,
		ClientSecret: oauth.ClientSecret,
		RedirectURL:  oauth.RedirectURI,
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://accounts.spotify.com/api/token",
		},
	}

	return config.Exchange(context.Background(), code)
}

func (oauth *OAuth) RefreshToken(refreshToken string) (*oauth2.Token, error) {
	config := &oauth2.Config{
		ClientID:     oauth.ClientID,
		ClientSecret: oauth.ClientSecret,
		RedirectURL:  oauth.RedirectURI,
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://accounts.spotify.com/api/token",
		},
	}
	return config.TokenSource(context.Background(), &oauth2.Token{
		RefreshToken: refreshToken,
	}).Token()
}