package spotify

import (
	"net/http"
	"net/url"
	"strings"
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
		"user-read-playback-position", // Playback position
		"user-top-read",               // Top tracks and artists
		"user-read-recently-played",   // Listening history

		"playlist-read-private", // Read private playlists
		"playlist-read-collaborative",
		"playlist-modify-private", // Modify private playlists
		"playlist-modify-public",  // Modify public playlists

		"user-follow-read", // Read followed artists/users

		"user-library-read",   // Read saved tracks and albums
		"user-library-modify", // Modify saved tracks and albums
	}
	params := url.Values{}
	params.Add("client_id", oauth.ClientID)
	params.Add("response_type", "code")
	params.Add("redirect_uri", oauth.RedirectURI)
	params.Add("scope", strings.Join(scopes, " "))
	params.Add("show_dialog", "true")
	return "https://accounts.spotify.com/authorize?" + params.Encode()
}

func (oauth *OAuth) ExchangeCodeForToken(code string) (string, error) {
	_, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", nil)
	if err != nil {
		return "", err
	}
	return "<access_token>", nil
}
