package handlers

import (
	"net/http"
	"strconv"

	"search-engine-service/internal/database/models"
	"search-engine-service/internal/services"

	"github.com/gin-gonic/gin"
)

// SearchHandler handles search-related HTTP requests
type SearchHandler struct {
	searchService  *services.SearchService
	scoringService *services.ScoringService
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(searchService *services.SearchService, scoringService *services.ScoringService) *SearchHandler {
	return &SearchHandler{
		searchService:  searchService,
		scoringService: scoringService,
	}
}

// Search handles search requests
func (sh *SearchHandler) Search(c *gin.Context) {
	// Get query parameters
	query := c.Query("q")
	contentType := models.ContentType(c.Query("type"))
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	// Parse pagination parameters
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	// Validate content type
	if contentType != "" && contentType != models.ContentTypeVideo && contentType != models.ContentTypeText && contentType != "all" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid content type. Must be 'video', 'text', or 'all'",
		})
		return
	}

	// Validate query length
	if len(query) > 500 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Query too long. Maximum 500 characters allowed.",
		})
		return
	}

	// Perform search
	result, err := sh.searchService.Search(query, contentType, page, limit)
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

// GetContentByID handles requests to get a specific content by ID
func (sh *SearchHandler) GetContentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid content ID",
		})
		return
	}

	content, err := sh.searchService.GetContentByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Content not found",
		})
		return
	}

	// Get score breakdown
	scoreBreakdown := sh.scoringService.GetScoreBreakdown(content)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"content":        content,
			"score_breakdown": scoreBreakdown,
		},
	})
}

// GetPopularContent handles requests to get popular content
func (sh *SearchHandler) GetPopularContent(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	contents, err := sh.searchService.GetPopularContent(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get popular content",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    contents,
	})
}

// SearchWithFilters handles advanced search with filters
func (sh *SearchHandler) SearchWithFilters(c *gin.Context) {
	var request struct {
		Query   string                 `json:"query"`
		Filters map[string]interface{} `json:"filters"`
		Page    int                    `json:"page"`
		Limit   int                    `json:"limit"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	// Validate parameters
	if request.Page < 1 {
		request.Page = 1
	}
	if request.Limit < 1 || request.Limit > 100 {
		request.Limit = 10
	}

	// Perform search with filters
	result, err := sh.searchService.SearchWithFilters(request.Query, request.Filters, request.Page, request.Limit)
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