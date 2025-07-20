package repository

import (
	"math"
	"strings"

	"search-engine-service/internal/database/models"

	"gorm.io/gorm"
)

type ContentRepositoryImpl struct {
	db *gorm.DB
}

func NewContentRepository(db *gorm.DB) models.ContentRepository {
	return &ContentRepositoryImpl{db: db}
}

func (r *ContentRepositoryImpl) Create(content *models.Content) error {
	return r.db.Create(content).Error
}

func (r *ContentRepositoryImpl) Update(content *models.Content) error {
	return r.db.Save(content).Error
}

func (r *ContentRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&models.Content{}, id).Error
}

func (r *ContentRepositoryImpl) FindByID(id uint) (*models.Content, error) {
	var content models.Content
	err := r.db.First(&content, id).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func (r *ContentRepositoryImpl) FindByProviderID(provider, providerID string) (*models.Content, error) {
	var content models.Content
	err := r.db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&content).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func (r *ContentRepositoryImpl) Search(query string, contentType models.ContentType, page, limit int) (*models.SearchResult, error) {
	var contents []models.Content
	var total int64

	// Build query
	dbQuery := r.db.Model(&models.Content{})

	// Add search conditions
	if query != "" {
		searchQuery := "%" + strings.ToLower(query) + "%"
		dbQuery = dbQuery.Where(
			"LOWER(title) LIKE ? OR LOWER(description) LIKE ? OR LOWER(tags) LIKE ?",
			searchQuery, searchQuery, searchQuery,
		)
	}

	// Add content type filter
	if contentType != "" && contentType != "all" {
		dbQuery = dbQuery.Where("type = ?", contentType)
	}

	// Count total
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	// Calculate pagination
	offset := (page - 1) * limit
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	// Get results
	err := dbQuery.Order("final_score DESC").
		Offset(offset).
		Limit(limit).
		Find(&contents).Error

	if err != nil {
		return nil, err
	}

	return &models.SearchResult{
		Contents:    contents,
		Total:       total,
		Page:        page,
		Limit:       limit,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}, nil
}

func (r *ContentRepositoryImpl) GetPopular(limit int) ([]models.Content, error) {
	var contents []models.Content
	err := r.db.Order("final_score DESC").Limit(limit).Find(&contents).Error
	return contents, err
}

func (r *ContentRepositoryImpl) UpdateScores() error {
	// This method will be called periodically to update scores
	// For now, we'll just return nil as scores are calculated on-the-fly
	return nil
}

func (r *ContentRepositoryImpl) BulkUpsert(contents []models.Content) error {
	if len(contents) == 0 {
		return nil
	}

	// Use transaction for bulk operations
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, content := range contents {
			// Check if content exists
			var existing models.Content
			err := tx.Where("provider = ? AND provider_id = ?", content.Provider, content.ProviderID).First(&existing).Error
			
			if err == gorm.ErrRecordNotFound {
				// Create new content
				if err := tx.Create(&content).Error; err != nil {
					return err
				}
			} else if err == nil {
				// Update existing content
				content.ID = existing.ID
				if err := tx.Save(&content).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}
		return nil
	})
} 