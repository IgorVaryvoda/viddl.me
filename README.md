# viddl.me

A secure video downloader for YouTube, Twitter, Instagram, TikTok, and other social media platforms.

## Features

- **Multiple Platform Support**: YouTube, Twitter/X, Instagram, TikTok, Facebook, Vimeo, Reddit, Twitch
- **Multi-Video Support**: Automatically detects and allows selection from Twitter posts with multiple videos
- **Secure Backend**: Go with Gin framework
- **Input Sanitization**: URL validation and domain whitelisting
- **Rate Limiting**: Prevents abuse with IP-based rate limiting (3 requests/minute per IP)
- **Quality Selection**: Choose from 144p to 4K resolution
- **MP4-Only Output**: All videos automatically converted to MP4 format
- **Smart Format Merging**: Automatically merges video and audio streams for best quality
- **Download Progress**: Real-time progress tracking during downloads
- **Lightweight Frontend**: Vue.js 3 with Vite and dark theme UI
- **Automatic Cleanup**: Temporary files cleaned up every 5 minutes

## Security Features

- **Thread-safe rate limiting** with automatic cleanup
- **URL validation and sanitization** with domain whitelist
- **Command injection protection** with strict input validation
- **Security headers** (X-Content-Type-Options, X-Frame-Options, CSP, etc.)
- **File size limits** (configurable, 2GB default)
- **Automatic temp file cleanup** (5-minute intervals)
- **CORS protection** with configurable origins
- **Structured logging** with severity levels
- **Health check monitoring** (yt-dlp and filesystem checks)

## Prerequisites

- Go 1.21 or higher
- Node.js 18 or higher
- yt-dlp installed and available in PATH
- ffmpeg (required for merging video/audio streams)

### Installing yt-dlp

**Linux/macOS:**
```bash
sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
sudo chmod a+rx /usr/local/bin/yt-dlp
```

**Or via pip:**
```bash
pip install yt-dlp
```

### Installing ffmpeg

**Linux (Debian/Ubuntu):**
```bash
sudo apt update
sudo apt install ffmpeg
```

**macOS:**
```bash
brew install ffmpeg
```

**Arch Linux:**
```bash
sudo pacman -S ffmpeg
```

## Setup

### Backend (Go)

```bash
cd backend
go mod download
cp .env.example .env
# Edit .env if needed
PORT=3000 go run main.go
```

The backend will start on port 3000 by default and listen on `http://localhost:3000`.

### Frontend (Vue.js)

```bash
cd frontend
npm install
npm run dev
```

The frontend will start on `http://localhost:5173` with hot module replacement enabled. The Vite dev server automatically proxies `/api` requests to the backend.

## Environment Variables

Create a `.env` file in the backend directory or set these environment variables:

```env
# Server Configuration
PORT=3000                                    # Server port (default: 3000)

# CORS Configuration
ALLOWED_ORIGINS=http://localhost:5173,https://viddl.me  # Comma-separated allowed origins

# Domain Whitelist (optional)
ALLOWED_DOMAINS=youtube.com,youtu.be,twitter.com,x.com,instagram.com,facebook.com,tiktok.com,vimeo.com,reddit.com,twitch.tv
# If not set, uses default list above

# Download Limits
MAX_DOWNLOAD_SIZE=2G                         # Maximum file size (e.g., 2G, 500M) (default: 2G)

# yt-dlp Configuration
YTDLP_COOKIES=/path/to/cookies.txt          # Optional: Path to cookies file for authenticated downloads
```

### Environment Variable Details

- **PORT**: The port on which the backend server runs (default: 3000)
- **ALLOWED_ORIGINS**: Comma-separated list of allowed CORS origins for the frontend
- **ALLOWED_DOMAINS**: Comma-separated list of allowed video platform domains (overrides defaults)
- **MAX_DOWNLOAD_SIZE**: Maximum allowed file size for downloads (uses yt-dlp syntax: K, M, G)
- **YTDLP_COOKIES**: Path to a Netscape-format cookies file for downloading age-restricted or private videos

## Production Deployment

### Backend

```bash
cd backend
go build -o viddl-server main.go
./viddl-server
```

Or use systemd service. Create `/etc/systemd/system/viddl.service`:

