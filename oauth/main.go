package facebook

import (
	"golang.org/x/oauth2"
)

const (
	facebookClientID = ""
	facebookClientSecret = ""
	
	authorizeEndpoint = "https://www.facebook.com/dialog/oauth"
	tokenEndpoint     = "https://graph.facebook.com/oauth/access_token"
)

// GetConnect 接続する
func GetConnect() *oauth2.Config {
	config := &oauth2.Config {
		ClientID: facebookClientID,
		ClientSecret: facebookClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL: authorizeEndpoint,
			TokenURL: tokenEndpoint,
		},
		Scopes: []string{"email"},
		RedirectURL: "http://localhost:8080/facebook/callback",
	}
	
	return config
}

