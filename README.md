# Basecoat

A miniature painting tracker for tabletop gamers. Organise your collection by game and box, track the painting status of every miniature, log paints and techniques, and store progress photos. Uses the Claude API to automatically suggest miniature lists when you add a new box.

## Features

- **Collection hierarchy**: Collection → Game → Box → Miniature
- **Painting status tracking**: unpainted → primed → basecoated → shaded → detailed → finished
- **Grey pile of shame meter**: dashboard shows how much of your collection is still unpainted
- **Claude-powered mini lookup**: add a box name and Claude suggests the miniatures inside it
- **Shared catalog**: crowd-sourced mini lists are contributed to a shared catalog so the next person gets instant results with no API call
- **Paint log**: record which paints you used on each miniature and for what purpose
- **Photo storage**: attach photos at different painting stages

## Prerequisites

| Tool | Version | Install |
|------|---------|---------|
| Go | 1.26+ | https://go.dev/dl |
| Node.js | 18+ | https://nodejs.org |
| `air` | latest | `go install github.com/air-verse/air@latest` |

SQLite is embedded via `go-sqlite3` (CGO) — no separate database installation needed.

## Project structure

```
backend/    Go REST API (port 8080)
frontend/   React + TypeScript + Vite (port 5173 in dev)
```

## Setup

### 1. Backend environment

Create `backend/.env`:

```env
DATABASE_PATH=basecoat.db
JWT_SECRET=your-secret-key-here
STORAGE_PATH=uploads
PORT=8080
BASE_URL=http://localhost:8080
MIGRATIONS_DIR=migrations
ANTHROPIC_API_KEY=sk-ant-...   # optional — mini lookup is disabled if not set
```

`JWT_SECRET` can be any long random string. Generate one with:
```bash
openssl rand -hex 32
```

`ANTHROPIC_API_KEY` is available at https://console.anthropic.com. The app works without it — Claude lookup is simply disabled and you can add minis manually.

### 2. Install frontend dependencies

```bash
cd frontend
npm install
```

## Running in development

Open two terminals:

**Terminal 1 — backend** (hot reload via `air`):
```bash
cd backend
air
```

**Terminal 2 — frontend** (hot reload via Vite):
```bash
cd frontend
npm run dev
```

Then open http://localhost:5173.

On first startup the backend seeds a test account and sample data:

| Field | Value |
|-------|-------|
| Email | `test@test.com` |
| Password | `password` |

## Running in production

```bash
# Build the frontend
cd frontend && npm run build

# Build and run the backend (serves the built frontend via /uploads)
cd backend && go build -o basecoat ./cmd/server && ./basecoat
```

For production, set a strong `JWT_SECRET` and an absolute path for `DATABASE_PATH`.

## Environment variables

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DATABASE_PATH` | No | `basecoat.db` | Path to the SQLite file |
| `JWT_SECRET` | Yes | — | Secret used to sign JWT tokens |
| `STORAGE_PATH` | No | `uploads` | Directory for uploaded images |
| `PORT` | No | `8080` | HTTP port |
| `BASE_URL` | No | `http://localhost:8080` | Base URL used to construct image URLs |
| `MIGRATIONS_DIR` | No | auto-detected | Path to the `migrations/` directory |
| `ANTHROPIC_API_KEY` | No | — | Enables Claude mini lookup |

## Tech stack

- **Backend**: Go, [Chi](https://github.com/go-chi/chi) router, SQLite via `go-sqlite3`
- **Frontend**: React, TypeScript, Vite
- **Auth**: JWT (HS256), bcrypt password hashing
- **AI**: Anthropic Claude API (`claude-sonnet-4-6`) for miniature list lookup
- **Hot reload**: `air` (backend), Vite HMR (frontend)
