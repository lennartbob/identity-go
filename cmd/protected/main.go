package main

import (
	_ "github.com/vondr/identity-go/docs"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/vondr/identity-go/internal/api/middleware"
	"github.com/vondr/identity-go/internal/api/protected"
	"github.com/vondr/identity-go/internal/core"
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

	forwardAuthHandler := protected.NewForwardAuthHandler(
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		cfg.AuthLoginURL,
		cfg.ErrorLoginRedirect,
	)

	r.Any("/auth/verify", forwardAuthHandler.Verify)

	r.Use(middleware.AdminAuthMiddleware(cfg.AdminToken))
	r.Use(middleware.ExtractVondrContext())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
