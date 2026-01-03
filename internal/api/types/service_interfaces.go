package types

import (
	"context"
	"time"

	"github.com/vondr/identity-go/internal/core"
)

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
	GetByID(ctx context.Context, memberID string) (*Member, error)
	LinkMicrosoftAccount(ctx context.Context, email, microsoftID, firstName, lastName string) (*Member, error)
	CreateSystemMember(ctx context.Context, orgID, orgName, email, microsoftID, firstName, lastName string) (*Member, error)
}

type OrganizationService interface {
	GetByID(ctx context.Context, orgID string) (*Organization, error)
}

type AppService interface {
	GetByToken(ctx context.Context, token string) (*App, error)
	GetByID(ctx context.Context, appID string) (*App, error)
	GetAllowedDomainsForOrg(ctx context.Context, orgID, hostname string) ([]string, error)
	GetDomainAppMap(ctx context.Context, orgID, hostname string) (map[string]*App, error)
}

type AppAllowedCountryService interface {
	ListCountryCodes(ctx context.Context, appID string) ([]string, error)
}

type GeoIPService interface {
	LookupCountry(ip string) (string, error)
	IsEnabled() bool
}

type SessionService interface {
	RecordLogin(ctx context.Context, microsoftID, email, orgID string) error
	GetLastLoginForMember(ctx context.Context, memberID string) (*time.Time, error)
	GetLastLoginsBatch(ctx context.Context, memberIDs []string) (map[string]*time.Time, error)
}

type Member struct {
	ID             string
	Email          string
	FirstName      *string
	LastName       *string
	MicrosoftID    *string
	OrganizationID string
	Role           core.MemberRole
}

type Organization struct {
	ID       string
	Name     string
	Hostname string
}

type App struct {
	ID              string
	OrganizationID  string
	Name            string
	SubdomainLabels []string
	MainLabel       string
	IsPlatformApp   bool
}
