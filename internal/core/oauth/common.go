package oauth

import (
	"github.com/vondr/identity-go/internal/core"
)

var redirectURL string

func SetRedirectURL(url string) {
	redirectURL = url
}

func getRedirectURL() string {
	if redirectURL != "" {
		return redirectURL
	}
	return core.GetConfig().OAuthCallbackURL
}
