package main

import (
	_ "github.com/vondr/identity-go/docs"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/vondr/identity-go/internal/api/middleware"
	"github.com/vondr/identity-go/internal/api/public"
	"github.com/vondr/identity-go/internal/core"
	"github.com/vondr/identity-go/internal/core/oauth"
	"github.com/vondr/identity-go/internal/infrastructure/cache"
	"github.com/vondr/identity-go/internal/infrastructure/database"
	"github.com/vondr/identity-go/internal/infrastructure/geoip"
)

func main() {
	cfg, err := core.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := database.InitDB(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := cache.InitRedis(cfg.KeyDBURL); err != nil {
		log.Fatalf("Failed to initialize redis: %v", err)
	}

	if err := geoip.InitGeoIP(cfg.GeoIPDBPath); err != nil {
		log.Printf("Warning: Failed to initialize GeoIP: %v", err)
	}

	r := gin.Default()

	allowedOrigins := cfg.CORSOrigins()
	if len(allowedOrigins) > 0 {
		r.Use(middleware.CORSMiddleware(allowedOrigins))
	}

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := r.Group("/auth")
	{
		oauthConfig := &oauth.MicrosoftOAuthConfig{
			ClientID:     cfg.MicrosoftClientID,
			ClientSecret: cfg.MicrosoftClientSecret,
			TenantID:     cfg.MicrosoftTenantID,
			CallbackURL:  cfg.OAuthCallbackURL,
		}
		authHandler := public.NewAuthHandler(
			oauth.NewMicrosoftOAuthConfig(oauthConfig),
			cfg.OAuthCallbackURL,
			cfg.PostLoginRedirectURL,
			cfg.CookieDomain,
			cfg.CookieSecure,
			http.SameSiteLaxMode,
			cfg.SessionTTLDays,
			nil,
			nil,
			nil,
			nil,
			cfg.SystemEmails(),
			"",
			"",
		)
		auth.GET("/microsoft/login", authHandler.MicrosoftLogin)
		auth.GET("/microsoft/callback", authHandler.MicrosoftCallback)
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/me", authHandler.Me)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
