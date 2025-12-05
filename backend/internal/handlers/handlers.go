package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"viddl.me/backend/internal/cleanup"
	"viddl.me/backend/internal/config"
	"viddl.me/backend/internal/downloader"
	"viddl.me/backend/internal/models"
)

type Handler struct {
	cfg        *config.Config
	downloader *downloader.Downloader
}

func New(cfg *config.Config) *Handler {
	dl := downloader.New(cfg.TmpDir, cfg.CookiesFile, cfg.MaxDownloadSize)
	return &Handler{
		cfg:        cfg,
		downloader: dl,
	}
}

func (h *Handler) GetVideoInfo(c *gin.Context) {
	var req models.DownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	sanitizedURL, err := downloader.SanitizeURL(req.URL, h.cfg.AllowedDomains)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	info, err := h.downloader.GetVideoInfo(sanitizedURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}

func (h *Handler) DownloadVideo(c *gin.Context) {
	var req models.DownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	sanitizedURL, err := downloader.SanitizeURL(req.URL, h.cfg.AllowedDomains)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	format := downloader.SanitizeFormat(req.Format)
	log.Printf("INFO: Download request from %s for URL: %s, format: %s",
		c.ClientIP(), sanitizedURL, format)

	result, err := h.downloader.Download(sanitizedURL, format, req.VideoIndex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("INFO: Serving file: %s to client: %s", result.FileName, c.ClientIP())
	cleanup.ScheduleFileRemoval(result.FilePath, 60*time.Second)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+result.FileName)
	c.Header("Content-Type", result.ContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", result.FileSize))
	c.File(result.FilePath)
}

func (h *Handler) ExtractAudio(c *gin.Context) {
	var req models.AudioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	sanitizedURL, err := downloader.SanitizeURL(req.URL, h.cfg.AllowedDomains)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("INFO: Audio extraction request from %s for URL: %s, format: %s",
		c.ClientIP(), sanitizedURL, req.AudioFormat)

	result, err := h.downloader.ExtractAudio(sanitizedURL, req.AudioFormat, req.VideoIndex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("INFO: Serving audio file: %s to client: %s", result.FileName, c.ClientIP())
	cleanup.ScheduleFileRemoval(result.FilePath, 60*time.Second)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+result.FileName)
	c.Header("Content-Type", result.ContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", result.FileSize))
	c.File(result.FilePath)
}

func (h *Handler) HealthCheck(c *gin.Context) {
	if err := h.downloader.CheckHealth(); err != nil {
		c.JSON(http.StatusServiceUnavailable, models.HealthResponse{
			Status: "unhealthy",
			Error:  "yt-dlp not available",
		})
		return
	}

	testFile := filepath.Join(h.cfg.TmpDir, ".health_check")
	if err := os.WriteFile(testFile, []byte("ok"), 0644); err != nil {
		c.JSON(http.StatusServiceUnavailable, models.HealthResponse{
			Status: "unhealthy",
			Error:  "tmp directory not writable",
		})
		return
	}
	os.Remove(testFile)

	c.JSON(http.StatusOK, models.HealthResponse{
		Status:  "healthy",
		Version: "1.0.0",
	})
}
