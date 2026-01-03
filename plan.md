# Identity Service - Go Rewrite Plan

## Overview
This document outlines the plan to rewrite the Vondr Identity Management API from Python/FastAPI to Go.

## Technology Stack Translation

| Python Component | Go Equivalent |
|------------------|---------------|
| FastAPI | Gin Web Framework |
| SQLAlchemy ORM | GORM |
| PostgreSQL Driver | pgx |
| Redis/KeyDB | go-redis |
| authlib (OAuth) | golang.org/x/oauth2 |
| GeoIP2 | oschwald/geoip2-golang |
| pydantic-settings | Viper |
| httpx | net/http + http.Client |
| Microsoft Graph API | go-microsoft-graph-client |

## Project Structure

```
identity-go/
├── cmd/
│   ├── public/
│   │   └── main.go          # Public API entrypoint
│   └── protected/
│       └── main.go          # Protected API entrypoint
├── internal/
│   ├── api/
│   │   ├── public/           # Public auth endpoints
│   │   │   ├── handlers.go
│   │   │   └── router.go
│   │   ├── protected/        # Protected API endpoints
│   │   │   ├── organizations/
│   │   │   │   └── handlers.go
│   │   │   ├── members/
│   │   │   │   └── handlers.go
│   │   │   ├── apps/
│   │   │   │   └── handlers.go
│   │   │   ├── groups/
│   │   │   │   └── handlers.go
│   │   │   ├── invitations/
│   │   │   │   └── handlers.go
│   │   │   └── router.go
│   │   └── middleware/
│   │       ├── admin_auth.go
│   │       └── cors.go
│   ├── application/
│   │   ├── services/
│   │   │   ├── member_service.go
│   │   │   ├── organization_service.go
│   │   │   ├── app_service.go
│   │   │   ├── session_service.go
│   │   │   ├── session_manager.go
│   │   │   ├── user_group_service.go
│   │   │   ├── invitation_service.go
│   │   │   └── app_allowed_country_service.go
│   │   └── repositories/
│   │       ├── member_repository.go
│   │       ├── organization_repository.go
│   │       ├── app_repository.go
│   │       ├── session_repository.go
│   │       ├── user_group_repository.go
│   │       ├── user_group_member_repository.go
│   │       └── invitation_repository.go
│   ├── core/
│   │   ├── config.go
│   │   ├── oauth/
│   │   │   ├── microsoft.go
│   │   │   └── relatics.go
│   │   ├── roles.go
│   │   ├── constants.go
│   │   └── errors.go
│   ├── infrastructure/
│   │   ├── database/
│   │   │   ├── postgres.go
│   │   │   └── models/
│   │   │       ├── organization.go
│   │   │       ├── member.go
│   │   │       ├── app.go
│   │   │       ├── session.go
│   │   │       ├── user_group.go
│   │   │       ├── user_group_member.go
│   │   │       ├── app_allowed_country.go
│   │   │       └── invitation.go
│   │   ├── cache/
│   │   │   ├── redis.go
│   │   │   └── cached_repositories.go
│   │   └── geoip/
│   │       └── geoip.go
│   └── services/
│       ├── email.go          # Microsoft Graph email
│       └── utils.go
├── pkg/
│   └── domain/
│       ├── entities.go       # Domain entities
│       └── types.go
├── scripts/
│   ├── migrations/
│   └── init_database.go
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── .env.example
```

## Implementation Tasks

### 1. Project Setup & Dependencies
- [ ] Initialize Go module: `go mod init github.com/vondr/identity-go`
- [ ] Create directory structure
- [ ] Add dependencies:
  - `github.com/gin-gonic/gin` - Web framework
  - `gorm.io/gorm` - ORM
  - `gorm.io/driver/postgres` - PostgreSQL driver
  - `github.com/redis/go-redis/v9` - Redis client
  - `golang.org/x/oauth2` - OAuth2 client
  - `github.com/oschwald/geoip2-golang` - GeoIP
  - `github.com/spf13/viper` - Configuration management
  - `github.com/golang-jwt/jwt/v5` - JWT tokens
  - `github.com/google/uuid` - UUID generation
  - `go.uber.org/zap` - Logging

### 2. Configuration Management
- [ ] Create `config.go` with Viper integration
- [ ] Support environment variables and .env file
- [ ] Configuration struct with all settings:
  - Database URL
  - KeyDB/Redis URL
  - Microsoft OAuth (client_id, client_secret, tenant_id)
  - Relatics OAuth
  - Session settings (TTL, secret key)
  - Cookie settings (domain, secure, samesite)
  - CORS origins
  - Admin token
  - System emails (comma-separated)
  - Redirect URLs
  - GeoIP database path

### 3. Database Models (GORM)
- [ ] Organization model
- [ ] OrganizationMember model
- [ ] App model
- [ ] Session model
- [ ] UserGroup model
- [ ] UserGroupMember model
- [ ] AppAllowedCountry model
- [ ] Invitation model
- [ ] Database connection setup with pgx driver
- [ ] Auto-migration support

