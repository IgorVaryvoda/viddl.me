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
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var allowedDomains = []string{
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
	limiters map[string]*rate.Limiter
}

func NewIPRateLimiter() *IPRateLimiter {
	return &IPRateLimiter{
		limiters: make(map[string]*rate.Limiter),
	}
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	if limiter, exists := i.limiters[ip]; exists {
		return limiter
	}
	limiter := rate.NewLimiter(rate.Every(time.Minute/3), 3)
	i.limiters[ip] = limiter
	return limiter
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
	re := regexp.MustCompile(`^[a-zA-Z0-9+]+$`)
	if format == "" || !re.MatchString(format) {
		return "best"
	}
	if len(format) > 50 {
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
		log.Printf("Using cookies file: %s", cookiesFile)
		args = append(args, "--cookies", cookiesFile)
	} else {
		log.Printf("WARNING: No cookies file configured (YTDLP_COOKIES not set)")
	}

	args = append(args, sanitizedURL)

	log.Printf("Running yt-dlp with args: %v", args)
	cmd := exec.Command("yt-dlp", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("yt-dlp error: %v, output: %s", err, string(output))
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

	tmpDir := filepath.Join(".", "tmp")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
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

	log.Printf("Downloading with format: %s for URL: %s", formatSpec, sanitizedURL)

	dlArgs := []string{"-f", formatSpec, "-o", outputTemplate, "--no-playlist", "--merge-output-format", "mp4"}

	cookiesFile := os.Getenv("YTDLP_COOKIES")
	if cookiesFile != "" {
		dlArgs = append(dlArgs, "--cookies", cookiesFile)
	}

	dlArgs = append(dlArgs, sanitizedURL)

	cmd := exec.Command("yt-dlp", dlArgs...)

	if err := cmd.Run(); err != nil {
		log.Printf("yt-dlp download error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Download failed"})
		return
	}

	files, err := filepath.Glob(filepath.Join(tmpDir, sessionID+".*"))
	if err != nil || len(files) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Downloaded file not found"})
		return
	}

	filePath := files[0]
	fileName := filepath.Base(filePath)

	defer func() {
		time.Sleep(60 * time.Second)
		os.Remove(filePath)
	}()

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.File(filePath)
}

func cleanupOldFiles() {
	tmpDir := filepath.Join(".", "tmp")
	ticker := time.NewTicker(5 * time.Minute)

	for range ticker.C {
		files, err := os.ReadDir(tmpDir)
		if err != nil {
			continue
		}

		now := time.Now()
		for _, file := range files {
			filePath := filepath.Join(tmpDir, file.Name())
			info, err := os.Stat(filePath)
			if err != nil {
				continue
			}

			if now.Sub(info.ModTime()) > 5*time.Minute {
				os.Remove(filePath)
				log.Printf("Cleaned up old file: %s", file.Name())
			}
		}
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

	limiter := NewIPRateLimiter()

	r.POST("/api/info", rateLimitMiddleware(limiter), getVideoInfo)
	r.POST("/api/download", rateLimitMiddleware(limiter), downloadVideo)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	go cleanupOldFiles()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
