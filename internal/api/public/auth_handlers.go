package public

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vondr/identity-go/internal/core/oauth"
)

type AuthHandler struct {
	oauthClient    *oauth.MicrosoftOAuthClient
	callbackURL    string
	postLoginURL   string
	cookieDomain   string
	cookieSecure   bool
	cookieSameSite http.SameSite
	sessionTTL     int
	sessionManager SessionManager
	memberService  MemberService
	orgService     OrganizationService
	sessionService SessionService
	systemEmails   map[string]bool
	defaultOrgID   string
	defaultOrgName string
}

type SessionManager interface {
	CreateSession(ctx context.Context, memberID, email, orgID, microsoftID string) (string, error)
	GetSession(ctx context.Context, token string) (*SessionData, error)
	DeleteSession(ctx context.Context, token string) error
}

type SessionData struct {
	MemberID       string
	Email          string
	OrganizationID string
	MicrosoftID    string
}

type MemberService interface {
	GetByMicrosoftID(ctx context.Context, microsoftID string) (*Member, error)
	GetByEmail(ctx context.Context, email string) (*Member, error)
	LinkMicrosoftAccount(ctx context.Context, email, microsoftID, firstName, lastName string) (*Member, error)
	CreateSystemMember(ctx context.Context, orgID, orgName, email, microsoftID, firstName, lastName string) (*Member, error)
	GetByID(ctx context.Context, memberID string) (*Member, error)
}

type OrganizationService interface {
	GetByID(ctx context.Context, orgID string) (*Organization, error)
}

type SessionService interface {
	RecordLogin(ctx context.Context, microsoftID, email, orgID string) error
}

type Member struct {
	ID             string
	Email          string
	FirstName      *string
	LastName       *string
	MicrosoftID    *string
	OrganizationID string
	Role           string
}

type Organization struct {
	ID       string
	Name     string
	Hostname *string
}

func NewAuthHandler(
	oauthClient *oauth.MicrosoftOAuthClient,
	callbackURL string,
	postLoginURL string,
	cookieDomain string,
	cookieSecure bool,
	cookieSameSite http.SameSite,
	sessionTTL int,
	sessionManager SessionManager,
	memberService MemberService,
	orgService OrganizationService,
	sessionService SessionService,
	systemEmails []string,
	defaultOrgID string,
	defaultOrgName string,
) *AuthHandler {
	systemEmailsMap := make(map[string]bool)
	for _, email := range systemEmails {
		systemEmailsMap[strings.ToLower(email)] = true
	}
	return &AuthHandler{
		oauthClient:    oauthClient,
		callbackURL:    callbackURL,
		postLoginURL:   postLoginURL,
		cookieDomain:   cookieDomain,
		cookieSecure:   cookieSecure,
		cookieSameSite: cookieSameSite,
		sessionTTL:     sessionTTL,
	}
}

func generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func (h *AuthHandler) MicrosoftLogin(c *gin.Context) {
	returnTo := c.Query("return_to")

	state, err := generateState()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate state"})
		return
	}

	authURL := h.oauthClient.GetAuthURL(h.callbackURL, state, returnTo)
	c.Redirect(http.StatusFound, authURL)
}

func (h *AuthHandler) MicrosoftCallback(c *gin.Context) {
	ctx := c.Request.Context()

	code := c.Query("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing authorization code"})
		return
	}

	userInfo, err := h.oauthClient.ExchangeCode(ctx, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange authorization code: " + err.Error()})
		return
	}

	email := strings.ToLower(userInfo.Email)
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not provided by Microsoft"})
		return
	}

	microsoftID := userInfo.ID

	c.JSON(http.StatusOK, gin.H{"message": "OAuth callback received", "email": email, "microsoft_id": microsoftID})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	ctx := c.Request.Context()

	sessionToken, err := c.Cookie("session_token")
	if err == nil && sessionToken != "" {
		_ = h.sessionManager.DeleteSession(ctx, sessionToken)
	}

	c.SetCookie(
		"session_token",
		"",
		-1,
		"/",
		h.cookieDomain,
		h.cookieSecure,
		h.cookieSameSite == http.SameSiteNoneMode,
	)

	c.JSON(http.StatusOK, gin.H{"status": "logged_out"})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": userID,
	})
}
