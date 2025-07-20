package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"search-engine-service/internal/database/models"
)

// JSONProvider implements Provider interface for JSON data sources
type JSONProvider struct {
	url     string
	timeout time.Duration
}

// JSONVideoResponse represents the structure of JSON video data
type JSONVideoResponse struct {
	Videos []JSONVideo `json:"videos"`
}

// JSONVideo represents a single video from JSON provider
type JSONVideo struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Views       int       `json:"views"`
	Likes       int       `json:"likes"`
	Duration    int       `json:"duration"`
	Tags        string    `json:"tags"`
	Language    string    `json:"language"`
	PublishedAt time.Time `json:"published_at"`
}

// NewJSONProvider creates a new JSON provider
func NewJSONProvider(url string, timeout time.Duration) *JSONProvider {
	return &JSONProvider{
		url:     url,
		timeout: timeout,
	}
}

// GetName returns the provider name
func (jp *JSONProvider) GetName() string {
	return "json_provider"
}

// GetURL returns the provider URL
func (jp *JSONProvider) GetURL() string {
	return jp.url
}

// GetTimeout returns the provider timeout
func (jp *JSONProvider) GetTimeout() time.Duration {
	return jp.timeout
}

// FetchContent fetches content from the JSON provider
func (jp *JSONProvider) FetchContent(ctx context.Context) ([]models.Content, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: jp.timeout,
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", jp.url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "SearchEngineService/1.0")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var jsonResponse JSONVideoResponse
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Convert to Content models
	var contents []models.Content
	for _, video := range jsonResponse.Videos {
		content := models.Content{
			ProviderID:  video.ID,
			Title:       video.Title,
			Description: video.Description,
			URL:         video.URL,
			Type:        models.ContentTypeVideo,
			Views:       video.Views,
			Likes:       video.Likes,
			Duration:    video.Duration,
			Tags:        video.Tags,
			Language:    video.Language,
			PublishedAt: video.PublishedAt,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		contents = append(contents, content)
	}

	return contents, nil
} 