```ini
[Unit]
Description=Viddl.me Backend
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/viddl.me/backend
ExecStart=/var/www/viddl.me/backend/viddl-server
Restart=always
RestartSec=5s

Environment="PORT=3000"
Environment="ALLOWED_ORIGINS=https://viddl.me,https://www.viddl.me"
Environment="MAX_DOWNLOAD_SIZE=2G"
Environment="YTDLP_COOKIES=/var/www/viddl.me/cookies.txt"

# Logging
StandardOutput=journal
StandardError=journal
SyslogIdentifier=viddl

# Security
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

**Enable and start the service:**
```bash
sudo systemctl daemon-reload
sudo systemctl enable viddl.service
sudo systemctl start viddl.service
sudo systemctl status viddl.service
```

**View logs:**
```bash
sudo journalctl -u viddl.service -f
```

### Frontend

```bash
cd frontend
npm run build
```

Serve the `dist` folder with nginx or any static file server.

### Nginx Configuration Example

```nginx
server {
    listen 80;
    server_name viddl.me;

    location / {
        root /var/www/viddl.me/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_cache_bypass $http_upgrade;
    }
}
```

## API Endpoints

### POST /api/info

Get video information without downloading.

**Request:**
```json
{
  "url": "https://youtube.com/watch?v=..."
}
```

**Response (Single Video):**
```json
{
  "title": "Video Title",
  "thumbnail": "https://...",
  "duration": 180,
  "uploader": "Channel Name",
  "is_multi_video": false,
  "formats": [
    {
      "format_id": "137",
      "ext": "mp4",
      "quality": "1080p",
      "filesize": 52428800
    }
  ]
}
```

**Response (Multiple Videos - e.g., Twitter post with multiple videos):**
```json
{
  "title": "Multiple videos (3)",
  "is_multi_video": true,
  "multi_videos": [
    {
      "index": 1,
      "title": "Video 1",
      "thumbnail": "https://...",
      "duration": 60
    },
    {
      "index": 2,
      "title": "Video 2",
      "thumbnail": "https://...",
      "duration": 45
    },
    {
      "index": 3,
      "title": "Video 3",
      "thumbnail": "https://...",
      "duration": 30
    }
  ]
}
```

### POST /api/download

Download video.

**Request (Single Video or Best Quality):**
```json
{
  "url": "https://youtube.com/watch?v=...",
  "format": "best"
}
```

**Request (Specific Video from Multi-Video Post):**
```json
{
  "url": "https://twitter.com/user/status/...",
  "format": "best",
  "video_index": 2
}
```

**Response:** File download

### GET /health

Health check endpoint that verifies yt-dlp availability and filesystem writability.

**Response (Healthy):**
```json
{
  "status": "healthy",
  "version": "1.0.0"
}
```

**Response (Unhealthy):**
```json
{
  "status": "unhealthy",
  "error": "yt-dlp not available"
}
```

## Development

### Backend Development

```bash
cd backend
go run main.go
```

### Frontend Development

```bash
cd frontend
npm run dev
```

The frontend dev server proxies API requests to the backend.

## Security Considerations

- Always run behind a reverse proxy (nginx/caddy) in production
- Use HTTPS in production
- Keep yt-dlp updated regularly
- Monitor disk usage for the tmp directory
- Consider implementing user accounts and quotas for production
- Add CAPTCHA for additional abuse prevention
- Monitor and adjust rate limits based on your needs
- Set appropriate MAX_DOWNLOAD_SIZE to prevent disk space exhaustion
- Regularly review logs for suspicious activity
- Use the /health endpoint for monitoring and alerting

## Recent Security Improvements

- **v1.0.0** (2025-10-21):
  - Added thread-safe rate limiting with mutex protection
  - Fixed memory leak in rate limiter with automatic cleanup
  - Strengthened format parameter validation to prevent command injection
  - Fixed blocking file cleanup issue by using goroutines
  - Added comprehensive security headers middleware
  - Implemented proper health check monitoring
  - Added file size limits with configurable MAX_DOWNLOAD_SIZE
  - Improved structured logging with severity levels
  - Made domain whitelist configurable via environment variables
  - Optimized file cleanup mechanism with better error handling

## License

MIT
