package handlers

import (
	"net/http"
	"strconv"

	"search-engine-service/internal/database/models"
	"search-engine-service/internal/services"

	"github.com/gin-gonic/gin"
)

// DashboardHandler handles dashboard-related HTTP requests
type DashboardHandler struct {
	searchService  *services.SearchService
	scoringService *services.ScoringService
}

// NewDashboardHandler creates a new dashboard handler
func NewDashboardHandler(searchService *services.SearchService, scoringService *services.ScoringService) *DashboardHandler {
	return &DashboardHandler{
		searchService:  searchService,
		scoringService: scoringService,
	}
}

// Dashboard handles the main dashboard page
func (dh *DashboardHandler) Dashboard(c *gin.Context) {
	// Get popular content for dashboard
	contents, err := dh.searchService.GetPopularContent(10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to load dashboard data",
		})
		return
	}

	// Get content statistics
	stats, err := dh.searchService.GetContentStats()
	if err != nil {
		stats = map[string]interface{}{
			"total_content": 0,
			"video_count":   0,
			"text_count":    0,
		}
	}

	// Get provider information
	providers := dh.searchService.GetProviders()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"popular_content": contents,
			"statistics":      stats,
			"providers":       providers,
		},
	})
}

// DashboardSearch handles search requests from the dashboard
func (dh *DashboardHandler) DashboardSearch(c *gin.Context) {
	query := c.Query("q")
	contentType := c.Query("type")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 50 {
		limit = 10
	}

	// Use the search service to perform the search
	result, err := dh.searchService.Search(query, models.ContentType(contentType), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to perform search",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    result,
	})
}

// GetDashboardStats returns detailed statistics for the dashboard
func (dh *DashboardHandler) GetDashboardStats(c *gin.Context) {
	// Get basic stats
	stats, err := dh.searchService.GetContentStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get statistics",
		})
		return
	}

	// Get provider information
	providers := dh.searchService.GetProviders()

	// Get top content by score
	topContent, err := dh.searchService.GetPopularContent(5)
	if err != nil {
		topContent = []models.Content{}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"statistics":  stats,
			"providers":   providers,
			"top_content": topContent,
		},
	})
} 