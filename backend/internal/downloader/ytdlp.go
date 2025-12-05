package downloader

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"viddl.me/backend/internal/models"
)

type Downloader struct {
	tmpDir        string
	cookiesFile   string
	maxFilesize   string
	healthChecked bool
	healthError   error
}

func New(tmpDir, cookiesFile, maxFilesize string) *Downloader {
	return &Downloader{
		tmpDir:      tmpDir,
		cookiesFile: cookiesFile,
		maxFilesize: maxFilesize,
	}
}

func (d *Downloader) GetVideoInfo(videoURL string) (*models.VideoInfo, error) {
	// Skip multi-video check for YouTube single videos (not playlists)
	isYouTube := strings.Contains(strings.ToLower(videoURL), "youtube.com") ||
		strings.Contains(strings.ToLower(videoURL), "youtu.be")
	isPlaylist := strings.Contains(videoURL, "/playlist") || strings.Contains(videoURL, "list=")

	if !isYouTube || isPlaylist {
		multiVideos, err := d.checkMultipleVideos(videoURL)
		if err == nil && len(multiVideos) > 1 {
			return &models.VideoInfo{
				Title:        fmt.Sprintf("Multiple videos (%d)", len(multiVideos)),
				IsMultiVideo: true,
				MultiVideos:  multiVideos,
			}, nil
		}
	}

	return d.getSingleVideoInfo(videoURL)
}

func (d *Downloader) checkMultipleVideos(videoURL string) ([]models.VideoEntry, error) {
	args := []string{"--flat-playlist", "--dump-json", "--no-warnings"}

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
	args := []string{"--dump-json", "--no-playlist", "--no-warnings"}

	isYouTube := strings.Contains(strings.ToLower(videoURL), "youtube.com") ||
		strings.Contains(strings.ToLower(videoURL), "youtu.be")

	if isYouTube {
		if d.cookiesFile != "" {
			args = append(args, "--extractor-args", "youtube:player_client=default,web_safari")
		} else {
			args = append(args, "--extractor-args", "youtube:player_client=web_safari")
		}
	}

	if d.cookiesFile != "" {
		log.Printf("INFO: Using cookies file: %s", d.cookiesFile)
		args = append(args, "--cookies", d.cookiesFile)
	}
	args = append(args, videoURL)

	log.Printf("INFO: Running yt-dlp with args: %v", args)
	cmd := exec.Command("yt-dlp", args...)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("ERROR: yt-dlp error: %v", err)
		return nil, fmt.Errorf("failed to fetch video information")
	}

	var ytdlpInfo models.YtDlpInfo
	if err := json.Unmarshal(output, &ytdlpInfo); err != nil {
		log.Printf("ERROR: JSON parse error: %v, output: %s", err, string(output[:min(500, len(output))]))
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

		// Use consistent quality labels based on height
		quality := getQualityLabel(f.Height)

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
	// 10 minute timeout for downloads
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()
	if err := os.MkdirAll(d.tmpDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	sessionID, err := generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	// Use session ID prefix + title for filename
	outputTemplate := filepath.Join(d.tmpDir, sessionID+"_%(title).80s.%(ext)s")
	args := d.buildDownloadArgs(videoURL, format, outputTemplate, videoIndex)

	// Retry logic with exponential backoff
	var output []byte
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			log.Printf("INFO: Retry attempt %d/%d after %v", attempt+1, maxRetries, backoff)
			time.Sleep(backoff)
		}

		log.Printf("INFO: Running yt-dlp download with args: %v", args)
		cmd := exec.CommandContext(ctx, "yt-dlp", args...)
		output, err = cmd.CombinedOutput()
		if err == nil {
			break
		}

		// Check if error is retryable (network issues, temporary failures)
		outputStr := string(output)
		if strings.Contains(outputStr, "HTTP Error 5") ||
			strings.Contains(outputStr, "timed out") ||
			strings.Contains(outputStr, "Connection reset") {
			log.Printf("WARN: Retryable error on attempt %d: %v", attempt+1, err)
			continue
		}

		// Check if format-specific error - try fallback to best
		if format != "best" && (strings.Contains(outputStr, "format") || strings.Contains(outputStr, "unavailable")) {
			log.Printf("WARN: Format %s failed, trying fallback to best", format)
			fallbackArgs := d.buildDownloadArgs(videoURL, "best", outputTemplate, videoIndex)
			cmd := exec.CommandContext(ctx, "yt-dlp", fallbackArgs...)
			output, err = cmd.CombinedOutput()
			if err == nil {
				break
			}
		}

		// Non-retryable error, fail immediately
		log.Printf("ERROR: yt-dlp download error: %v, output: %s", err, outputStr)
		return nil, fmt.Errorf("download failed or file exceeds size limit")
	}

	if err != nil {
		log.Printf("ERROR: yt-dlp download failed after %d retries: %v, output: %s", maxRetries, err, string(output))
		return nil, fmt.Errorf("download failed or file exceeds size limit")
	}

	files, err := filepath.Glob(filepath.Join(d.tmpDir, sessionID+"_*"))
	if err != nil || len(files) == 0 {
		return nil, fmt.Errorf("downloaded file not found")
	}

	filePath := files[0]
	// Extract just the title part for the download filename (remove session ID prefix)
	baseName := filepath.Base(filePath)
	fileName := strings.TrimPrefix(baseName, sessionID+"_")

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
	isYouTube := strings.Contains(strings.ToLower(videoURL), "youtube.com") ||
		strings.Contains(strings.ToLower(videoURL), "youtu.be")

	formatSpec := "bv*[ext=mp4]+ba[ext=m4a]/b[ext=mp4]/bv*+ba/b"
	if format != "best" {
		formatSpec = fmt.Sprintf("%s+ba/%s", format, format)
	}

	if isInstagram {
		formatSpec = "1/best[vcodec^=avc]/best[ext=mp4]/best"
		log.Printf("INFO: Instagram URL detected, using format 1 (H.264 baseline for iPhone)")
	}

	log.Printf("INFO: Downloading with format: %s", formatSpec)
	args := []string{"-f", formatSpec, "-o", outputTemplate, "--merge-output-format", "mp4", "--no-warnings", "--restrict-filenames"}

	if isYouTube {
		if d.cookiesFile != "" {
			args = append(args, "--extractor-args", "youtube:player_client=default,web_safari")
		} else {
			args = append(args, "--extractor-args", "youtube:player_client=web_safari")
		}
		log.Printf("INFO: YouTube URL detected, using web_safari player client")
	}

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

