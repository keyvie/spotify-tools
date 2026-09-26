package spotify

import (
	"context"
	"time"
    "github.com/zmb3/spotify/v2"
    "golang.org/x/oauth2"
)

const RequestTimeout = 30 * time.Second

func MakeContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(
		context.Background(),
		RequestTimeout,
	)
}

func NewClient(accessToken string) *spotify.Client {
	token := &oauth2.Token{
		AccessToken: accessToken,
	}
	
	httpClient := oauth2.NewClient(
		context.Background(),
		oauth2.StaticTokenSource(token),
	) 
	
	return spotify.New(
		httpClient,
		spotify.WithRetry(true),
	) 
}