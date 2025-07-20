package api

import (
	"search-engine-service/internal/api/handlers"
	"search-engine-service/internal/services"
)

// Handler holds all the API handlers
type Handler struct {
	SearchHandler    *handlers.SearchHandler
	ProviderHandler  *handlers.ProviderHandler
	DashboardHandler *handlers.DashboardHandler
}

// NewHandler creates a new API handler
func NewHandler(searchService *services.SearchService, scoringService *services.ScoringService) *Handler {
	return &Handler{
		SearchHandler:    handlers.NewSearchHandler(searchService, scoringService),
		ProviderHandler:  handlers.NewProviderHandler(searchService),
		DashboardHandler: handlers.NewDashboardHandler(searchService, scoringService),
	}
} 