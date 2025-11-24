package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"viddl.me/backend/internal/cleanup"
	"viddl.me/backend/internal/config"
	"viddl.me/backend/internal/handlers"
	"viddl.me/backend/internal/middleware"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetPrefix("[viddl.me] ")
}

func main() {
	cfg := config.Load()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Use(middleware.SecurityHeaders())
	r.Use(middleware.Gzip())

	limiter := middleware.NewIPRateLimiter(rate.Every(time.Minute/3), 3)

	h := handlers.New(cfg)

	r.POST("/api/info", middleware.RateLimit(limiter), h.GetVideoInfo)
	r.POST("/api/download", middleware.RateLimit(limiter), h.DownloadVideo)
	r.GET("/health", h.HealthCheck)

	cleaner := cleanup.New(cfg.TmpDir, 5*time.Minute, 5*time.Minute)
	cleaner.Start()

	log.Printf("INFO: Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("FATAL: Server failed to start: %v", err)
	}
}
