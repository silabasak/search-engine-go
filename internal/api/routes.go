package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all the API routes
func SetupRoutes(router *gin.Engine, handler *Handler) {
	// API routes group
	api := router.Group("/api")
	{
		// Search routes
		search := api.Group("/search")
		{
			search.GET("", handler.SearchHandler.Search)
			search.POST("/filters", handler.SearchHandler.SearchWithFilters)
		}

		// Content routes
		content := api.Group("/content")
		{
			content.GET("/:id", handler.SearchHandler.GetContentByID)
			content.GET("/popular", handler.SearchHandler.GetPopularContent)
		}

		// Provider routes
		providers := api.Group("/providers")
		{
			providers.GET("", handler.ProviderHandler.GetProviders)
			providers.POST("/refresh", handler.ProviderHandler.RefreshProviders)
			providers.GET("/stats", handler.ProviderHandler.GetContentStats)
		}

		// Dashboard routes
		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("", handler.DashboardHandler.Dashboard)
			dashboard.GET("/search", handler.DashboardHandler.DashboardSearch)
			dashboard.GET("/stats", handler.DashboardHandler.GetDashboardStats)
		}
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "search-engine-service",
		})
	})

	// Serve static files
	router.Static("/static", "./web/static")
	
	// Serve dashboard HTML
	router.GET("/dashboard", func(c *gin.Context) {
		c.File("./web/templates/dashboard.html")
	})

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Search Engine Service API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"search":     "/api/search",
				"content":    "/api/content",
				"providers":  "/api/providers",
				"dashboard":  "/api/dashboard",
				"health":     "/health",
			},
			"dashboard": "/dashboard",
		})
	})
} 