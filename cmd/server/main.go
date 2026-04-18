package main

import (
	"context"
	"crypto/tls"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zmashaly/cp-leaderboard/internal/database"
	"github.com/zmashaly/cp-leaderboard/internal/models"
	"github.com/zmashaly/cp-leaderboard/internal/repository"
	"github.com/zmashaly/cp-leaderboard/internal/service"
	"github.com/zmashaly/cp-leaderboard/internal/sync"
)

func main() {
	// ── 1. Logger ─────────────────────────────────────────────────────────
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// ── 2. Database ────────────────────────────────────────────────────────
	db, err := database.Connect(database.Config{})
	if err != nil {
		log.Fatalf("❌ DB connection failed: %v", err)
	}
	if err := database.AutoMigrate(db, &models.User{}, &models.DailySnapshot{}); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	// ── 3. Redis ───────────────────────────────────────────────────────────
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("❌ REDIS_URL environment variable is required")
	}

	tlsCfg := (*tls.Config)(nil)
	if os.Getenv("REDIS_TLS") == "true" {
		tlsCfg = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:      redisURL,
		Password:  os.Getenv("REDIS_PASSWORD"),
		TLSConfig: tlsCfg,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("❌ Redis connection failed: %v", err)
	}
	logger.Info("✅ Redis connected")

	// ── 4. Repositories & Services ─────────────────────────────────────────
	userRepo := repository.NewUserRepository(db)
	snapRepo := repository.NewSnapshotRepository(db)
	leaderboardSvc := service.NewLeaderboardService(rdb)

	// ── 5. Sync Engine (background goroutine) ──────────────────────────────
	syncCtx, syncCancel := context.WithCancel(context.Background())
	defer syncCancel()

	engine := sync.New(userRepo, snapRepo, rdb, logger)
	go func() {
		for {
			select {
			case <-syncCtx.Done():
				logger.Info("🛑 Sync engine stopped")
				return
			default:
				engine.Run(syncCtx)
				select {
				case <-time.After(5 * time.Minute):
				case <-syncCtx.Done():
					return
				}
			}
		}
	}()

	// ── 6. HTTP Server ─────────────────────────────────────────────────────
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:5173"
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{allowedOrigin},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// ── 7. Routes ──────────────────────────────────────────────────────────

	// Root endpoint for Hugging Face Readiness Probe
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "Running",
			"message": "Leaderboard API is UP and actively syncing!",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":       "healthy",
			"version":      "1.0.0",
			"developed_by": "Ziad Mashaly",
			"position":     "ICPC Delta University Community President",
			"message":      "Made with love for the competitive programming community",
		})
	})

	// Leaderboard endpoint
	r.GET("/api/leaderboard", func(c *gin.Context) {
		users, err := leaderboardSvc.GetLeaderboard(c.Request.Context(), 0)
		if err != nil {
			logger.Error("Failed to fetch leaderboard", "err", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to fetch data from server",
			})
			return
		}
		c.JSON(http.StatusOK, users)
	})

	// ── 8. Graceful Shutdown ───────────────────────────────────────────────
	port := os.Getenv("PORT")
	if port == "" {
		port = "7860" // ✅ Hugging Face standard port
	}

	srv := &http.Server{
		Addr:         "0.0.0.0:" + port, // ✅ Must explicitly bind to 0.0.0.0
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("🚀 Server started", "port", port, "host", "0.0.0.0")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("⏳ Shutting down gracefully...")
	syncCancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("❌ Server forced to shutdown", "err", err)
	}
	logger.Info("✅ Server stopped")
}
