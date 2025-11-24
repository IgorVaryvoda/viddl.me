package config

import (
	"log"
	"os"
	"strings"
)

type Config struct {
	Port            string
	AllowedOrigins  []string
	AllowedDomains  []string
	MaxDownloadSize string
	CookiesFile     string
	TmpDir          string
}

var defaultDomains = []string{
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

func Load() *Config {
	cfg := &Config{
		Port:            getEnv("PORT", "3000"),
		MaxDownloadSize: getEnv("MAX_DOWNLOAD_SIZE", "2G"),
		CookiesFile:     os.Getenv("YTDLP_COOKIES"),
		TmpDir:          getEnv("TMP_DIR", "./tmp"),
	}

	// Parse allowed origins
	origins := getEnv("ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:8080")
	cfg.AllowedOrigins = strings.Split(origins, ",")

	// Parse allowed domains
	domainsEnv := os.Getenv("ALLOWED_DOMAINS")
	if domainsEnv != "" {
		domains := strings.Split(domainsEnv, ",")
		for i, domain := range domains {
			domains[i] = strings.TrimSpace(domain)
		}
		cfg.AllowedDomains = domains
	} else {
		cfg.AllowedDomains = defaultDomains
	}

	log.Printf("INFO: Configuration loaded - Port: %s, Domains: %v", cfg.Port, cfg.AllowedDomains)
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
