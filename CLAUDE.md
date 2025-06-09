# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Architecture Overview

Skyclerk is a SaaS bookkeeping application with:
- **Backend**: Go (Gin framework) serving RESTful API at `/api/v3`
- **Frontend**: Angular 7.2 app in `frontend/` 
- **Admin**: Angular 8.1 app in `centcom/`
- **Database**: MySQL/MariaDB with GORM ORM
- **Cache**: Redis
- **File Storage**: AWS S3/CloudFront
- **Payments**: Stripe integration

## Development Commands

### Backend
```bash
# Start local services (MySQL, Redis, Mailhog)
cd backend/docker
docker-compose up -d

# Run backend server
cd backend
go run main.go

# Run all tests
go test ./...

# Run specific test
go test ./controllers -run TestAccountIndex

# Run tests with timeout (recommended)
go test ./... -timeout 10m

# Run specific test packages for reliability
go test ./library/... -timeout 5m     # Library tests (all pass)
go test ./controllers -timeout 5m     # Main controller tests  
go test ./controllers/admin -timeout 5m # Admin controller tests

# Run with sequential execution to avoid test isolation issues
go test ./... -timeout 10m -p 1
```

### Frontend
```bash
cd frontend
npm install
npm start          # Dev server at http://localhost:4200
npm run build      # Production build
npm test           # Unit tests
npm run e2e        # E2E tests
```

## Key Patterns

### Backend Structure
- Controllers in `backend/controllers/` handle HTTP requests
- Models in `backend/models/` define data structures and DB operations
- One test file per source file (`file.go` → `file_test.go`)
- Environment config via `.env` file (copy from `.env.sample`)

### API Conventions
- Multi-tenant: URLs include account ID: `/api/v3/accounts/{accountId}/...`
- Authentication: Bearer token in Authorization header
- JSON responses with `response` wrapper for data

### Testing
- Backend: Standard Go testing with `nbio/st` for assertions
- Test database created/destroyed per test run
- Mock external services with `gock`

### Frontend Routing
- Main app routes: `/ledger`, `/settings`, `/contacts`, `/snapclerk`
- Auth required except for `/login`, `/register`, `/forgot-password`
- Angular services in `frontend/src/app/services/`

## Commit Message Guidelines
- Create commit messages that are detailed and explain the purpose of the change
- Follow a clear and consistent structure that provides context
- Include references to related issues or tickets when applicable