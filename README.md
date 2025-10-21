# viddl.me

A secure video downloader for YouTube, Twitter, Instagram, TikTok, and other social media platforms.

## Features

- **Multiple Platform Support**: YouTube, Twitter/X, Instagram, TikTok, Facebook, Vimeo, Reddit, Twitch
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

- URL validation and sanitization
- Domain whitelist (only allowed platforms)
- Rate limiting (3 requests per minute per IP)
- Input validation and regex-based sanitization for all parameters
- Shell injection prevention with proper argument escaping
- Automatic temp file cleanup (5-minute intervals)
- CORS protection with configurable origins
- No file size limits (handles large 4K videos)

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

Create a `.env` file in the backend directory:

```env
PORT=3000
ALLOWED_ORIGINS=http://localhost:5173,https://viddl.me
```

## Production Deployment

### Backend

```bash
cd backend
go build -o viddl-server main.go
./viddl-server
```

Or use systemd service:

```ini
[Unit]
Description=Viddl.me Backend
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/var/www/viddl.me/backend
ExecStart=/var/www/viddl.me/backend/viddl-server
Restart=on-failure
Environment="PORT=3000"
Environment="ALLOWED_ORIGINS=https://viddl.me"

[Install]
WantedBy=multi-user.target
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

**Response:**
```json
{
  "title": "Video Title",
  "thumbnail": "https://...",
  "duration": 180,
  "uploader": "Channel Name",
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

### POST /api/download

Download video.

**Request:**
```json
{
  "url": "https://youtube.com/watch?v=...",
  "format": "best"
}
```

**Response:** File download

### GET /health

Health check endpoint.

**Response:**
```json
{
  "status": "ok"
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

## License

MIT
