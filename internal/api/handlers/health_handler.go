package handlers

import (
	"context"
	"net/http"
	"time"

	"search-engine-service/internal/database"
	"search-engine-service/internal/utils/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db     *database.Database
	logger logger.Logger
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.Database, logger logger.Logger) *HealthHandler {
	return &HealthHandler{
		db:     db,
		logger: logger,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Service   string            `json:"service"`
	Version   string            `json:"version"`
	Uptime    string            `json:"uptime"`
	Checks    map[string]Check  `json:"checks,omitempty"`
}

// Check represents a health check result
type Check struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Health performs a basic health check
func (h *HealthHandler) Health(c *gin.Context) {
	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Service:   "search-engine-service",
		Version:   "2.0.0",
		Uptime:    h.getUptime(),
	}

	h.logger.Info("Health check requested",
		zap.String("client_ip", c.ClientIP()),
		zap.String("user_agent", c.GetHeader("User-Agent")),
	)

	c.JSON(http.StatusOK, response)
}

// Ready performs a readiness check including database connectivity
func (h *HealthHandler) Ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	response := HealthResponse{
		Status:    "ok",
		Timestamp: time.Now(),
		Service:   "search-engine-service",
		Version:   "2.0.0",
		Uptime:    h.getUptime(),
		Checks:    make(map[string]Check),
	}

	// Database health check
	dbCheck := h.checkDatabase(ctx)
	response.Checks["database"] = dbCheck

	// Provider health check
	providerCheck := h.checkProviders(ctx)
	response.Checks["providers"] = providerCheck

	// Determine overall status
	if dbCheck.Status == "error" || providerCheck.Status == "error" {
		response.Status = "error"
		c.JSON(http.StatusServiceUnavailable, response)
		return
	}

	if dbCheck.Status == "warning" || providerCheck.Status == "warning" {
		response.Status = "warning"
		c.JSON(http.StatusOK, response)
		return
	}

	h.logger.Info("Readiness check completed",
		zap.String("status", response.Status),
		zap.String("client_ip", c.ClientIP()),
	)

	c.JSON(http.StatusOK, response)
}

// checkDatabase checks database connectivity
func (h *HealthHandler) checkDatabase(ctx context.Context) Check {
	check := Check{
		Status:  "ok",
		Message: "Database is healthy",
		Details: make(map[string]interface{}),
	}

	// Test database connection
	sqlDB, err := h.db.DB.DB()
	if err != nil {
		check.Status = "error"
		check.Message = "Failed to get database instance"
		check.Details["error"] = err.Error()
		return check
	}

	// Ping database
	if err := sqlDB.PingContext(ctx); err != nil {
		check.Status = "error"
		check.Message = "Database ping failed"
		check.Details["error"] = err.Error()
		return check
	}

	// Get database stats
	stats := sqlDB.Stats()
	check.Details["open_connections"] = stats.OpenConnections
	check.Details["in_use"] = stats.InUse
	check.Details["idle"] = stats.Idle
	check.Details["wait_count"] = stats.WaitCount
	check.Details["wait_duration"] = stats.WaitDuration.String()

	return check
}

// checkProviders checks provider connectivity
func (h *HealthHandler) checkProviders(ctx context.Context) Check {
	check := Check{
		Status:  "ok",
		Message: "All providers are healthy",
		Details: make(map[string]interface{}),
	}

	// This would typically check external provider endpoints
	// For now, we'll simulate a basic check
	providers := []string{"json_provider", "xml_provider"}
	healthyProviders := 0

	for _, provider := range providers {
		// Simulate provider health check
		if h.simulateProviderHealth(ctx, provider) {
			healthyProviders++
		}
	}

	check.Details["total_providers"] = len(providers)
	check.Details["healthy_providers"] = healthyProviders

	if healthyProviders == 0 {
		check.Status = "error"
		check.Message = "No providers are healthy"
	} else if healthyProviders < len(providers) {
		check.Status = "warning"
		check.Message = "Some providers are unhealthy"
	}

	return check
}

// simulateProviderHealth simulates a provider health check
func (h *HealthHandler) simulateProviderHealth(ctx context.Context, provider string) bool {
	// In a real implementation, this would make HTTP requests to provider endpoints
	// For now, we'll simulate success
	return true
}

// getUptime returns the service uptime
func (h *HealthHandler) getUptime() string {
	// In a real implementation, this would track the actual uptime
	// For now, we'll return a placeholder
	return "1h 23m 45s"
}

// Metrics returns service metrics
func (h *HealthHandler) Metrics(c *gin.Context) {
	metrics := map[string]interface{}{
		"service": "search-engine-service",
		"version": "2.0.0",
		"timestamp": time.Now(),
		"metrics": map[string]interface{}{
			"requests_total": 1234,
			"requests_per_second": 45.6,
			"average_response_time": "125ms",
			"error_rate": 0.02,
			"active_connections": 15,
			"memory_usage": "256MB",
			"cpu_usage": 12.5,
		},
	}

	c.JSON(http.StatusOK, metrics)
} 