package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

type MicrosoftOAuthConfig struct {
	ClientID     string
	ClientSecret string
	TenantID     string
	CallbackURL  string
}

type MicrosoftUserInfo struct {
	ID         string `json:"sub"`
	Email      string `json:"email"`
	GivenName  string `json:"given_name"`
	FamilyName string `json:"family_name"`
	Name       string `json:"name"`
}

type MicrosoftOAuthClient struct {
	config     *oauth2.Config
	httpClient *http.Client
}

func NewMicrosoftOAuthConfig(cfg *MicrosoftOAuthConfig) *MicrosoftOAuthClient {
	return &MicrosoftOAuthClient{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.CallbackURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     microsoft.AzureADEndpoint(cfg.TenantID),
		},
		httpClient: http.DefaultClient,
	}
}

func (c *MicrosoftOAuthClient) GetAuthURL(redirectURL, state, returnTo string) string {
	if redirectURL != "" {
		c.config.RedirectURL = redirectURL
	}
	return c.config.AuthCodeURL(state, oauth2.SetAuthURLParam("return_to", returnTo))
}

func (c *MicrosoftOAuthClient) ExchangeCode(ctx context.Context, code string) (*MicrosoftUserInfo, error) {
	if c.config.RedirectURL == "" {
		return nil, fmt.Errorf("redirect URL not configured")
	}

	token, err := c.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	client := c.config.Client(ctx, token)
	resp, err := client.Get("https://graph.microsoft.com/oidc/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to get user info: status %d, body: %s", resp.StatusCode, string(body))
	}

	var userInfo MicrosoftUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}
