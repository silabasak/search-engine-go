package providers

import (
	"context"
	"time"

	"search-engine-service/internal/database/models"
)

// Provider interface defines the contract for content providers
type Provider interface {
	GetName() string
	FetchContent(ctx context.Context) ([]models.Content, error)
	GetURL() string
	GetTimeout() time.Duration
}

// ProviderManager manages multiple providers
type ProviderManager struct {
	providers []Provider
	config    ProviderConfig
}

// ProviderConfig holds provider configuration
type ProviderConfig struct {
	JSONURL   string
	XMLURL    string
	Timeout   time.Duration
	RateLimit int
}

// NewManager creates a new provider manager
func NewManager(config ProviderConfig) *ProviderManager {
	manager := &ProviderManager{
		config: config,
	}

	// Initialize providers
	jsonProvider := NewJSONProvider(config.JSONURL, config.Timeout)
	xmlProvider := NewXMLProvider(config.XMLURL, config.Timeout)

	manager.providers = []Provider{
		jsonProvider,
		xmlProvider,
	}

	return manager
}

// GetAllProviders returns all registered providers
func (pm *ProviderManager) GetAllProviders() []Provider {
	return pm.providers
}

// FetchAllContent fetches content from all providers
func (pm *ProviderManager) FetchAllContent(ctx context.Context) ([]models.Content, error) {
	var allContent []models.Content

	for _, provider := range pm.providers {
		content, err := provider.FetchContent(ctx)
		if err != nil {
			// Log error but continue with other providers
			continue
		}

		// Set provider name for each content
		for i := range content {
			content[i].Provider = provider.GetName()
		}

		allContent = append(allContent, content...)
	}

	return allContent, nil
}

// RefreshContent fetches fresh content from all providers
func (pm *ProviderManager) RefreshContent(ctx context.Context) ([]models.Content, error) {
	return pm.FetchAllContent(ctx)
}

// GetProviderByName returns a specific provider by name
func (pm *ProviderManager) GetProviderByName(name string) Provider {
	for _, provider := range pm.providers {
		if provider.GetName() == name {
			return provider
		}
	}
	return nil
} 