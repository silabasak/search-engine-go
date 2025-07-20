package tests

import (
	"testing"
	"time"

	"search-engine-service/internal/database/models"
	"search-engine-service/internal/services"
)

func TestScoringService(t *testing.T) {
	scoringService := services.NewScoringService()

	t.Run("Test Video Scoring", func(t *testing.T) {
		content := &models.Content{
			Type:        models.ContentTypeVideo,
			Views:       10000,
			Likes:       800,
			Duration:    1800,
			PublishedAt: time.Now().AddDate(0, 0, -3), // 3 days ago
		}

		scoringService.CalculateScore(content)

		// Expected calculations:
		// Base score: 10000/1000 + 800/100 = 10 + 8 = 18
		// Type multiplier: 1.5
		// Freshness score: 5 (within 1 week)
		// Engagement score: (800/10000) * 10 = 0.8
		// Final score: (18 * 1.5) + 5 + 0.8 = 27 + 5 + 0.8 = 32.8

		if content.BaseScore != 18.0 {
			t.Errorf("Expected base score 18.0, got %f", content.BaseScore)
		}

		if content.TypeMultiplier != 1.5 {
			t.Errorf("Expected type multiplier 1.5, got %f", content.TypeMultiplier)
		}

		if content.FreshnessScore != 5.0 {
			t.Errorf("Expected freshness score 5.0, got %f", content.FreshnessScore)
		}

		if content.EngagementScore != 0.8 {
			t.Errorf("Expected engagement score 0.8, got %f", content.EngagementScore)
		}

		expectedFinalScore := 32.8
		if content.FinalScore != expectedFinalScore {
			t.Errorf("Expected final score %f, got %f", expectedFinalScore, content.FinalScore)
		}
	})

	t.Run("Test Text Scoring", func(t *testing.T) {
		content := &models.Content{
			Type:        models.ContentTypeText,
			ReadingTime: 10,
			Reactions:   200,
			PublishedAt: time.Now().AddDate(0, 0, -20), // 20 days ago
		}

		scoringService.CalculateScore(content)

		// Expected calculations:
		// Base score: 10 + 200/50 = 10 + 4 = 14
		// Type multiplier: 1.0
		// Freshness score: 3 (within 1 month)
		// Engagement score: (200/10) * 5 = 100
		// Final score: (14 * 1.0) + 3 + 100 = 14 + 3 + 100 = 117

		if content.BaseScore != 14.0 {
			t.Errorf("Expected base score 14.0, got %f", content.BaseScore)
		}

		if content.TypeMultiplier != 1.0 {
			t.Errorf("Expected type multiplier 1.0, got %f", content.TypeMultiplier)
		}

		if content.FreshnessScore != 3.0 {
			t.Errorf("Expected freshness score 3.0, got %f", content.FreshnessScore)
		}

		if content.EngagementScore != 100.0 {
			t.Errorf("Expected engagement score 100.0, got %f", content.EngagementScore)
		}

		expectedFinalScore := 117.0
		if content.FinalScore != expectedFinalScore {
			t.Errorf("Expected final score %f, got %f", expectedFinalScore, content.FinalScore)
		}
	})

	t.Run("Test Old Content Scoring", func(t *testing.T) {
		content := &models.Content{
			Type:        models.ContentTypeVideo,
			Views:       5000,
			Likes:       300,
			PublishedAt: time.Now().AddDate(0, 0, -100), // 100 days ago
		}

		scoringService.CalculateScore(content)

		// Expected calculations:
		// Base score: 5000/1000 + 300/100 = 5 + 3 = 8
		// Type multiplier: 1.5
		// Freshness score: 0 (older than 3 months)
		// Engagement score: (300/5000) * 10 = 0.6
		// Final score: (8 * 1.5) + 0 + 0.6 = 12 + 0 + 0.6 = 12.6

		if content.FreshnessScore != 0.0 {
			t.Errorf("Expected freshness score 0.0, got %f", content.FreshnessScore)
		}

		expectedFinalScore := 12.6
		if content.FinalScore != expectedFinalScore {
			t.Errorf("Expected final score %f, got %f", expectedFinalScore, content.FinalScore)
		}
	})
}

func TestScoreBreakdown(t *testing.T) {
	scoringService := services.NewScoringService()

	content := &models.Content{
		Type:        models.ContentTypeVideo,
		Views:       10000,
		Likes:       800,
		PublishedAt: time.Now().AddDate(0, 0, -5),
	}

	breakdown := scoringService.GetScoreBreakdown(content)

	expectedKeys := []string{"base_score", "type_multiplier", "freshness_score", "engagement_score", "final_score"}
	for _, key := range expectedKeys {
		if _, exists := breakdown[key]; !exists {
			t.Errorf("Expected breakdown to contain key: %s", key)
		}
	}

	if breakdown["final_score"] <= 0 {
		t.Errorf("Expected final score to be positive, got %f", breakdown["final_score"])
	}
}

func TestNormalizeScore(t *testing.T) {
	scoringService := services.NewScoringService()

	tests := []struct {
		score    float64
		maxScore float64
		expected float64
	}{
		{50.0, 100.0, 50.0},
		{25.0, 100.0, 25.0},
		{0.0, 100.0, 0.0},
		{100.0, 100.0, 100.0},
		{0.0, 0.0, 0.0}, // Edge case: maxScore is 0
	}

	for _, test := range tests {
		result := scoringService.NormalizeScore(test.score, test.maxScore)
		if result != test.expected {
			t.Errorf("NormalizeScore(%f, %f) = %f, expected %f", test.score, test.maxScore, result, test.expected)
		}
	}
} 