package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminAuthMiddleware(adminToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("x-vondr-admin-token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "admin token required"})
			return
		}

		if token != adminToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid admin token"})
			return
		}

		c.Next()
	}
}

type VondrContext struct {
	UserID         string
	Email          string
	OrganizationID string
	IsSystem       bool
}

func ExtractVondrContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("x-vondr-user-id")
		email := c.GetHeader("x-vondr-email")
		orgID := c.GetHeader("x-vondr-organization-id")

		ctx := &VondrContext{
			UserID:         userID,
			Email:          email,
			OrganizationID: orgID,
			IsSystem:       false,
		}

		c.Set("vondr_ctx", ctx)
		c.Next()
	}
}

func GetVondrContext(c *gin.Context) *VondrContext {
	ctx, _ := c.Get("vondr_ctx")
	return ctx.(*VondrContext)
}

func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowed := false

		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, x-vondr-admin-token, x-vondr-user-id, x-vondr-email, x-vondr-organization-id")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
