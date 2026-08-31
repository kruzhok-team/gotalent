package talentoauth2

import (
	"os"

	"golang.org/x/oauth2"
)

func Config(scopes []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("TALENT_CLIENT_ID"),
		ClientSecret: os.Getenv("TALENT_CLIENT_SECRET"),
		Endpoint:     Endpoint(),
		Scopes:       scopes,
	}
}

func Endpoint() oauth2.Endpoint {
	tokenURL := os.Getenv("TALENT_API_URL")
	if tokenURL == "" {
		tokenURL = "http://t2-api:8000/v2"
	}
	tokenURL += "/oauth/issue-token"
	return oauth2.Endpoint{
		TokenURL:  tokenURL,
		AuthStyle: oauth2.AuthStyleInParams,
	}
}
