package core

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	AppName     string `mapstructure:"APP_NAME"`
	Environment string `mapstructure:"ENVIRONMENT"`
	RootDomain  string `mapstructure:"ROOT_DOMAIN"`
	AdminToken  string `mapstructure:"ADMIN_TOKEN"`

	DatabaseURL string `mapstructure:"DATABASE_URL"`
	KeyDBURL    string `mapstructure:"KEYDB_URL"`

	MicrosoftClientID     string `mapstructure:"MS_CLIENT_ID"`
	MicrosoftClientSecret string `mapstructure:"MS_CLIENT_SECRET"`
	MicrosoftTenantID     string `mapstructure:"MS_TENANT_ID"`

	RelaticsClientID     string `mapstructure:"RELATICS_CLIENT_ID"`
	RelaticsClientSecret string `mapstructure:"RELATICS_CLIENT_SECRET"`
	RelaticsRedirectURI  string `mapstructure:"RELATICS_REDIRECT_URI"`
	RelaticsRealm        string `mapstructure:"RELATICS_REALM"`

	PostLoginRedirectURL string `mapstructure:"POST_LOGIN_REDIRECT_URL"`
	ErrorLoginRedirect   string `mapstructure:"ERROR_LOGIN_REDIRECT"`
	OAuthCallbackURL     string `mapstructure:"OAUTH_CALLBACK_URL"`
	AuthLoginURL         string `mapstructure:"AUTH_LOGIN_URL"`
	CookieDomain         string `mapstructure:"COOKIE_DOMAIN"`
	CookieSecure         bool   `mapstructure:"COOKIE_SECURE"`
	CookieSameSite       string `mapstructure:"COOKIE_SAMESITE"`
	SessionTTLDays       int    `mapstructure:"SESSION_TTL_DAYS"`
	SessionSecretKey     string `mapstructure:"SESSION_SECRET_KEY"`

	SystemEmailsRaw string `mapstructure:"SYSTEM_EMAILS"`

	CORSOriginsRaw string `mapstructure:"CORS_ORIGINS"`

	MicrosoftEmailTenantID     string `mapstructure:"MICROSOFT_EMAIL_TENANT_ID"`
	MicrosoftEmailClientID     string `mapstructure:"MICROSOFT_EMAIL_CLIENT_ID"`
	MicrosoftEmailClientSecret string `mapstructure:"MICROSOFT_EMAIL_CLIENT_SECRET"`
	MicrosoftEmailSender       string `mapstructure:"MICROSOFT_EMAIL_SENDER"`

	GeoIPDBPath string `mapstructure:"GEOIP_DB_PATH"`
}

var settings *Config

func LoadConfig() (*Config, error) {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")

	viper.AutomaticEnv()
	viper.SetEnvPrefix("")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	config := &Config{
		AppName:                    viper.GetString("APP_NAME"),
		Environment:                viper.GetString("ENVIRONMENT"),
		RootDomain:                 viper.GetString("ROOT_DOMAIN"),
		AdminToken:                 viper.GetString("ADMIN_TOKEN"),
		DatabaseURL:                viper.GetString("DATABASE_URL"),
		KeyDBURL:                   viper.GetString("KEYDB_URL"),
		MicrosoftClientID:          viper.GetString("MS_CLIENT_ID"),
		MicrosoftClientSecret:      viper.GetString("MS_CLIENT_SECRET"),
		MicrosoftTenantID:          viper.GetString("MS_TENANT_ID"),
		RelaticsClientID:           viper.GetString("RELATICS_CLIENT_ID"),
		RelaticsClientSecret:       viper.GetString("RELATICS_CLIENT_SECRET"),
		RelaticsRedirectURI:        viper.GetString("RELATICS_REDIRECT_URI"),
		RelaticsRealm:              viper.GetString("RELATICS_REALM"),
		PostLoginRedirectURL:       viper.GetString("POST_LOGIN_REDIRECT_URL"),
		ErrorLoginRedirect:         viper.GetString("ERROR_LOGIN_REDIRECT"),
		OAuthCallbackURL:           viper.GetString("OAUTH_CALLBACK_URL"),
		AuthLoginURL:               viper.GetString("AUTH_LOGIN_URL"),
		CookieDomain:               viper.GetString("COOKIE_DOMAIN"),
		CookieSecure:               viper.GetBool("COOKIE_SECURE"),
		CookieSameSite:             viper.GetString("COOKIE_SAMESITE"),
		SessionTTLDays:             viper.GetInt("SESSION_TTL_DAYS"),
		SessionSecretKey:           viper.GetString("SESSION_SECRET_KEY"),
		SystemEmailsRaw:            viper.GetString("SYSTEM_EMAILS"),
		CORSOriginsRaw:             viper.GetString("CORS_ORIGINS"),
		MicrosoftEmailTenantID:     viper.GetString("MICROSOFT_EMAIL_TENANT_ID"),
		MicrosoftEmailClientID:     viper.GetString("MICROSOFT_EMAIL_CLIENT_ID"),
		MicrosoftEmailClientSecret: viper.GetString("MICROSOFT_EMAIL_CLIENT_SECRET"),
		MicrosoftEmailSender:       viper.GetString("MICROSOFT_EMAIL_SENDER"),
		GeoIPDBPath:                viper.GetString("GEOIP_DB_PATH"),
	}

	config.SetDefaults()
	settings = config
	return config, nil
}

func (c *Config) SetDefaults() {
	if c.AppName == "" {
		c.AppName = "vondr-identity"
	}
	if c.Environment == "" {
		c.Environment = "production"
	}
	if c.RootDomain == "" {
		c.RootDomain = "vondr.ai"
	}
	if c.CookieSameSite == "" {
		c.CookieSameSite = "lax"
	}
	if c.SessionTTLDays == 0 {
		c.SessionTTLDays = 7
	}
	if c.SessionSecretKey == "" {
		c.SessionSecretKey = "change-me-in-production"
	}
	if c.MicrosoftEmailSender == "" {
		c.MicrosoftEmailSender = "noreply@vondr.ai"
	}
	if c.RelaticsRealm == "" {
		c.RelaticsRealm = "cpmconsultancy"
	}
}

func (c *Config) SystemEmails() []string {
	if c.SystemEmailsRaw == "" {
		return []string{}
	}

	emails := strings.Split(c.SystemEmailsRaw, ",")
	result := make([]string, 0, len(emails))

	for _, email := range emails {
		trimmed := strings.TrimSpace(email)
		if trimmed != "" {
			result = append(result, strings.ToLower(trimmed))
		}
	}

	return result
}

func (c *Config) CORSOrigins() []string {
	if c.CORSOriginsRaw == "" {
		return []string{}
	}

	origins := strings.Split(c.CORSOriginsRaw, ",")
	result := make([]string, 0, len(origins))

	for _, origin := range origins {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

func GetConfig() *Config {
	return settings
}
