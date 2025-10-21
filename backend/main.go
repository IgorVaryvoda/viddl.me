package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var allowedDomains []string

func init() {
	// Configure structured logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("[viddl.me] ")

	// Load allowed domains from environment variable or use defaults
	domainsEnv := os.Getenv("ALLOWED_DOMAINS")
	if domainsEnv != "" {
		allowedDomains = strings.Split(domainsEnv, ",")
		for i, domain := range allowedDomains {
			allowedDomains[i] = strings.TrimSpace(domain)
		}
	} else {
		allowedDomains = []string{
			"youtube.com",
			"youtu.be",
			"twitter.com",
			"x.com",
			"instagram.com",
			"facebook.com",
			"tiktok.com",
			"vimeo.com",
			"reddit.com",
			"twitch.tv",
		}
	}
	log.Printf("INFO: Allowed domains configured: %v", allowedDomains)
}

type VideoInfo struct {
	Title     string       `json:"title"`
	Thumbnail string       `json:"thumbnail"`
	Duration  float64      `json:"duration"`
	Uploader  string       `json:"uploader"`
	Formats   []FormatInfo `json:"formats"`
}

type FormatInfo struct {
	FormatID string `json:"format_id"`
	Ext      string `json:"ext"`
	Quality  string `json:"quality"`
	Filesize int64  `json:"filesize"`
}

