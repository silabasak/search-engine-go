package handlers

import (
	"context"
	"net/http"

	"search-engine-service/internal/services"

	"github.com/gin-gonic/gin"
)

// ProviderHandler handles provider-related HTTP requests
type ProviderHandler struct {
	searchService *services.SearchService
}

// NewProviderHandler creates a new provider handler
func NewProviderHandler(searchService *services.SearchService) *ProviderHandler {
	return &ProviderHandler{
		searchService: searchService,
	}
}

// GetProviders returns information about all available providers
func (ph *ProviderHandler) GetProviders(c *gin.Context) {
	providers := ph.searchService.GetProviders()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    providers,
	})
}

// RefreshProviders fetches fresh content from all providers
func (ph *ProviderHandler) RefreshProviders(c *gin.Context) {
	ctx := context.Background()

	if err := ph.searchService.RefreshContent(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to refresh content from providers",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Content refreshed successfully",
	})
}

// GetContentStats returns statistics about the content
func (ph *ProviderHandler) GetContentStats(c *gin.Context) {
	stats, err := ph.searchService.GetContentStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get content statistics",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
} 