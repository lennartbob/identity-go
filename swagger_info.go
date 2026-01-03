package docs

// @title           Vondr Identity API
// @version         1.0
// @description     Vondr Identity Management API - Go implementation
// @termsOfService  https://vondr.ai/terms

// @contact.name   Vondr Support
// @contact.url    https://vondr.ai/support
// @contact.email  support@vondr.ai

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8089
// @BasePath  /

// @securityDefinitions.apikey AdminToken
// @in header
// @name x-vondr-admin-token
// @description Admin authentication token

// @securityDefinitions.apikey SessionToken
// @in cookie
// @name session_token
// @description User session token (HTTP-only cookie)
