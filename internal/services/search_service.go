package services

import (
	"context"
	"log"
	"time"

	"search-engine-service/internal/database/models"
	"search-engine-service/internal/database/repository"
	"search-engine-service/internal/providers"

	"gorm.io/gorm"
)

// SearchService handles search operations and content management
type SearchService struct {
	contentRepo     models.ContentRepository
	providerManager *providers.ProviderManager
	scoringService  *ScoringService
}

// NewSearchService creates a new search service
func NewSearchService(db *gorm.DB, providerManager *providers.ProviderManager) *SearchService {
	contentRepo := repository.NewContentRepository(db)
	scoringService := NewScoringService()
	
	return &SearchService{
		contentRepo:     contentRepo,
		providerManager: providerManager,
		scoringService:  scoringService,
	}
}

// Search performs a search operation with the given parameters
func (ss *SearchService) Search(query string, contentType models.ContentType, page, limit int) (*models.SearchResult, error) {
	// Validate parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Perform search
	result, err := ss.contentRepo.Search(query, contentType, page, limit)
	if err != nil {
		return nil, err
	}

	// Calculate scores for all results
	for i := range result.Contents {
		ss.scoringService.CalculateScore(&result.Contents[i])
	}

	return result, nil
}

// GetContentByID retrieves a specific content by ID
func (ss *SearchService) GetContentByID(id uint) (*models.Content, error) {
	content, err := ss.contentRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Calculate score
	ss.scoringService.CalculateScore(content)

	return content, nil
}

// GetPopularContent retrieves popular content
func (ss *SearchService) GetPopularContent(limit int) ([]models.Content, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}

	contents, err := ss.contentRepo.GetPopular(limit)
	if err != nil {
		return nil, err
	}

	// Calculate scores
	for i := range contents {
		ss.scoringService.CalculateScore(&contents[i])
	}

	return contents, nil
}

// RefreshContent fetches fresh content from all providers and updates the database
func (ss *SearchService) RefreshContent(ctx context.Context) error {
	log.Println("Starting content refresh...")

	// Fetch content from all providers
	contents, err := ss.providerManager.FetchAllContent(ctx)
	if err != nil {
		return err
	}

	log.Printf("Fetched %d content items from providers", len(contents))

	// Calculate scores for all content
	ss.scoringService.CalculateScoresForBatch(contents)

	// Bulk upsert to database
	if err := ss.contentRepo.BulkUpsert(contents); err != nil {
		return err
	}

	log.Printf("Successfully updated %d content items", len(contents))
	return nil
}

// GetProviders returns information about all available providers
func (ss *SearchService) GetProviders() []map[string]interface{} {
	providers := ss.providerManager.GetAllProviders()
	var result []map[string]interface{}

	for _, provider := range providers {
		result = append(result, map[string]interface{}{
			"name":    provider.GetName(),
			"url":     provider.GetURL(),
			"timeout": provider.GetTimeout().String(),
		})
	}

	return result
}

// GetContentStats returns statistics about the content
func (ss *SearchService) GetContentStats() (map[string]interface{}, error) {
	// This is a simplified implementation
	// In a real system, you might want to add more sophisticated statistics
	
	stats := map[string]interface{}{
		"total_content": 0,
		"video_count":   0,
		"text_count":    0,
		"last_updated":  time.Now(),
	}

	return stats, nil
}

// SearchWithFilters performs a search with additional filters
func (ss *SearchService) SearchWithFilters(query string, filters map[string]interface{}, page, limit int) (*models.SearchResult, error) {
	// Extract content type from filters
	contentType := models.ContentType("")
	if typeStr, ok := filters["type"].(string); ok {
		contentType = models.ContentType(typeStr)
	}

	// Perform basic search
	result, err := ss.Search(query, contentType, page, limit)
	if err != nil {
		return nil, err
	}

	// Apply additional filters if needed
	// This is where you could add more sophisticated filtering logic

	return result, nil
}

// AutoRefresh starts a background goroutine to periodically refresh content
func (ss *SearchService) AutoRefresh(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Auto refresh stopped")
			return
		case <-ticker.C:
			if err := ss.RefreshContent(ctx); err != nil {
				log.Printf("Auto refresh failed: %v", err)
			}
		}
	}
} 