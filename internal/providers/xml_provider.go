package providers

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"

	"search-engine-service/internal/database/models"
)

// XMLProvider implements Provider interface for XML data sources
type XMLProvider struct {
	url     string
	timeout time.Duration
}

// XMLArticleResponse represents the structure of XML article data
type XMLArticleResponse struct {
	XMLName xml.Name     `xml:"articles"`
	Articles []XMLArticle `xml:"article"`
}

// XMLArticle represents a single article from XML provider
type XMLArticle struct {
	ID          string    `xml:"id"`
	Title       string    `xml:"title"`
	Description string    `xml:"description"`
	URL         string    `xml:"url"`
	ReadingTime int       `xml:"reading_time"`
	Reactions   int       `xml:"reactions"`
	Tags        string    `xml:"tags"`
	Language    string    `xml:"language"`
	PublishedAt time.Time `xml:"published_at"`
}

// NewXMLProvider creates a new XML provider
func NewXMLProvider(url string, timeout time.Duration) *XMLProvider {
	return &XMLProvider{
		url:     url,
		timeout: timeout,
	}
}

// GetName returns the provider name
func (xp *XMLProvider) GetName() string {
	return "xml_provider"
}

// GetURL returns the provider URL
func (xp *XMLProvider) GetURL() string {
	return xp.url
}

// GetTimeout returns the provider timeout
func (xp *XMLProvider) GetTimeout() time.Duration {
	return xp.timeout
}

// FetchContent fetches content from the XML provider
func (xp *XMLProvider) FetchContent(ctx context.Context) ([]models.Content, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: xp.timeout,
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", xp.url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/xml")
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

	// Parse XML response
	var xmlResponse XMLArticleResponse
	if err := xml.Unmarshal(body, &xmlResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML: %w", err)
	}

	// Convert to Content models
	var contents []models.Content
	for _, article := range xmlResponse.Articles {
		content := models.Content{
			ProviderID:  article.ID,
			Title:       article.Title,
			Description: article.Description,
			URL:         article.URL,
			Type:        models.ContentTypeText,
			ReadingTime: article.ReadingTime,
			Reactions:   article.Reactions,
			Tags:        article.Tags,
			Language:    article.Language,
			PublishedAt: article.PublishedAt,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		contents = append(contents, content)
	}

	return contents, nil
} 