func (d *Downloader) ExtractAudio(videoURL, audioFormat string, videoIndex int) (*DownloadResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := os.MkdirAll(d.tmpDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	sessionID, err := generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	if audioFormat == "" {
		audioFormat = "mp3"
	}

	allowedFormats := map[string]bool{"mp3": true, "m4a": true, "aac": true, "opus": true, "vorbis": true, "flac": true, "wav": true}
	if !allowedFormats[audioFormat] {
		audioFormat = "mp3"
	}

	outputTemplate := filepath.Join(d.tmpDir, sessionID+"_%(title).80s.%(ext)s")
	args := d.buildAudioArgs(videoURL, audioFormat, outputTemplate, videoIndex)

	log.Printf("INFO: Running yt-dlp audio extraction with args: %v", args)
	cmd := exec.CommandContext(ctx, "yt-dlp", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("ERROR: yt-dlp audio extraction error: %v, output: %s", err, string(output))
		return nil, fmt.Errorf("audio extraction failed")
	}

	files, err := filepath.Glob(filepath.Join(d.tmpDir, sessionID+"_*"))
	if err != nil || len(files) == 0 {
		return nil, fmt.Errorf("extracted audio file not found")
	}

	filePath := files[0]
	baseName := filepath.Base(filePath)
	fileName := strings.TrimPrefix(baseName, sessionID+"_")

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file info: %w", err)
	}

	contentType := getAudioContentType(audioFormat)
	log.Printf("INFO: Audio file size: %d bytes (%.2f MB)", fileInfo.Size(), float64(fileInfo.Size())/(1024*1024))

	return &DownloadResult{
		FilePath:    filePath,
		FileName:    fileName,
		FileSize:    fileInfo.Size(),
		ContentType: contentType,
	}, nil
}

func (d *Downloader) buildAudioArgs(videoURL, audioFormat, outputTemplate string, videoIndex int) []string {
	isYouTube := strings.Contains(strings.ToLower(videoURL), "youtube.com") ||
		strings.Contains(strings.ToLower(videoURL), "youtu.be")

	args := []string{"-x", "--audio-format", audioFormat, "-o", outputTemplate, "--no-warnings", "--restrict-filenames"}

	if isYouTube {
		if d.cookiesFile != "" {
			args = append(args, "--extractor-args", "youtube:player_client=default,web_safari")
		} else {
			args = append(args, "--extractor-args", "youtube:player_client=web_safari")
		}
	}

	if videoIndex > 0 {
		args = append(args, "--playlist-items", fmt.Sprintf("%d", videoIndex))
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

func getAudioContentType(format string) string {
	switch format {
	case "mp3":
		return "audio/mpeg"
	case "m4a", "aac":
		return "audio/mp4"
	case "opus":
		return "audio/opus"
	case "vorbis":
		return "audio/ogg"
	case "flac":
		return "audio/flac"
	case "wav":
		return "audio/wav"
	default:
		return "audio/mpeg"
	}
}

func (d *Downloader) CheckHealth() error {
	if d.healthChecked {
		return d.healthError
	}
	cmd := exec.Command("yt-dlp", "--version")
	d.healthError = cmd.Run()
	d.healthChecked = true
	return d.healthError
}

func generateSessionID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func getQualityLabel(height int) string {
	switch {
	case height >= 2160:
		return "4K"
	case height >= 1440:
		return "1440p"
	case height >= 1080:
		return "1080p"
	case height >= 720:
		return "720p"
	case height >= 480:
		return "480p"
	case height >= 360:
		return "360p"
	case height >= 240:
		return "240p"
	case height >= 144:
		return "144p"
	default:
		return fmt.Sprintf("%dp", height)
	}
}

