package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	AllowedOrigins  []string
	AllowedDomains  []string
	MaxDownloadSize string
	CookiesFile     string
	TmpDir          string
	APIKey          string
}

var defaultOrigins = []string{
	"http://localhost:5173",
	"http://localhost:8080",
	"https://sirv-ai-tools.vercel.app",
	"https://*.sirv.com",
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
	"threads.net",
	"sirv.com",
	"v3.fal.media",
	"fal.media",
}

func Load() *Config {
	godotenv.Load()

	cfg := &Config{
		Port:            getEnv("PORT", "3000"),
		MaxDownloadSize: getEnv("MAX_DOWNLOAD_SIZE", "2G"),
		CookiesFile:     os.Getenv("YTDLP_COOKIES"),
		TmpDir:          getEnv("TMP_DIR", "./tmp"),
		APIKey:          os.Getenv("API_KEY"),
	}

	// Parse allowed origins (env var adds to defaults)
	cfg.AllowedOrigins = append([]string{}, defaultOrigins...)
	if originsEnv := os.Getenv("ALLOWED_ORIGINS"); originsEnv != "" {
		for _, origin := range strings.Split(originsEnv, ",") {
			origin = strings.TrimSpace(origin)
			// Add https:// if no protocol specified
			if origin != "*" && !strings.HasPrefix(origin, "http://") && !strings.HasPrefix(origin, "https://") {
				origin = "https://" + origin
			}
			cfg.AllowedOrigins = append(cfg.AllowedOrigins, origin)
		}
	}

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