type YtDlpFormat struct {
	FormatID   string  `json:"format_id"`
	Ext        string  `json:"ext"`
	FormatNote string  `json:"format_note"`
	Quality    float64 `json:"quality"`
	Filesize   int64   `json:"filesize"`
	VCodec     string  `json:"vcodec"`
	ACodec     string  `json:"acodec"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	Fps        float64 `json:"fps"`
	Resolution string  `json:"resolution"`
}

type YtDlpInfo struct {
	Title     string        `json:"title"`
	Thumbnail string        `json:"thumbnail"`
	Duration  float64       `json:"duration"`
	Uploader  string        `json:"uploader"`
	Formats   []YtDlpFormat `json:"formats"`
}

type DownloadRequest struct {
	URL    string `json:"url" binding:"required"`
	Format string `json:"format"`
}

type IPRateLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*rate.Limiter
	lastSeen map[string]time.Time
}

func NewIPRateLimiter() *IPRateLimiter {
	limiter := &IPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
		lastSeen: make(map[string]time.Time),
	}
	go limiter.cleanupOldLimiters()
	return limiter
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	limiter, exists := i.limiters[ip]
	i.mu.RUnlock()

	if exists {
		i.mu.Lock()
		i.lastSeen[ip] = time.Now()
		i.mu.Unlock()
		return limiter
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	// Double-check in case another goroutine created it
	if limiter, exists := i.limiters[ip]; exists {
		i.lastSeen[ip] = time.Now()
		return limiter
	}

	limiter = rate.NewLimiter(rate.Every(time.Minute/3), 3)
	i.limiters[ip] = limiter
	i.lastSeen[ip] = time.Now()
	return limiter
}

func (i *IPRateLimiter) cleanupOldLimiters() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		i.mu.Lock()
		now := time.Now()
		for ip, lastSeen := range i.lastSeen {
			if now.Sub(lastSeen) > 30*time.Minute {
				delete(i.limiters, ip)
				delete(i.lastSeen, ip)
				log.Printf("INFO: Cleaned up rate limiter for IP: %s", ip)
			}
		}
		i.mu.Unlock()
	}
}

func securityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		c.Next()
	}
}

func rateLimitMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.GetLimiter(ip).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func sanitizeURL(inputURL string) (string, error) {
	if len(inputURL) > 2048 {
		return "", fmt.Errorf("URL too long")
	}

	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL format")
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("invalid protocol")
	}

	hostname := strings.TrimPrefix(parsedURL.Hostname(), "www.")
	allowed := false
	for _, domain := range allowedDomains {
		if hostname == domain || strings.HasSuffix(hostname, "."+domain) {
			allowed = true
			break
		}
	}

	if !allowed {
		return "", fmt.Errorf("domain not allowed")
	}

	return parsedURL.String(), nil
}

func sanitizeFormat(format string) string {
	// Use strict whitelist for format IDs - only alphanumeric characters
	// No special characters allowed to prevent any potential injection
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if format == "" || !re.MatchString(format) {
		return "best"
	}
	if len(format) > 20 {
		return "best"
	}
	return format
}

func getVideoInfo(c *gin.Context) {
	var req DownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	sanitizedURL, err := sanitizeURL(req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	args := []string{"--dump-json", "--no-playlist"}

	cookiesFile := os.Getenv("YTDLP_COOKIES")
	if cookiesFile != "" {
		log.Printf("INFO: Using cookies file: %s", cookiesFile)
		args = append(args, "--cookies", cookiesFile)
	} else {
		log.Printf("WARN: No cookies file configured (YTDLP_COOKIES not set)")
	}

	args = append(args, sanitizedURL)

	log.Printf("INFO: Running yt-dlp with args: %v", args)
	cmd := exec.Command("yt-dlp", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("ERROR: yt-dlp error: %v, output: %s", err, string(output))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch video information"})
		return
	}

	var ytdlpInfo YtDlpInfo
	if err := json.Unmarshal(output, &ytdlpInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse video information"})
		return
	}

	formats := []FormatInfo{}
	seen := make(map[int]bool)

	for _, f := range ytdlpInfo.Formats {
		if f.VCodec == "none" || f.VCodec == "" {
			continue
		}

		if strings.Contains(strings.ToLower(f.FormatNote), "storyboard") {
			continue
		}

		if f.Height <= 0 {
			continue
		}

		if seen[f.Height] {
			continue
		}

		quality := f.FormatNote
		if quality == "" && f.Height > 0 {
			quality = fmt.Sprintf("%dp", f.Height)
		} else if quality == "" && f.Resolution != "" {
			quality = f.Resolution
		}
		if quality == "" || quality == "0" {
			quality = "unknown"
		}

		seen[f.Height] = true

		estimatedSize := f.Filesize
		if estimatedSize == 0 && f.Height > 0 {
			estimatedSize = int64(f.Height * f.Height * 100)
		}

		formats = append(formats, FormatInfo{
			FormatID: f.FormatID,
			Ext:      "mp4",
			Quality:  quality,
			Filesize: estimatedSize,
		})

		if len(formats) >= 15 {
			break
		}
	}

	info := VideoInfo{
		Title:     ytdlpInfo.Title,
		Thumbnail: ytdlpInfo.Thumbnail,
		Duration:  ytdlpInfo.Duration,
		Uploader:  ytdlpInfo.Uploader,
		Formats:   formats,
	}

	c.JSON(http.StatusOK, info)
}

func downloadVideo(c *gin.Context) {
	var req DownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	sanitizedURL, err := sanitizeURL(req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	format := sanitizeFormat(req.Format)
	log.Printf("INFO: Download request from %s for URL: %s, format: %s", c.ClientIP(), sanitizedURL, format)

	tmpDir := filepath.Join(".", "tmp")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		log.Printf("ERROR: Failed to create temp directory: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temp directory"})
		return
	}

	sessionBytes := make([]byte, 16)
	if _, err := rand.Read(sessionBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session ID"})
		return
	}
	sessionID := hex.EncodeToString(sessionBytes)

	outputTemplate := filepath.Join(tmpDir, sessionID+".%(ext)s")

	formatSpec := "bestvideo[ext=mp4]+bestaudio[ext=m4a]/bestvideo+bestaudio/best"
	if format != "best" {
		formatSpec = fmt.Sprintf("%s+bestaudio/%s+bestaudio[ext=m4a]", format, format)
	}

	log.Printf("INFO: Downloading with format: %s for URL: %s", formatSpec, sanitizedURL)

	dlArgs := []string{"-f", formatSpec, "-o", outputTemplate, "--no-playlist", "--merge-output-format", "mp4"}

	// Add max filesize limit (2GB default, configurable via env)
	maxFilesize := os.Getenv("MAX_DOWNLOAD_SIZE")
	if maxFilesize == "" {
		maxFilesize = "2G"
	}
	dlArgs = append(dlArgs, "--max-filesize", maxFilesize)

	cookiesFile := os.Getenv("YTDLP_COOKIES")
	if cookiesFile != "" {
		dlArgs = append(dlArgs, "--cookies", cookiesFile)
	}

	dlArgs = append(dlArgs, sanitizedURL)

	log.Printf("INFO: Running yt-dlp download with args: %v", dlArgs)
	cmd := exec.Command("yt-dlp", dlArgs...)

	if err := cmd.Run(); err != nil {
		log.Printf("ERROR: yt-dlp download error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Download failed or file exceeds size limit"})
		return
	}

	files, err := filepath.Glob(filepath.Join(tmpDir, sessionID+".*"))
	if err != nil || len(files) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Downloaded file not found"})
		return
	}

	filePath := files[0]
	fileName := filepath.Base(filePath)

	// Cleanup file asynchronously after serving
	log.Printf("INFO: Serving file: %s to client: %s", fileName, c.ClientIP())
	go func(path string) {
		time.Sleep(60 * time.Second)
		if err := os.Remove(path); err != nil {
			log.Printf("ERROR: Failed to cleanup file %s: %v", path, err)
		} else {
			log.Printf("INFO: Cleaned up downloaded file: %s", path)
		}
	}(filePath)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "video/mp4")
	c.File(filePath)
}

func cleanupOldFiles() {
	tmpDir := filepath.Join(".", "tmp")
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// Initial cleanup on startup
	cleanupOnce(tmpDir)

	for range ticker.C {
		cleanupOnce(tmpDir)
	}
}

func cleanupOnce(tmpDir string) {
	files, err := os.ReadDir(tmpDir)
	if err != nil {
		log.Printf("ERROR: Error reading tmp directory: %v", err)
		return
	}

	now := time.Now()
	cleanedCount := 0
	var totalSize int64

	for _, file := range files {
		// Skip hidden files like .health_check
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		filePath := filepath.Join(tmpDir, file.Name())
		info, err := os.Stat(filePath)
		if err != nil {
			log.Printf("ERROR: Error stating file %s: %v", file.Name(), err)
			continue
		}

		// Clean up files older than 5 minutes
		if now.Sub(info.ModTime()) > 5*time.Minute {
			fileSize := info.Size()
			if err := os.Remove(filePath); err != nil {
				log.Printf("ERROR: Error removing file %s: %v", file.Name(), err)
			} else {
				cleanedCount++
				totalSize += fileSize
				log.Printf("INFO: Cleaned up old file: %s (size: %d bytes, age: %v)",
					file.Name(), fileSize, now.Sub(info.ModTime()).Round(time.Second))
			}
		}
	}

	if cleanedCount > 0 {
		log.Printf("INFO: Cleanup complete: removed %d files, freed %d MB",
			cleanedCount, totalSize/(1024*1024))
	}
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:5173,http://localhost:8080"
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(allowedOrigins, ","),
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Add security headers middleware
	r.Use(securityHeadersMiddleware())

	limiter := NewIPRateLimiter()

	r.POST("/api/info", rateLimitMiddleware(limiter), getVideoInfo)
	r.POST("/api/download", rateLimitMiddleware(limiter), downloadVideo)
	r.GET("/health", func(c *gin.Context) {
		// Check if yt-dlp is available
		cmd := exec.Command("yt-dlp", "--version")
		if err := cmd.Run(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  "yt-dlp not available",
			})
			return
		}

		// Check if tmp directory is writable
		tmpDir := filepath.Join(".", "tmp")
		testFile := filepath.Join(tmpDir, ".health_check")
		if err := os.WriteFile(testFile, []byte("ok"), 0644); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unhealthy",
				"error":  "tmp directory not writable",
			})
			return
		}
		os.Remove(testFile)

		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	go cleanupOldFiles()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("INFO: Server starting on port %s", port)
	log.Printf("INFO: Environment: ALLOWED_ORIGINS=%s, MAX_DOWNLOAD_SIZE=%s",
		os.Getenv("ALLOWED_ORIGINS"), os.Getenv("MAX_DOWNLOAD_SIZE"))
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("FATAL: Server failed to start: %v", err)
	}
}
