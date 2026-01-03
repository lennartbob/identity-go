package protected

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vondr/identity-go/internal/api/types"
)

type ForwardAuthHandler struct {
	sessionManager types.SessionManager
	memberService  types.MemberService
	appService     types.AppService
	orgService     types.OrganizationService
	countryService types.AppAllowedCountryService
	geoipService   types.GeoIPService
	authLoginURL   string
	errorLoginURL  string
}

func NewForwardAuthHandler(
	sessionManager types.SessionManager,
	memberService types.MemberService,
	appService types.AppService,
	orgService types.OrganizationService,
	countryService types.AppAllowedCountryService,
	geoipService types.GeoIPService,
	authLoginURL string,
	errorLoginURL string,
) *ForwardAuthHandler {
	return &ForwardAuthHandler{
		sessionManager: sessionManager,
		memberService:  memberService,
		appService:     appService,
		orgService:     orgService,
		countryService: countryService,
		geoipService:   geoipService,
		authLoginURL:   authLoginURL,
		errorLoginURL:  errorLoginURL,
	}
}

func extractClientIP(c *gin.Context) string {
	forwardedFor := c.GetHeader("x-forwarded-for")
	if forwardedFor != "" {
		parts := strings.Split(forwardedFor, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	realIP := c.GetHeader("x-real-ip")
	if realIP != "" {
		return strings.TrimSpace(realIP)
	}
	return c.ClientIP()
}

func checkCountryAccess(
	ctx context.Context,
	app *types.App,
	clientIP string,
	countryService types.AppAllowedCountryService,
	geoipService types.GeoIPService,
) string {
	allowedCodes, err := countryService.ListCountryCodes(ctx, app.ID)
	if err != nil || len(allowedCodes) == 0 {
		return ""
	}
	if !geoipService.IsEnabled() {
		return "GeoIP database not configured while country restrictions are enabled."
	}
	if clientIP == "" {
		return "Unable to determine client IP address for country validation."
	}

	ip := net.ParseIP(clientIP)
	if ip != nil && (ip.IsPrivate() || ip.IsLoopback()) {
		return ""
	}

	countryCode, err := geoipService.LookupCountry(clientIP)
	if err != nil || countryCode == "" {
		return "Could not resolve country for the provided IP address."
	}

	upperCode := strings.ToUpper(countryCode)
	for _, allowed := range allowedCodes {
		if strings.ToUpper(allowed) == upperCode {
			return ""
		}
	}
	return "Access from country '" + countryCode + "' is not allowed for this application."
}

func getForwardedValue(value string) string {
	if value == "" {
		return ""
	}
	parts := strings.Split(value, ",")
	if len(parts) > 0 {
		return strings.TrimSpace(parts[0])
	}
	return ""
}

func buildOriginalRequestURL(c *gin.Context) string {
	forwardedProto := getForwardedValue(c.GetHeader("x-forwarded-proto"))
	forwardedHost := getForwardedValue(c.GetHeader("x-forwarded-host"))
	forwardedURI := getForwardedValue(c.GetHeader("x-forwarded-uri"))

	if forwardedProto == "" || forwardedHost == "" {
		return ""
	}

	proto := strings.ToLower(forwardedProto)
	if proto != "http" && proto != "https" {
		return ""
	}

	path := forwardedURI
	if path == "" {
		path = c.GetHeader("x-original-uri")
		if path == "" {
			path = c.Request.URL.Path
		}
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return proto + "://" + forwardedHost + path
}

func buildLoginRedirectURL(c *gin.Context, authLoginURL string) string {
	originalURL := buildOriginalRequestURL(c)
	if originalURL == "" {
		return authLoginURL
	}
	return authLoginURL + "?auto=1&return_to=" + url.QueryEscape(originalURL)
}

func buildErrorRedirectURL(errorLoginURL string, errorCode string) string {
	return errorLoginURL + "/" + errorCode
}

func (h *ForwardAuthHandler) Verify(c *gin.Context) {
	ctx := c.Request.Context()

	originalMethod := c.GetHeader("x-forwarded-method")
	if originalMethod == "" {
		originalMethod = c.Request.Method
	}

	isPreflight := strings.ToUpper(originalMethod) == "OPTIONS" || c.GetHeader("access-control-request-method") != ""
	sessionToken, _ := c.Cookie("session_token")

	if isPreflight && sessionToken == "" {
		c.Status(http.StatusOK)
		return
	}

	acceptHeader := c.GetHeader("accept")
	isBrowserRequest := strings.Contains(acceptHeader, "text/html")

	m2mToken := c.GetHeader("x-vondr-auth")
	if m2mToken != "" {
		h.handleM2MAuth(ctx, c, m2mToken, isBrowserRequest)
		return
	}

	if sessionToken == "" {
		h.handleNoSession(c, isBrowserRequest)
		return
	}

	sessionData, err := h.sessionManager.GetSession(ctx, sessionToken)
	if err != nil {
		h.handleInvalidSession(c, isBrowserRequest)
		return
	}

	member, err := h.memberService.GetByID(ctx, sessionData.MemberID)
	if err != nil {
		h.handleInvalidSession(c, isBrowserRequest)
		return
	}

	forwardedHost := c.GetHeader("x-forwarded-host")
	forwardedProto := c.GetHeader("x-forwarded-proto")

	if member.Role == "system" {
		c.Header("x-vondr-user-id", member.ID)
		c.Header("x-vondr-email", member.Email)
		c.Header("x-vondr-organization-id", member.OrganizationID)
		c.Status(http.StatusOK)
		return
	}

	if forwardedProto == "" || strings.ToLower(forwardedProto) != "https" {
		h.handleHTTPSRequired(c, isBrowserRequest)
		return
	}

	if forwardedHost != "" {
		org, err := h.orgService.GetByID(ctx, member.OrganizationID)
		if err != nil {
			h.handleInvalidSession(c, isBrowserRequest)
			return
		}

		allowedDomains, _ := h.appService.GetAllowedDomainsForOrg(ctx, member.OrganizationID, org.Hostname)

		allowed := false
		for _, domain := range allowedDomains {
			if domain == forwardedHost {
				allowed = true
				break
			}
		}

		if !allowed {
			h.handleAppNotAllowed(c, isBrowserRequest)
			return
		}

		domainMap, _ := h.appService.GetDomainAppMap(ctx, member.OrganizationID, org.Hostname)
		if domainMap != nil {
			targetApp := domainMap[forwardedHost]
			if targetApp != nil {
				clientIP := extractClientIP(c)
				countryError := checkCountryAccess(ctx, targetApp, clientIP, h.countryService, h.geoipService)
				if countryError != "" {
					h.handleCountryBlocked(c, isBrowserRequest)
					return
				}
			}
		}
	}

	c.Header("x-vondr-user-id", member.ID)
	c.Header("x-vondr-email", member.Email)
	c.Header("x-vondr-organization-id", member.OrganizationID)
	c.Status(http.StatusOK)
}

func (h *ForwardAuthHandler) handleM2MAuth(ctx context.Context, c *gin.Context, m2mToken string, isBrowserRequest bool) {
	app, err := h.appService.GetByToken(ctx, m2mToken)
	if err != nil {
		h.handleUnauthorized(c, isBrowserRequest, "Invalid authentication token")
		return
	}

	rawUserID := c.GetHeader("x-vondr-user-id")
	if rawUserID == "" {
		h.handleUnauthorized(c, isBrowserRequest, "x-vondr-user-id header is required when using API token authentication")
		return
	}

	member, err := h.memberService.GetByID(ctx, rawUserID)
	if err != nil {
		h.handleUnauthorized(c, isBrowserRequest, "Member not found for provided x-vondr-user-id")
		return
	}

	if member.OrganizationID != app.OrganizationID {
		h.handleUnauthorized(c, isBrowserRequest, "Member is not allowed to use this application token")
		return
	}

	forwardedHost := c.GetHeader("x-forwarded-host")
	if forwardedHost != "" {
		org, err := h.orgService.GetByID(ctx, app.OrganizationID)
		if err != nil {
			h.handleUnauthorized(c, isBrowserRequest, "Organization associated with this application no longer exists")
			return
		}

		allowedDomains, _ := h.appService.GetAllowedDomainsForOrg(ctx, app.OrganizationID, org.Hostname)

		allowed := false
		for _, domain := range allowedDomains {
			if domain == forwardedHost {
				allowed = true
				break
			}
		}

		if !allowed {
			h.handleUnauthorized(c, isBrowserRequest, "Access to this domain is not allowed for this application token")
			return
		}

		clientIP := extractClientIP(c)
		countryError := checkCountryAccess(ctx, app, clientIP, h.countryService, h.geoipService)
		if countryError != "" {
			h.handleUnauthorized(c, isBrowserRequest, countryError)
			return
		}
	}

	c.Header("x-vondr-user-id", member.ID)
	c.Header("x-vondr-email", member.Email)
	c.Header("x-vondr-organization-id", member.OrganizationID)
	c.Status(http.StatusOK)
}

func (h *ForwardAuthHandler) handleNoSession(c *gin.Context, isBrowserRequest bool) {
	if isBrowserRequest {
		loginURL := buildLoginRedirectURL(c, h.authLoginURL)
		c.Redirect(http.StatusFound, loginURL)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No session cookie found"})
	}
}

func (h *ForwardAuthHandler) handleInvalidSession(c *gin.Context, isBrowserRequest bool) {
	if isBrowserRequest {
		loginURL := buildLoginRedirectURL(c, h.authLoginURL)
		c.Redirect(http.StatusFound, loginURL)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
	}
}

func (h *ForwardAuthHandler) handleHTTPSRequired(c *gin.Context, isBrowserRequest bool) {
	if isBrowserRequest {
		errorURL := buildErrorRedirectURL(h.errorLoginURL, "403/app_not_allowed")
		c.Redirect(http.StatusFound, errorURL)
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only HTTPS requests are allowed"})
	}
}

func (h *ForwardAuthHandler) handleAppNotAllowed(c *gin.Context, isBrowserRequest bool) {
	if isBrowserRequest {
		errorURL := buildErrorRedirectURL(h.errorLoginURL, "403/app_not_allowed")
		c.Redirect(http.StatusFound, errorURL)
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access to this application is not allowed for your organization"})
	}
}

func (h *ForwardAuthHandler) handleCountryBlocked(c *gin.Context, isBrowserRequest bool) {
	if isBrowserRequest {
		errorURL := buildErrorRedirectURL(h.errorLoginURL, "403/app_country_blocked")
		c.Redirect(http.StatusFound, errorURL)
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access from this country is not allowed for this application"})
	}
}

func (h *ForwardAuthHandler) handleUnauthorized(c *gin.Context, isBrowserRequest bool, message string) {
	if isBrowserRequest {
		loginURL := buildLoginRedirectURL(c, h.authLoginURL)
		c.Redirect(http.StatusFound, loginURL)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": message})
	}
}
