# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Basecoat** is a miniature painting tracker web app. Users can:
- Search for miniature sets/games via the Claude API to import mini lists into their collection
- Track painting status per mini (primed, basecoated, shaded, detailed, finished, etc.)
- Log paints used and technique notes per mini
- Store photos at different painting stages
- View a dashboard with in-progress work and a "grey pile of shame meter" for unpainted/unfinished minis

## Tech Stack

- **Backend**: Go (REST API)
- **Frontend**: React (TypeScript)
- **Database**: SQLite via `database/sql` with an abstraction layer to ease future migration to Postgres or similar
- **Image storage**: Local filesystem with a storage interface to ease future migration to S3 or similar
- **Auth**: User accounts with JWT-based authentication
- **External API**: Anthropic Claude API for importing miniature lists by game/set name

## Repository Structure (planned)

```
/
├── backend/          # Go REST API server
│   ├── cmd/server/   # Main entrypoint
│   ├── internal/
│   │   ├── api/      # HTTP handlers and routing
│   │   ├── auth/     # JWT auth middleware and helpers
│   │   ├── db/       # Database layer (sql abstraction, migrations)
│   │   ├── models/   # Domain types
│   │   ├── claude/   # Claude API client for mini list lookup
│   │   └── storage/  # File storage interface (local FS impl, easy S3 swap)
│   └── migrations/   # SQL migration files
├── frontend/         # React + TypeScript app
│   ├── src/
│   │   ├── api/      # API client (typed fetch wrappers)
│   │   ├── components/
│   │   ├── pages/
│   │   └── types/    # Shared TypeScript types matching backend models
│   └── public/
└── CLAUDE.md
```

## Key Architecture Decisions

### Database abstraction
Use a `Repository` interface pattern in `internal/db/` so that swapping SQLite for Postgres only requires a new implementation, not changes to business logic. All queries go through repository interfaces, never raw SQL in handlers.

### Storage abstraction
Define a `Storage` interface in `internal/storage/` with `Put`, `Get`, and `Delete` methods. The local filesystem implementation lives in the same package. Future S3 implementation drops in without touching upload handlers.

### Claude API integration
The `internal/claude/` package wraps the Anthropic API. Given a game name or set name, it returns a structured list of miniatures. The prompt should ask Claude to return JSON with mini names, unit types, and count. This is the only place the Anthropic SDK is used.

### Auth
JWT tokens issued on login, validated via middleware applied to all protected routes. Store hashed passwords (bcrypt) in the users table. Tokens should carry user ID and be validated on every request.

### Image uploads
Images are associated with a specific mini and a painting stage. Store metadata (mini ID, stage, file path/URL, timestamp) in the DB; store the actual file via the Storage interface. Return URLs from the API; the frontend never constructs file paths directly.

## Development Commands

### Backend
```bash
cd backend
air                          # Run dev server with hot reload (reads .env automatically)
go run ./cmd/server          # Run without hot reload
go test ./...                # Run all tests
go test ./internal/api/...   # Run tests for a specific package
go build ./cmd/server        # Build the binary
```

### Frontend
```bash
cd frontend
npm install                  # Install dependencies
npm run dev                  # Start Vite dev server (proxies /api and /uploads to :8080)
npm run build                # Production build
npm run lint                 # Lint
```

### Database migrations
Migrations are plain SQL files in `backend/migrations/`. Apply them in order on startup or via a dedicated migrate command (to be determined once migration tooling is chosen — consider `golang-migrate`).

## Environment Variables

The backend reads configuration from environment variables (or a `.env` file in development):

| Variable | Description |
|---|---|
| `DATABASE_PATH` | Path to SQLite file |
| `ANTHROPIC_API_KEY` | API key for Claude mini lookup |
| `JWT_SECRET` | Secret for signing JWTs |
| `STORAGE_PATH` | Root directory for local image storage |
| `PORT` | HTTP port (default 8080) |
