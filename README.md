# Vondr Identity Service - Go Implementation

This is a Go rewrite of the Vondr Identity Management API, originally written in Python/FastAPI.

## Project Status

The following components have been implemented:

- ✅ Project structure and Go modules
- ✅ Docker configuration (docker-compose.yml, Dockerfile)
- ✅ Configuration management with Viper
- ✅ Database models (Organization, Member, App, UserGroup, Session, Invitation)
- ✅ Repository layer (GORM-based PostgreSQL repositories)
- ✅ Cache layer (Redis/KeyDB for session storage)
- ✅ Application services (Member, Organization, App, Session services)
- ✅ Application services (UserGroup, AppAllowedCountry services)
- ✅ OAuth integration (Microsoft and Relatics)
- ✅ GeoIP service for country-based access control
- ✅ Session management
- ✅ Admin token middleware
- ✅ Main entry points (public and protected APIs)
- ✅ Public API handlers (Microsoft login, callback, logout)
- ✅ Forward auth endpoint (/auth/verify) for Traefik
- ✅ Route registration in main.go files

## Technology Stack

| Component | Technology |
|-----------|-----------|
| Web Framework | Gin |
| ORM | GORM |
| Database | PostgreSQL (via pgx driver) |
| Cache | Redis/KeyDB |
| OAuth | golang.org/x/oauth2 |
| GeoIP | oschwald/geoip2-golang |
| Config | Viper |

## Project Structure

```
identity-go/
├── cmd/
│   ├── public/main.go          # Public API entrypoint (auth only)
│   └── protected/main.go       # Protected API entrypoint (with forward auth)
├── internal/
│   ├── api/
│   │   └── middleware/         # Auth and CORS middleware
│   ├── application/
│   │   ├── repositories/        # Data access layer
│   │   └── services/           # Business logic
│   ├── core/
│   │   ├── oauth/              # OAuth configuration
│   │   ├── config.go           # Settings management
│   │   ├── roles.go            # Member roles
│   │   └── errors.go           # Common errors
│   └── infrastructure/
│       ├── database/
│       │   └── models/          # GORM models
│       ├── cache/               # Redis client
│       └── geoip/              # GeoIP service
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── .env.example
```

## Getting Started

### Prerequisites

- Go 1.25+
- Docker and Docker Compose
- PostgreSQL 16
- Redis/KeyDB

### Environment Configuration

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

Required environment variables:
- `DATABASE_URL` - PostgreSQL connection string
- `KEYDB_URL` - Redis connection string
- `ADMIN_TOKEN` - Admin authentication token
- `MS_CLIENT_ID` - Microsoft OAuth client ID
- `MS_CLIENT_SECRET` - Microsoft OAuth client secret
- `MS_TENANT_ID` - Microsoft Azure tenant ID

### Running with Docker

```bash
docker-compose up -d
```

Services:
- `db` - PostgreSQL database
- `keydb` - Redis cache
- `identity-api-public` - Public authentication API (port 8088)
- `identity-api` - Protected API (port 8089)

### Local Development

```bash
# Install dependencies
go mod download

# Run public API
go run cmd/public/main.go

# Run protected API
go run cmd/protected/main.go
```

### Building

```bash
# Build public API
go build -o bin/public-api ./cmd/public

# Build protected API
go build -o bin/protected-api ./cmd/protected
```

## API Endpoints

### Public API (Authentication Only)

- `GET /healthz` - Health check
- `GET /auth/microsoft/login` - Microsoft OAuth login
- `GET /auth/microsoft/callback` - Microsoft OAuth callback
- `POST /auth/logout` - Logout

### Protected API (Authenticated Endpoints)

- `GET /healthz` - Health check
- `GET /auth/verify` - Traefik forward auth endpoint
- `GET /api/v1/auth/me` - Get current user

#### Organizations (under `/api/v1/organizations`)

- `GET /` - List organizations
- `POST /` - Create organization
- `GET /{id}` - Get organization details
- `PUT /{id}` - Update organization
- `DELETE /{id}` - Delete organization

#### Members (under `/api/v1/organizations/{org_id}/members`)

- `GET /` - List members
- `POST /` - Create member
- `GET /{id}` - Get member
- `PUT /{id}` - Update member
- `DELETE /{id}` - Delete member

#### Apps (under `/api/v1/organizations/{org_id}/apps`)

- `GET /` - List apps
- `POST /` - Create app
- `GET /{id}` - Get app
- `PUT /{id}` - Update app
- `DELETE /{id}` - Delete app

#### Groups (under `/api/v1/organizations/{org_id}/groups`)

- `GET /` - List groups
- `POST /` - Create group
- `GET /{id}` - Get group
- `PUT /{id}` - Update group
- `DELETE /{id}` - Delete group

## Architecture

### Layered Design

1. **API Layer** (`internal/api/`) - HTTP handlers and middleware
2. **Service Layer** (`internal/application/services/`) - Business logic
3. **Repository Layer** (`internal/application/repositories/`) - Data access
4. **Infrastructure Layer** (`internal/infrastructure/`) - External services

### Multi-Tenancy

- Organizations are isolated by `organization_id`
- Members belong to exactly one organization
- Apps are scoped to organizations

### Authentication

1. **Session-based auth** - Browser requests use HTTP-only cookies
2. **M2M auth** - API requests use `x-vondr-auth` token header
3. **Forward auth** - Traefik calls `/auth/verify` for protected routes

### Session Management

- Sessions stored in Redis
- TTL configurable (default: 7 days)
- Includes member_id, email, organization_id, microsoft_id

## Remaining Work

The following components from the plan are **not yet implemented**:

- [ ] Public API handlers (login, callback, logout) - scaffolding only
- [ ] Protected API handlers (organizations, members, apps, groups, invitations)
- [ ] Forward auth endpoint (`/auth/verify`) for Traefik
- [ ] Email service for invitations (Microsoft Graph API)

These need to be implemented to match the full functionality of the Python version.

## Differences from Python Version

1. **Framework**: FastAPI → Gin
2. **ORM**: SQLAlchemy → GORM
3. **Config**: pydantic-settings → Viper
4. **Session**: Starlette middleware → Redis-based
5. **Structure**: Similar layered architecture maintained

## Testing

```bash
# Run tests (when implemented)
go test ./...

# Run with coverage
go test -cover ./...
```

## License

Same as the original Python implementation.
