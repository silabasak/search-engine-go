package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"search-engine-service/internal/api/handlers"
	"search-engine-service/internal/config"
	"search-engine-service/internal/database"
	"search-engine-service/internal/services"
	"search-engine-service/internal/utils/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestEnvironment sets up the test environment
func setupTestEnvironment(t *testing.T) (*gin.Engine, func()) {
	// Set test mode
	gin.SetMode(gin.TestMode)
	
	// Load test config
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: "8080",
		},
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     "3306",
			User:     "test_user",
			Password: "test_password",
			Name:     "test_db",
		},
		Environment: "test",
	}
	
	// Initialize logger
	log, err := logger.NewLogger(cfg)
	require.NoError(t, err)
	
	// Initialize database (use test database)
	db, err := database.NewDatabase(cfg)
	require.NoError(t, err)
	
	// Initialize services
	scoringService := services.NewScoringService()
	searchService := services.NewSearchService(db.DB, nil) // Provider manager is nil for tests
	
	// Initialize handlers
	searchHandler := handlers.NewSearchHandler(searchService, scoringService)
	healthHandler := handlers.NewHealthHandler(db, log)
	
	// Setup router
	router := gin.New()
	router.Use(gin.Recovery())
	
	// Setup routes
	api := router.Group("/api/v1")
	{
		search := api.Group("/search")
		{
			search.GET("", searchHandler.Search)
		}
		
		content := api.Group("/content")
		{
			content.GET("/:id", searchHandler.GetContentByID)
			content.GET("/popular", searchHandler.GetPopularContent)
		}
	}
	
	router.GET("/health", healthHandler.Health)
	
	// Cleanup function
	cleanup := func() {
		// Clean up test data
		db.DB.Exec("DELETE FROM contents")
	}
	
	return router, cleanup
}

// TestHealthEndpoint tests the health endpoint
func TestHealthEndpoint(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()
	
	// Create request
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)
	
	// Create response recorder
	w := httptest.NewRecorder()
	
	// Serve request
	router.ServeHTTP(w, req)
	
	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "ok", response["status"])
	assert.Equal(t, "search-engine-service", response["service"])
	assert.Equal(t, "2.0.0", response["version"])
}

// TestSearchEndpoint tests the search endpoint
func TestSearchEndpoint(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()
	
	// Create request
	req, err := http.NewRequest("GET", "/api/v1/search?q=golang&type=video&page=1&limit=10", nil)
	require.NoError(t, err)
	
	// Create response recorder
	w := httptest.NewRecorder()
	
	// Serve request
	router.ServeHTTP(w, req)
	
	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.True(t, response["success"].(bool))
	assert.NotNil(t, response["data"])
}

// TestSearchEndpointInvalidParams tests search with invalid parameters
func TestSearchEndpointInvalidParams(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()
	
	testCases := []struct {
		name     string
		query    string
		expected int
	}{
		{
			name:     "Invalid page number",
			query:    "/api/v1/search?q=test&page=invalid",
			expected: http.StatusOK, // Should default to page 1
		},
		{
			name:     "Invalid limit",
			query:    "/api/v1/search?q=test&limit=1000", // Too high
			expected: http.StatusOK, // Should default to max limit
		},
		{
			name:     "Invalid content type",
			query:    "/api/v1/search?q=test&type=invalid",
			expected: http.StatusBadRequest,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tc.query, nil)
			require.NoError(t, err)
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			
			assert.Equal(t, tc.expected, w.Code)
		})
	}
}

// TestContentEndpoints tests content-related endpoints
func TestContentEndpoints(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()
	
	// Test popular content endpoint
	t.Run("Popular Content", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/content/popular?limit=5", nil)
		require.NoError(t, err)
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		assert.Equal(t, http.StatusOK, w.Code)
		
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)
		
		assert.True(t, response["success"].(bool))
	})
	
	// Test content by ID endpoint
	t.Run("Content By ID", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/v1/content/999", nil)
		require.NoError(t, err)
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		// Should return 404 for non-existent content
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

// TestSearchWithFilters tests search with filters
func TestSearchWithFilters(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()
	
	// Test search with filters
	filterData := map[string]interface{}{
		"query":       "programming",
		"content_type": "video",
		"min_score":   50.0,
		"max_score":   100.0,
		"page":        1,
		"limit":       10,
	}
	
	jsonData, err := json.Marshal(filterData)
	require.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/api/v1/search/filters", bytes.NewBuffer(jsonData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.True(t, response["success"].(bool))
}

// TestRateLimiting tests rate limiting functionality
func TestRateLimiting(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()
	
	// Make multiple requests quickly
	for i := 0; i < 105; i++ { // Exceed rate limit
		req, err := http.NewRequest("GET", "/health", nil)
		require.NoError(t, err)
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		if i < 100 {
			assert.Equal(t, http.StatusOK, w.Code)
		} else {
			assert.Equal(t, http.StatusTooManyRequests, w.Code)
		}
	}
}

// TestErrorHandling tests error handling
func TestErrorHandling(t *testing.T) {
	router, cleanup := setupTestEnvironment(t)
	defer cleanup()
	
	// Test malformed JSON
	req, err := http.NewRequest("POST", "/api/v1/search/filters", bytes.NewBufferString("invalid json"))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// BenchmarkSearchEndpoint benchmarks the search endpoint
func BenchmarkSearchEndpoint(b *testing.B) {
	router, cleanup := setupTestEnvironment(&testing.T{})
	defer cleanup()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest("GET", "/api/v1/search?q=test&page=1&limit=10", nil)
		if err != nil {
			b.Fatal(err)
		}
		
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		
		if w.Code != http.StatusOK {
			b.Fatalf("Expected status 200, got %d", w.Code)
		}
	}
} 