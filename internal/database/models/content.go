package models

import (
	"time"

	"gorm.io/gorm"
)

// ContentType represents the type of content
type ContentType string

const (
	ContentTypeVideo ContentType = "video"
	ContentTypeText  ContentType = "text"
)

// Content represents the main content model
type Content struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" gorm:"size:255;not null"`
	Description string         `json:"description" gorm:"type:text"`
	URL         string         `json:"url" gorm:"size:500"`
	Type        ContentType    `json:"type" gorm:"size:20;not null"`
	Provider    string         `json:"provider" gorm:"size:100;not null"`
	ProviderID  string         `json:"provider_id" gorm:"size:100;not null"`
	
	// Video specific fields
	Views       int    `json:"views" gorm:"default:0"`
	Likes       int    `json:"likes" gorm:"default:0"`
	Duration    int    `json:"duration" gorm:"default:0"` // in seconds
	
	// Text specific fields
	ReadingTime int    `json:"reading_time" gorm:"default:0"` // in minutes
	Reactions   int    `json:"reactions" gorm:"default:0"`
	
	// Scoring fields
	BaseScore      float64 `json:"base_score" gorm:"type:decimal(10,4);default:0"`
	TypeMultiplier float64 `json:"type_multiplier" gorm:"type:decimal(10,4);default:1"`
	FreshnessScore float64 `json:"freshness_score" gorm:"type:decimal(10,4);default:0"`
	EngagementScore float64 `json:"engagement_score" gorm:"type:decimal(10,4);default:0"`
	FinalScore     float64 `json:"final_score" gorm:"type:decimal(10,4);default:0"`
	
	// Metadata
	Tags        string    `json:"tags" gorm:"type:text"`
	Language    string    `json:"language" gorm:"size:10;default:'en'"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName specifies the table name for Content
func (Content) TableName() string {
	return "contents"
}

// SearchResult represents a search result with pagination
type SearchResult struct {
	Contents    []Content `json:"contents"`
	Total       int64     `json:"total"`
	Page        int       `json:"page"`
	Limit       int       `json:"limit"`
	TotalPages  int       `json:"total_pages"`
	HasNext     bool      `json:"has_next"`
	HasPrevious bool      `json:"has_previous"`
}

// ContentRepository interface defines the methods for content operations
type ContentRepository interface {
	Create(content *Content) error
	Update(content *Content) error
	Delete(id uint) error
	FindByID(id uint) (*Content, error)
	FindByProviderID(provider, providerID string) (*Content, error)
	Search(query string, contentType ContentType, page, limit int) (*SearchResult, error)
	GetPopular(limit int) ([]Content, error)
	UpdateScores() error
	BulkUpsert(contents []Content) error
} 