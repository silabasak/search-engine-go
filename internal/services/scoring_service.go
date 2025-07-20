package services

import (
	"time"

	"search-engine-service/internal/database/models"
)

// ScoringService handles content scoring calculations
type ScoringService struct{}

// NewScoringService creates a new scoring service
func NewScoringService() *ScoringService {
	return &ScoringService{}
}

// CalculateScore calculates the final score for a content item
func (ss *ScoringService) CalculateScore(content *models.Content) {
	// Calculate base score
	baseScore := ss.calculateBaseScore(content)
	
	// Calculate type multiplier
	typeMultiplier := ss.calculateTypeMultiplier(content.Type)
	
	// Calculate freshness score
	freshnessScore := ss.calculateFreshnessScore(content.PublishedAt)
	
	// Calculate engagement score
	engagementScore := ss.calculateEngagementScore(content)
	
	// Calculate final score
	finalScore := (baseScore * typeMultiplier) + freshnessScore + engagementScore
	
	// Update content scores
	content.BaseScore = baseScore
	content.TypeMultiplier = typeMultiplier
	content.FreshnessScore = freshnessScore
	content.EngagementScore = engagementScore
	content.FinalScore = finalScore
}

// calculateBaseScore calculates the base score based on content type
func (ss *ScoringService) calculateBaseScore(content *models.Content) float64 {
	switch content.Type {
	case models.ContentTypeVideo:
		// Video: views / 1000 + (likes / 100)
		viewsScore := float64(content.Views) / 1000.0
		likesScore := float64(content.Likes) / 100.0
		return viewsScore + likesScore
		
	case models.ContentTypeText:
		// Text: reading_time + (reactions / 50)
		readingScore := float64(content.ReadingTime)
		reactionsScore := float64(content.Reactions) / 50.0
		return readingScore + reactionsScore
		
	default:
		return 0
	}
}

// calculateTypeMultiplier returns the type multiplier
func (ss *ScoringService) calculateTypeMultiplier(contentType models.ContentType) float64 {
	switch contentType {
	case models.ContentTypeVideo:
		return 1.5
	case models.ContentTypeText:
		return 1.0
	default:
		return 1.0
	}
}

// calculateFreshnessScore calculates the freshness score based on publication date
func (ss *ScoringService) calculateFreshnessScore(publishedAt time.Time) float64 {
	now := time.Now()
	age := now.Sub(publishedAt)
	
	// Convert to days
	days := age.Hours() / 24
	
	switch {
	case days <= 7: // 1 week
		return 5.0
	case days <= 30: // 1 month
		return 3.0
	case days <= 90: // 3 months
		return 1.0
	default:
		return 0.0
	}
}

// calculateEngagementScore calculates the engagement score
func (ss *ScoringService) calculateEngagementScore(content *models.Content) float64 {
	switch content.Type {
	case models.ContentTypeVideo:
		// Video: (likes / views) * 10
		if content.Views > 0 {
			engagement := float64(content.Likes) / float64(content.Views)
			return engagement * 10.0
		}
		return 0
		
	case models.ContentTypeText:
		// Text: (reactions / reading_time) * 5
		if content.ReadingTime > 0 {
			engagement := float64(content.Reactions) / float64(content.ReadingTime)
			return engagement * 5.0
		}
		return 0
		
	default:
		return 0
	}
}

// CalculateScoresForBatch calculates scores for multiple content items
func (ss *ScoringService) CalculateScoresForBatch(contents []models.Content) []models.Content {
	for i := range contents {
		ss.CalculateScore(&contents[i])
	}
	return contents
}

// GetScoreBreakdown returns a detailed breakdown of the scoring
func (ss *ScoringService) GetScoreBreakdown(content *models.Content) map[string]float64 {
	ss.CalculateScore(content)
	
	return map[string]float64{
		"base_score":       content.BaseScore,
		"type_multiplier":  content.TypeMultiplier,
		"freshness_score":  content.FreshnessScore,
		"engagement_score": content.EngagementScore,
		"final_score":      content.FinalScore,
	}
}

// NormalizeScore normalizes a score to a 0-100 range
func (ss *ScoringService) NormalizeScore(score float64, maxScore float64) float64 {
	if maxScore == 0 {
		return 0
	}
	return (score / maxScore) * 100
}

// GetMaxPossibleScore calculates the maximum possible score for a content type
func (ss *ScoringService) GetMaxPossibleScore(contentType models.ContentType) float64 {
	// This is a simplified calculation - in a real system, you might want to
	// use historical data to determine realistic maximums
	
	baseScore := 100.0 // Assuming max base score
	typeMultiplier := ss.calculateTypeMultiplier(contentType)
	freshnessScore := 5.0 // Max freshness score
	engagementScore := 10.0 // Max engagement score
	
	return (baseScore * typeMultiplier) + freshnessScore + engagementScore
} 