### 4. Repository Layer
- [ ] MemberRepository interface and implementation
- [ ] OrganizationRepository interface and implementation
- [ ] AppRepository interface and implementation
- [ ] SessionRepository (Redis-backed)
- [ ] UserGroupRepository
- [ ] UserGroupMemberRepository
- [ ] InvitationRepository
- [ ] Cached repository wrappers (optional, layer over base repos)

### 5. Cache Layer
- [ ] Redis connection setup
- [ ] Session storage implementation
- [ ] Cache decorators for repositories (if needed)
- [ ] SessionManager with TTL support

### 6. Core Services (Application Logic)
- [ ] MemberService
  - List, get by ID, get by email
  - Get by Microsoft ID
  - Link Microsoft account
  - Create system (super admin)
  - Create member
  - Update member
  - Delete (soft delete) member
- [ ] OrganizationService
  - List, get, create, update, delete
  - Get by hostname
  - Check unique constraints
- [ ] AppService
  - CRUD operations
  - Get by token
  - Get allowed domains for organization
  - Get domain-app mapping
- [ ] SessionService
  - Record login
  - Get last login for member
  - Get last logins batch (N+1 optimization)
- [ ] SessionManager
  - Create session
  - Get session
  - Delete session
- [ ] UserGroupService
  - CRUD operations
  - List groups for member
  - Add member to group
  - Remove member from group
- [ ] AppAllowedCountryService
  - CRUD operations
  - List country codes for app
- [ ] InvitationService
  - Create invitation
  - Validate invitation token
  - List pending invitations

### 7. OAuth Integration
- [ ] Microsoft OAuth config
- [ ] Relatics OAuth config
- [ ] Authorization URL builder
- [ ] Token exchange handler
- [ ] User info retrieval
- [ ] State parameter management (for OAuth flow)

### 8. GeoIP Service
- [ ] MaxMind GeoLite2 database reader
- [ ] Country lookup by IP
- [ ] Private IP detection (allow pass-through)
- [ ] Service enabled/disabled flag

### 9. Email Service
- [ ] Microsoft Graph API integration
- [ ] Send invitation emails
- [ ] Error handling

### 10. API Middleware
- [ ] AdminTokenMiddleware (validate x-vondr-admin-token header)
- [ ] CORSMiddleware
- [ ] Context extraction middleware (from Traefik headers)
- [ ] Error handling middleware

### 11. Public API Handlers
Routes:
- `GET /auth/microsoft/login` - Redirect to Microsoft OAuth
- `GET /auth/relatics/login` - Redirect to Relatics OAuth (test endpoint)
- `GET /auth/microsoft/callback` - Handle Microsoft OAuth callback
- `GET /auth/relatics/callback` - Handle Relatics OAuth callback (test endpoint)
- `POST /auth/logout` - Logout and clear session
- `GET /healthz` - Health check

Implementation details:
- OAuth flow with state parameter (return_to URL)
- Create/link member on callback
- Create session in Redis
- Set HTTP-only cookie
- Redirect to frontend

### 12. Protected API Handlers
Routes (under `/api/v1`):
- `GET /auth/verify` - Traefik forward auth endpoint
- `GET /auth/me` - Get current user

Organizations:
- `GET /api/v1/organizations` - List accessible organizations
- `GET /api/v1/organizations/{id}` - Get organization details (with members)
- `POST /api/v1/organizations` - Create organization (system only)
- `PUT /api/v1/organizations/{id}` - Update organization (system only)
- `DELETE /api/v1/organizations/{id}` - Delete organization (system only)

Members:
- `GET /api/v1/organizations/{org_id}/members` - List members
- `POST /api/v1/organizations/{org_id}/members` - Create member
- `GET /api/v1/organizations/{org_id}/members/{id}` - Get member
- `PUT /api/v1/organizations/{org_id}/members/{id}` - Update member
- `DELETE /api/v1/organizations/{org_id}/members/{id}` - Delete member

Apps:
- `GET /api/v1/organizations/{org_id}/apps` - List apps
- `POST /api/v1/organizations/{org_id}/apps` - Create app
- `GET /api/v1/organizations/{org_id}/apps/{id}` - Get app
- `PUT /api/v1/organizations/{org_id}/apps/{id}` - Update app
- `DELETE /api/v1/organizations/{org_id}/apps/{id}` - Delete app

Groups:
- `GET /api/v1/organizations/{org_id}/groups` - List groups
- `POST /api/v1/organizations/{org_id}/groups` - Create group
- `GET /api/v1/organizations/{org_id}/groups/{id}` - Get group
- `PUT /api/v1/organizations/{org_id}/groups/{id}` - Update group
- `DELETE /api/v1/organizations/{org_id}/groups/{id}` - Delete group
- `POST /api/v1/organizations/{org_id}/groups/{id}/members` - Add member
- `DELETE /api/v1/organizations/{org_id}/groups/{id}/members/{member_id}` - Remove member

