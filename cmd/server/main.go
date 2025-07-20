package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"search-engine-service/internal/api"
	"search-engine-service/internal/api/middleware"
	"search-engine-service/internal/config"
	"search-engine-service/internal/database"
	"search-engine-service/internal/providers"
	"search-engine-service/internal/services"
	"search-engine-service/internal/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger, err := logger.NewLogger(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Initialize database
	db, err := database.Init(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize providers
	providerConfig := providers.ProviderConfig{
		JSONURL:   cfg.Providers.JSONURL,
		XMLURL:    cfg.Providers.XMLURL,
		Timeout:   cfg.Providers.Timeout,
		RateLimit: cfg.Providers.RateLimit,
	}
	providerManager := providers.NewManager(providerConfig)
	
	// Initialize services
	searchService := services.NewSearchService(db, providerManager)
	scoringService := services.NewScoringService()

	// Initialize API
	apiHandler := api.NewHandler(searchService, scoringService)

	// Setup router
	router := gin.Default()
	
	// Initialize middleware
	securityMiddleware := middleware.NewSecurityMiddleware(logger)
	rateLimiter := middleware.NewRateLimiter(logger, 100)
	requestLogger := middleware.NewRequestLogger(logger)
	
	// Add security middleware
	router.Use(securityMiddleware.SecurityHeaders())
	router.Use(securityMiddleware.CORS())
	router.Use(securityMiddleware.RequestID())
	router.Use(securityMiddleware.InputSanitizer())
	router.Use(rateLimiter.Limit())
	
	// Add standard middleware
	router.Use(requestLogger.Log())
	router.Use(gin.Recovery())
	
	// Setup routes
	api.SetupRoutes(router, apiHandler)

	// Setup server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
} 