-- Create database if not exists
CREATE DATABASE IF NOT EXISTS search_engine CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Use the database
USE search_engine;

-- Create contents table
CREATE TABLE IF NOT EXISTS contents (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    url VARCHAR(500),
    type VARCHAR(20) NOT NULL,
    provider VARCHAR(100) NOT NULL,
    provider_id VARCHAR(100) NOT NULL,
    views INT DEFAULT 0,
    likes INT DEFAULT 0,
    duration INT DEFAULT 0,
    reading_time INT DEFAULT 0,
    reactions INT DEFAULT 0,
    base_score DECIMAL(10,4) DEFAULT 0,
    type_multiplier DECIMAL(10,4) DEFAULT 1,
    freshness_score DECIMAL(10,4) DEFAULT 0,
    engagement_score DECIMAL(10,4) DEFAULT 0,
    final_score DECIMAL(10,4) DEFAULT 0,
    tags TEXT,
    language VARCHAR(10) DEFAULT 'en',
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_provider_provider_id (provider, provider_id),
    INDEX idx_type (type),
    INDEX idx_final_score (final_score DESC),
    INDEX idx_published_at (published_at DESC),
    FULLTEXT idx_search (title, description, tags)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Insert sample data
INSERT INTO contents (title, description, url, type, provider, provider_id, views, likes, duration, reading_time, reactions, tags, language, published_at) VALUES
('Go Programming Tutorial for Beginners', 'Learn the basics of Go programming language with this comprehensive tutorial for beginners.', 'https://example.com/videos/go-tutorial', 'video', 'json_provider', 'video_1', 15000, 1200, 1800, 0, 0, 'golang,programming,tutorial,beginner', 'en', DATE_SUB(NOW(), INTERVAL 5 DAY)),
('Advanced Go Concurrency Patterns', 'Explore advanced concurrency patterns in Go including goroutines, channels, and select statements.', 'https://example.com/videos/go-concurrency', 'video', 'json_provider', 'video_2', 8500, 950, 2400, 0, 0, 'golang,concurrency,advanced,patterns', 'en', DATE_SUB(NOW(), INTERVAL 10 DAY)),
('Building REST APIs with Go and Gin', 'Learn how to build scalable REST APIs using Go and the Gin web framework.', 'https://example.com/videos/go-rest-api', 'video', 'json_provider', 'video_3', 22000, 1800, 2700, 0, 0, 'golang,api,rest,gin,web', 'en', DATE_SUB(NOW(), INTERVAL 2 DAY)),
('Understanding Go Modules and Dependency Management', 'A comprehensive guide to Go modules, dependency management, and best practices for modern Go development.', 'https://example.com/articles/go-modules', 'text', 'xml_provider', 'article_1', 0, 0, 0, 8, 450, 'golang,modules,dependencies,development', 'en', DATE_SUB(NOW(), INTERVAL 3 DAY)),
('Microservices Architecture with Go', 'Learn how to design and implement microservices using Go, including service discovery and communication patterns.', 'https://example.com/articles/go-microservices', 'text', 'xml_provider', 'article_2', 0, 0, 0, 12, 320, 'golang,microservices,architecture,distributed-systems', 'en', DATE_SUB(NOW(), INTERVAL 7 DAY)),
('Testing Strategies for Go Applications', 'Explore different testing strategies and tools for Go applications, from unit tests to integration tests.', 'https://example.com/articles/go-testing', 'text', 'xml_provider', 'article_3', 0, 0, 0, 10, 280, 'golang,testing,unit-tests,integration-tests', 'en', DATE_SUB(NOW(), INTERVAL 1 DAY)),
('Performance Optimization in Go', 'Learn techniques for optimizing Go applications, including profiling, benchmarking, and memory management.', 'https://example.com/articles/go-performance', 'text', 'xml_provider', 'article_4', 0, 0, 0, 15, 520, 'golang,performance,optimization,profiling', 'en', DATE_SUB(NOW(), INTERVAL 15 DAY)); 