Invitations:
- `POST /api/v1/organizations/{org_id}/invitations` - Send invitation
- `GET /api/v1/organizations/{org_id}/invitations` - List invitations
- `GET /api/v1/organizations/{org_id}/invitations/{token}` - Get invitation
- `DELETE /api/v1/organizations/{org_id}/invitations/{token}` - Cancel invitation

### 13. Forward Auth Implementation
The `/auth/verify` endpoint must:
1. Extract session token from cookie or M2M token from header
2. For M2M token (`x-vondr-auth`):
   - Validate token exists in App table
   - Verify `x-vondr-user-id` header is present and valid
   - Verify member exists and belongs to same org as app
   - Verify requested host is allowed for this app
   - Enforce country restrictions (if any)
   - Return 200 with user headers
3. For session token:
   - Validate session in Redis
   - Get member from database
   - Verify member exists and is active
   - Check domain access for organization
   - Enforce country restrictions for target app
   - Return 200 with user headers
4. For invalid/no auth:
   - If browser request (Accept: text/html): Return 302 to login
   - If API request: Return 401/403

### 14. Docker Configuration
- [ ] Dockerfile (multi-stage build)
- [ ] docker-compose.yml (Postgres, KeyDB, public API, protected API)
- [ ] Download GeoIP database in Dockerfile

### 15. Main Entry Points
**Public API** (`cmd/public/main.go`):
- Initialize dependencies
- Setup middleware
- Register public auth routes
- Run on port 8000

**Protected API** (`cmd/protected/main.go`):
- Initialize dependencies
- Setup middleware (AdminTokenMiddleware)
- Register forward auth endpoint
- Register protected API routes (under /api/v1)
- Run on port 8000

### 16. Testing & Validation
- [ ] Unit tests for services
- [ ] Integration tests for API endpoints
- [ ] Test OAuth flow
- [ ] Test forward auth
- [ ] Test country restrictions
- [ ] Load testing

## Key Features to Preserve

1. **Multi-tenancy**: Organizations with isolated data
2. **Microsoft OAuth**: SSO integration
3. **Traefik Forward Auth**: Middleware for protected services
4. **Session Management**: Redis-backed with TTL
5. **Country-based Access Control**: GeoIP-based restrictions per app
6. **User Groups**: Role-based group memberships
7. **Invitations**: Email-based member invitations
8. **M2M Auth**: API token-based authentication
9. **Admin Token**: Super admin access
10. **CORS Support**: Configurable origins

## Environment Variables

Required variables (copy from Python version):
```
# Application
ENVIRONMENT=production
ROOT_DOMAIN=vondr.ai
ADMIN_TOKEN=your-admin-token

# Database
DATABASE_URL=postgresql+psycopg://identity:identity@db:5432/identity
KEYDB_URL=redis://keydb:6379/0

# Microsoft OAuth
MS_CLIENT_ID=your-client-id
MS_CLIENT_SECRET=your-client-secret
MS_TENANT_ID=your-tenant-id

# Relatics OAuth (optional)
RELATICS_CLIENT_ID=
RELATICS_CLIENT_SECRET=
RELATICS_REDIRECT_URI=
RELATICS_REALM=cpmconsultancy

# Auth & Session
POST_LOGIN_REDIRECT_URL=https://your-frontend.com
ERROR_LOGIN_REDIRECT=https://your-frontend.com/error
OAUTH_CALLBACK_URL=https://auth.vondr.ai/auth/microsoft/callback
AUTH_LOGIN_URL=https://auth.vondr.ai
COOKIE_DOMAIN=vondr.ai
COOKIE_SECURE=true
COOKIE_SAMESITE=lax
SESSION_TTL_DAYS=7
SESSION_SECRET_KEY=your-long-random-secret-key

# System
SYSTEM_EMAILS=admin@vondr.ai,superadmin@vondr.ai

# Email (Microsoft Graph)
MICROSOFT_EMAIL_TENANT_ID=
MICROSOFT_EMAIL_CLIENT_ID=
MICROSOFT_EMAIL_CLIENT_SECRET=
MICROSOFT_EMAIL_SENDER=noreply@vondr.ai

# CORS
CORS_ORIGINS=https://your-frontend.com,https://another.com

# GeoIP
GEOIP_DB_PATH=/app/geoip/GeoLite2-Country.mmdb
```

## Migration Notes

1. **Database Schema**: Use existing schema from Python version
2. **Redis Keys**: Match existing session key format
3. **Cookie Names**: Use same cookie names for compatibility
4. **Headers**: Keep same Traefik forward auth headers
5. **API Routes**: Maintain exact route paths for compatibility

## Performance Considerations

1. **Database Connection Pooling**: GORM handles this automatically
2. **Redis Connection Pooling**: Use go-redis/v9 pool
3. **N+1 Queries**: Implement batch loading for related data
4. **Caching**: Use Redis for frequently accessed data
5. **JSON**: Use json.Marshal/Unmarshal for API responses

## Next Steps

Once Go is installed:
1. Run `go mod init github.com/vondr/identity-go`
2. Create directory structure
3. Implement tasks in order from 1-16
4. Test with existing infrastructure (Postgres, KeyDB)
5. Deploy to replace Python version
