# AGENTS.md

## Project
viddl.me video downloader.

- Path: `/home/igor/Projects/viddl.me`
- Stack: Go/Gin backend, yt-dlp + ffmpeg, Vue 3/Vite frontend.
- Backend: `backend/`
- Frontend: `frontend/`

## Commands

```bash
make dev-backend
make dev-frontend
make test
make build
make quality-gate
```

Equivalent raw commands:

```bash
cd backend && go test ./...
cd backend && go build -o viddl-server main.go
cd frontend && npm run build
```

## Quality gate

Run `make quality-gate` before finishing code changes. If only docs changed, skip explicitly.

## Agent workflow

1. Read `CLAUDE.md` first for architecture and security model.
2. Security-sensitive app: never loosen URL validation, CORS, rate limiting, concurrent limits, file-size limits, cleanup, or command execution without explicit reason.
3. Do not print `.env` values or secrets.
4. Do not commit generated binaries. `backend/backend` and `backend/viddl-server` should stay untracked/ignored.
5. For yt-dlp behavior changes, test metadata fetch/download paths with safe sample URLs where possible.
6. Keep frontend API assumptions aligned with backend endpoints.

## Current caveats

- The repo had an untracked `backend/backend` binary; ignore/remove generated binaries before committing.
