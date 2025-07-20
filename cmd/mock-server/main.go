package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Mock video data for JSON provider
type MockVideo struct {
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

// Mock article data for XML provider
type MockArticle struct {
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

// Mock response structures
type MockVideoResponse struct {
	Videos []MockVideo `json:"videos"`
}

type MockArticleResponse struct {
	XMLName xml.Name      `xml:"articles"`
	Articles []MockArticle `xml:"article"`
}

func main() {
	router := gin.Default()

	// JSON Provider (Videos)
	router.GET("/api/videos", func(c *gin.Context) {
		videos := []MockVideo{
			{
				ID:          "video_1",
				Title:       "Go Programming Tutorial for Beginners",
				Description: "Learn the basics of Go programming language with this comprehensive tutorial for beginners.",
				URL:         "https://example.com/videos/go-tutorial",
				Views:       15000,
				Likes:       1200,
				Duration:    1800, // 30 minutes
				Tags:        "golang,programming,tutorial,beginner",
				Language:    "en",
				PublishedAt: time.Now().AddDate(0, 0, -5),
			},
			{
				ID:          "video_2",
				Title:       "Advanced Go Concurrency Patterns",
				Description: "Explore advanced concurrency patterns in Go including goroutines, channels, and select statements.",
				URL:         "https://example.com/videos/go-concurrency",
				Views:       8500,
				Likes:       950,
				Duration:    2400, // 40 minutes
				Tags:        "golang,concurrency,advanced,patterns",
				Language:    "en",
				PublishedAt: time.Now().AddDate(0, 0, -10),
			},
			{
				ID:          "video_3",
				Title:       "Building REST APIs with Go and Gin",
				Description: "Learn how to build scalable REST APIs using Go and the Gin web framework.",
				URL:         "https://example.com/videos/go-rest-api",
				Views:       22000,
				Likes:       1800,
				Duration:    2700, // 45 minutes
				Tags:        "golang,api,rest,gin,web",
				Language:    "en",
				PublishedAt: time.Now().AddDate(0, 0, -2),
			},
		}

		response := MockVideoResponse{Videos: videos}
		c.JSON(http.StatusOK, response)
	})

	// XML Provider (Articles)
	router.GET("/api/articles", func(c *gin.Context) {
		articles := []MockArticle{
			{
				ID:          "article_1",
				Title:       "Understanding Go Modules and Dependency Management",
				Description: "A comprehensive guide to Go modules, dependency management, and best practices for modern Go development.",
				URL:         "https://example.com/articles/go-modules",
				ReadingTime: 8,
				Reactions:   450,
				Tags:        "golang,modules,dependencies,development",
				Language:    "en",
				PublishedAt: time.Now().AddDate(0, 0, -3),
			},
			{
				ID:          "article_2",
				Title:       "Microservices Architecture with Go",
				Description: "Learn how to design and implement microservices using Go, including service discovery and communication patterns.",
				URL:         "https://example.com/articles/go-microservices",
				ReadingTime: 12,
				Reactions:   320,
				Tags:        "golang,microservices,architecture,distributed-systems",
				Language:    "en",
				PublishedAt: time.Now().AddDate(0, 0, -7),
			},
			{
				ID:          "article_3",
				Title:       "Testing Strategies for Go Applications",
				Description: "Explore different testing strategies and tools for Go applications, from unit tests to integration tests.",
				URL:         "https://example.com/articles/go-testing",
				ReadingTime: 10,
				Reactions:   280,
				Tags:        "golang,testing,unit-tests,integration-tests",
				Language:    "en",
				PublishedAt: time.Now().AddDate(0, 0, -1),
			},
			{
				ID:          "article_4",
				Title:       "Performance Optimization in Go",
				Description: "Learn techniques for optimizing Go applications, including profiling, benchmarking, and memory management.",
				URL:         "https://example.com/articles/go-performance",
				ReadingTime: 15,
				Reactions:   520,
				Tags:        "golang,performance,optimization,profiling",
				Language:    "en",
				PublishedAt: time.Now().AddDate(0, 0, -15),
			},
		}

		response := MockArticleResponse{Articles: articles}
		c.Header("Content-Type", "application/xml")
		c.XML(http.StatusOK, response)
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "mock-provider-server",
		})
	})

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Mock Provider Server",
			"version": "1.0.0",
			"endpoints": gin.H{
				"videos":   "/api/videos",
				"articles": "/api/articles",
				"health":   "/health",
			},
		})
	})

	port := ":3001"
	fmt.Printf("Mock server starting on port %s\n", port)
	fmt.Printf("JSON Provider: http://localhost%s/api/videos\n", port)
	fmt.Printf("XML Provider: http://localhost%s/api/articles\n", port)
	
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start mock server:", err)
	}
} 