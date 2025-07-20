package di

import (
	"context"
	"time"

	"search-engine-service/internal/api/handlers"
	"search-engine-service/internal/api/middleware"
	"search-engine-service/internal/config"
	"search-engine-service/internal/database"
	"search-engine-service/internal/database/repository"
	"search-engine-service/internal/providers"
	"search-engine-service/internal/services"
	"search-engine-service/internal/utils/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

// Container represents the dependency injection container
type Container struct {
	container *dig.Container
}

// NewContainer creates a new DI container
func NewContainer() *Container {
	container := dig.New()
	
	// Register all dependencies
	registerDependencies(container)
	
	return &Container{
		container: container,
	}
}

// registerDependencies registers all dependencies in the container
func registerDependencies(container *dig.Container) {
	// Configuration
	container.Provide(config.LoadConfig)
	
	// Logger
	container.Provide(logger.NewLogger)
	
	// Database
	container.Provide(database.NewDatabase)
	
	// Repositories
	container.Provide(repository.NewContentRepository)
	
	// Providers
	container.Provide(providers.NewProviderManager)
	container.Provide(providers.NewJSONProvider)
	container.Provide(providers.NewXMLProvider)
	
	// Services
	container.Provide(services.NewScoringService)
	container.Provide(services.NewSearchService)
	
	// Middleware
	container.Provide(middleware.NewSecurityMiddleware)
	container.Provide(middleware.NewRateLimiter)
	container.Provide(middleware.NewRequestLogger)
	
	// Handlers
	container.Provide(handlers.NewSearchHandler)
	container.Provide(handlers.NewProviderHandler)
	container.Provide(handlers.NewDashboardHandler)
	container.Provide(handlers.NewHealthHandler)
	
	// Router
	container.Provide(NewRouter)
}

// NewRouter creates a new router with all dependencies
func NewRouter(
	searchHandler *handlers.SearchHandler,
	providerHandler *handlers.ProviderHandler,
	dashboardHandler *handlers.DashboardHandler,
	healthHandler *handlers.HealthHandler,
	securityMiddleware *middleware.SecurityMiddleware,
	rateLimiter *middleware.RateLimiter,
	requestLogger *middleware.RequestLogger,
) *gin.Engine {
	router := gin.New()
	
	// Add middleware
	router.Use(securityMiddleware.SecurityHeaders())
	router.Use(securityMiddleware.CORS())
	router.Use(rateLimiter.Limit())
	router.Use(requestLogger.Log())
	router.Use(gin.Recovery())
	
	// Setup routes
	setupRoutes(router, searchHandler, providerHandler, dashboardHandler, healthHandler)
	
	return router
}

// setupRoutes configures all routes
func setupRoutes(
	router *gin.Engine,
	searchHandler *handlers.SearchHandler,
	providerHandler *handlers.ProviderHandler,
	dashboardHandler *handlers.DashboardHandler,
	healthHandler *handlers.HealthHandler,
) {
	// API routes
	api := router.Group("/api/v1")
	{
		// Search routes
		search := api.Group("/search")
		{
			search.GET("", searchHandler.Search)
			search.POST("/filters", searchHandler.SearchWithFilters)
			search.GET("/suggestions", searchHandler.GetSuggestions)
		}
		
		// Content routes
		content := api.Group("/content")
		{
			content.GET("/:id", searchHandler.GetContentByID)
			content.GET("/popular", searchHandler.GetPopularContent)
			content.GET("/trending", searchHandler.GetTrendingContent)
		}
		
		// Provider routes
		providers := api.Group("/providers")
		{
			providers.GET("", providerHandler.GetProviders)
			providers.POST("/refresh", providerHandler.RefreshProviders)
			providers.GET("/stats", providerHandler.GetContentStats)
			providers.GET("/health", providerHandler.GetProviderHealth)
		}
		
		// Analytics routes
		analytics := api.Group("/analytics")
		{
			analytics.GET("/stats", dashboardHandler.GetAnalytics)
			analytics.GET("/trends", dashboardHandler.GetTrends)
		}
	}
	
	// Health check
	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)
	
	// Static files
	router.Static("/static", "./web/static")
	
	// Dashboard
	router.GET("/dashboard", dashboardHandler.Dashboard)
	
	// Root
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Search Engine Service API",
			"version": "2.0.0",
			"docs": "/api/v1/docs",
		})
	})
}

// Resolve resolves a dependency from the container
func (c *Container) Resolve(constructor interface{}) error {
	return c.container.Invoke(constructor)
}

// MustResolve resolves a dependency and panics on error
func (c *Container) MustResolve(constructor interface{}) {
	if err := c.Resolve(constructor); err != nil {
		panic(err)
	}
}

// Shutdown gracefully shuts down the container
func (c *Container) Shutdown(ctx context.Context) error {
	// Add any cleanup logic here
	return nil
} 