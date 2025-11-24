package cleanup

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Cleaner struct {
	tmpDir   string
	interval time.Duration
	maxAge   time.Duration
}

func New(tmpDir string, interval, maxAge time.Duration) *Cleaner {
	return &Cleaner{
		tmpDir:   tmpDir,
		interval: interval,
		maxAge:   maxAge,
	}
}

func (c *Cleaner) Start() {
	c.cleanOnce()

	ticker := time.NewTicker(c.interval)
	go func() {
		for range ticker.C {
			c.cleanOnce()
		}
	}()
}

func (c *Cleaner) cleanOnce() {
	files, err := os.ReadDir(c.tmpDir)
	if err != nil {
		log.Printf("ERROR: Error reading tmp directory: %v", err)
		return
	}

	now := time.Now()
	cleanedCount := 0
	var totalSize int64

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		filePath := filepath.Join(c.tmpDir, file.Name())
		info, err := os.Stat(filePath)
		if err != nil {
			log.Printf("ERROR: Error stating file %s: %v", file.Name(), err)
			continue
		}

		if now.Sub(info.ModTime()) > c.maxAge {
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

func ScheduleFileRemoval(filePath string, delay time.Duration) {
	go func() {
		time.Sleep(delay)
		if err := os.Remove(filePath); err != nil {
			log.Printf("ERROR: Failed to cleanup file %s: %v", filePath, err)
		} else {
			log.Printf("INFO: Cleaned up downloaded file: %s", filePath)
		}
	}()
}
