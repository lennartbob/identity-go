package oauth

import (
	"context"

	"golang.org/x/oauth2"
)

type RelaticsOAuthConfig struct {
	ClientID     string
	ClientSecret string
	Realm        string
	RedirectURI  string
}

func NewRelaticsOAuthConfig(cfg *RelaticsOAuthConfig) *oauth2.Config {
	authURL := "https://authenticate.relatics.com/auth/realms/" + cfg.Realm + "/protocol/openid-connect/auth"
	tokenURL := "https://authenticate.relatics.com/auth/realms/" + cfg.Realm + "/protocol/openid-connect/token"

	return &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURI,
		Scopes:       []string{"openid", "profile", "offline_access"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}
}

func GetRelaticsAuthURL(ctx context.Context, oauthConfig *oauth2.Config, state string) (string, error) {
	return oauthConfig.AuthCodeURL(state), nil
}

func ExchangeRelaticsCode(ctx context.Context, oauthConfig *oauth2.Config, code string) (*oauth2.Token, error) {
	return oauthConfig.Exchange(ctx, code)
}
