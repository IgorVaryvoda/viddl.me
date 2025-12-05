# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

viddl.me is a video downloader web application supporting YouTube, Twitter/X, Instagram, TikTok, and other platforms. It uses yt-dlp as the underlying download engine.

## Commands

### Backend (Go)
```bash
cd backend
go mod download          # Install dependencies
PORT=3000 go run main.go # Run development server
go build -o viddl-server main.go  # Build for production
go test ./...            # Run all tests
go test ./internal/downloader     # Run tests for a specific package
```

### Frontend (Vue.js)
```bash
cd frontend
npm install              # Install dependencies
npm run dev              # Run dev server (port 5173, proxies /api to backend)
npm run build            # Build for production (outputs to dist/)
```

## Architecture

### Backend Structure (`backend/`)
- `main.go` - Entry point, sets up Gin router with CORS, rate limiting, and routes
- `internal/config/` - Configuration loading from environment variables
- `internal/handlers/` - HTTP handlers for `/api/info`, `/api/download`, `/health`
- `internal/downloader/` - yt-dlp wrapper and URL sanitization logic
- `internal/middleware/` - Rate limiting (IP-based), concurrent download limiting, security headers, gzip
- `internal/models/` - Request/response structs
- `internal/cleanup/` - Automatic temp file cleanup (runs every 5 minutes)

### Frontend Structure (`frontend/`)
- Single-page Vue 3 app with Vite
- `src/App.vue` - Main component containing all UI logic
- `src/style.css` - Dark theme styling
- Vite proxies `/api` requests to backend during development

### API Endpoints
- `POST /api/info` - Get video metadata and available formats
- `POST /api/download` - Download video (supports `video_index` for multi-video posts)
- `GET /health` - Health check (verifies yt-dlp and filesystem)

## Key Implementation Details

- Rate limiting: 3 requests/minute per IP via `golang.org/x/time/rate`
- Concurrent download limit: 2 per IP
- All videos converted to MP4 format
- Temp files stored in `backend/tmp/` with automatic cleanup
- URL validation uses domain whitelist (configurable via `ALLOWED_DOMAINS` env var)
- File size limit configurable via `MAX_DOWNLOAD_SIZE` (default 2GB)

## Environment Variables

Key variables (set in backend `.env` or environment):
- `PORT` - Server port (default: 3000)
- `ALLOWED_ORIGINS` - CORS origins (comma-separated)
- `ALLOWED_DOMAINS` - Video platform whitelist (comma-separated)
- `MAX_DOWNLOAD_SIZE` - Max file size (e.g., "2G", "500M")
- `YTDLP_COOKIES` - Path to cookies file for authenticated downloads
