package downloader

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"viddl.me/backend/internal/models"
)

type Downloader struct {
	tmpDir      string
	cookiesFile string
	maxFilesize string
}

func New(tmpDir, cookiesFile, maxFilesize string) *Downloader {
	return &Downloader{
		tmpDir:      tmpDir,
		cookiesFile: cookiesFile,
		maxFilesize: maxFilesize,
	}
}

func (d *Downloader) GetVideoInfo(videoURL string) (*models.VideoInfo, error) {
	multiVideos, err := d.checkMultipleVideos(videoURL)
	if err == nil && len(multiVideos) > 1 {
		return &models.VideoInfo{
			Title:        fmt.Sprintf("Multiple videos (%d)", len(multiVideos)),
			IsMultiVideo: true,
			MultiVideos:  multiVideos,
		}, nil
	}

	return d.getSingleVideoInfo(videoURL)
}

func (d *Downloader) checkMultipleVideos(videoURL string) ([]models.VideoEntry, error) {
	args := []string{"--flat-playlist", "--dump-json"}

	if d.cookiesFile != "" {
		args = append(args, "--cookies", d.cookiesFile)
	}
	args = append(args, videoURL)

	log.Printf("INFO: Checking for multiple videos with args: %v", args)
	cmd := exec.Command("yt-dlp", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var videos []models.VideoEntry
	lines := strings.Split(string(output), "\n")
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		var entry map[string]any
		if json.Unmarshal([]byte(line), &entry) == nil {
			if entry["_type"] == "url" || entry["_type"] == "video" {
				video := models.VideoEntry{Index: i + 1}
				if t, ok := entry["title"].(string); ok {
					video.Title = t
				}
				if t, ok := entry["thumbnail"].(string); ok {
					video.Thumbnail = t
				}
				if d, ok := entry["duration"].(float64); ok {
					video.Duration = d
				}
				videos = append(videos, video)
			}
		}
	}
	return videos, nil
}

func (d *Downloader) getSingleVideoInfo(videoURL string) (*models.VideoInfo, error) {
	args := []string{"--dump-json", "--no-playlist"}

	if d.cookiesFile != "" {
		log.Printf("INFO: Using cookies file: %s", d.cookiesFile)
		args = append(args, "--cookies", d.cookiesFile)
	}
	args = append(args, videoURL)

	log.Printf("INFO: Running yt-dlp with args: %v", args)
	cmd := exec.Command("yt-dlp", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("ERROR: yt-dlp error: %v, output: %s", err, string(output))
		return nil, fmt.Errorf("failed to fetch video information")
	}

	var ytdlpInfo models.YtDlpInfo
	if err := json.Unmarshal(output, &ytdlpInfo); err != nil {
		return nil, fmt.Errorf("failed to parse video information")
	}

	isInstagram := strings.Contains(strings.ToLower(videoURL), "instagram.com")
	formats := d.extractFormats(ytdlpInfo, isInstagram)

	return &models.VideoInfo{
		Title:        ytdlpInfo.Title,
		Thumbnail:    ytdlpInfo.Thumbnail,
		Duration:     ytdlpInfo.Duration,
		Uploader:     ytdlpInfo.Uploader,
		Formats:      formats,
		IsMultiVideo: false,
	}, nil
}

func (d *Downloader) extractFormats(info models.YtDlpInfo, isInstagram bool) []models.FormatInfo {
	var formats []models.FormatInfo
	seen := make(map[int]bool)

	for _, f := range info.Formats {
		if f.VCodec == "none" || f.VCodec == "" {
			continue
		}
		if strings.Contains(strings.ToLower(f.FormatNote), "storyboard") {
			continue
		}
		if f.Height <= 0 || seen[f.Height] {
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
			if isInstagram {
				if info.Duration > 0 {
					estimatedSize = int64(info.Duration * float64(f.Height) * 200)
				} else {
					estimatedSize = int64(f.Height * f.Height * 80)
				}
			} else {
				estimatedSize = int64(f.Height * f.Height * 100)
			}
		}

		formats = append(formats, models.FormatInfo{
			FormatID: f.FormatID,
			Ext:      "mp4",
			Quality:  quality,
			Filesize: estimatedSize,
		})

		if len(formats) >= 15 {
			break
		}
	}
	return formats
}

type DownloadResult struct {
	FilePath    string
	FileName    string
	FileSize    int64
	ContentType string
}

func (d *Downloader) Download(videoURL, format string, videoIndex int) (*DownloadResult, error) {
	if err := os.MkdirAll(d.tmpDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	sessionID, err := generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	outputTemplate := filepath.Join(d.tmpDir, sessionID+".%(ext)s")
	args := d.buildDownloadArgs(videoURL, format, outputTemplate, videoIndex)

	log.Printf("INFO: Running yt-dlp download with args: %v", args)
	cmd := exec.Command("yt-dlp", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("ERROR: yt-dlp download error: %v, output: %s", err, string(output))
		return nil, fmt.Errorf("download failed or file exceeds size limit")
	}

	files, err := filepath.Glob(filepath.Join(d.tmpDir, sessionID+".*"))
	if err != nil || len(files) == 0 {
		return nil, fmt.Errorf("downloaded file not found")
	}

	filePath := files[0]
	fileName := filepath.Base(filePath)

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file info: %w", err)
	}

	contentType := "video/mp4"

	log.Printf("INFO: File size: %d bytes (%.2f MB)", fileInfo.Size(), float64(fileInfo.Size())/(1024*1024))

	return &DownloadResult{
		FilePath:    filePath,
		FileName:    fileName,
		FileSize:    fileInfo.Size(),
		ContentType: contentType,
	}, nil
}

func (d *Downloader) buildDownloadArgs(videoURL, format, outputTemplate string, videoIndex int) []string {
	isInstagram := strings.Contains(strings.ToLower(videoURL), "instagram.com")

	formatSpec := "bv*[ext=mp4]+ba[ext=m4a]/b[ext=mp4]/bv*+ba/b"
	if format != "best" {
		formatSpec = fmt.Sprintf("%s+ba/%s", format, format)
	}

	if isInstagram {
		formatSpec = "1/best[vcodec^=avc]/best[ext=mp4]/best"
		log.Printf("INFO: Instagram URL detected, using format 1 (H.264 baseline for iPhone)")
	}

	log.Printf("INFO: Downloading with format: %s", formatSpec)
	args := []string{"-f", formatSpec, "-o", outputTemplate, "--merge-output-format", "mp4"}

	if videoIndex > 0 {
		args = append(args, "--playlist-items", fmt.Sprintf("%d", videoIndex))
		log.Printf("INFO: Downloading video index: %d", videoIndex)
	} else {
		args = append(args, "--no-playlist")
	}

	args = append(args, "--max-filesize", d.maxFilesize)

	if d.cookiesFile != "" {
		args = append(args, "--cookies", d.cookiesFile)
	}

	args = append(args, videoURL)
	return args
}

func (d *Downloader) CheckHealth() error {
	cmd := exec.Command("yt-dlp", "--version")
	return cmd.Run()
}

func generateSessionID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
