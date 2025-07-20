package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"search-engine-service/internal/utils/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// SecurityMiddleware provides security-related middleware
type SecurityMiddleware struct {
	logger logger.Logger
}

// NewSecurityMiddleware creates a new security middleware
func NewSecurityMiddleware(logger logger.Logger) *SecurityMiddleware {
	return &SecurityMiddleware{
		logger: logger,
	}
}

// SecurityHeaders adds security headers to responses
func (sm *SecurityMiddleware) SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' https:; connect-src 'self' https:;")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		c.Header("X-Download-Options", "noopen")
		c.Header("X-Permitted-Cross-Domain-Policies", "none")
		
		c.Next()
	}
}

// CORS configures CORS middleware
func (sm *SecurityMiddleware) CORS() gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080", "https://yourdomain.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length", "X-Total-Count"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	
	return cors.New(config)
}

// RequestID adds a unique request ID to each request
func (sm *SecurityMiddleware) RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)
		
		// Add request ID to context
		ctx := context.WithValue(c.Request.Context(), "request_id", requestID)
		c.Request = c.Request.WithContext(ctx)
		
		c.Next()
	}
}

// InputSanitizer sanitizes input parameters
func (sm *SecurityMiddleware) InputSanitizer() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Sanitize query parameters
		for key, values := range c.Request.URL.Query() {
			for i, value := range values {
				sanitized := sm.sanitizeString(value)
				values[i] = sanitized
			}
			c.Request.URL.Query()[key] = values
		}
		
		// Sanitize path parameters
		for _, param := range c.Params {
			param.Value = sm.sanitizeString(param.Value)
		}
		
		c.Next()
	}
}

// sanitizeString removes potentially dangerous characters
func (sm *SecurityMiddleware) sanitizeString(input string) string {
	if input == "" {
		return input
	}
	
	// XSS protection
	dangerous := []string{
		"<script>", "</script>", "javascript:", "onload=", "onerror=", "onclick=",
		"onmouseover=", "onfocus=", "onblur=", "onchange=", "onsubmit=",
		"<iframe>", "</iframe>", "<object>", "</object>", "<embed>",
		"vbscript:", "data:", "mocha:", "livescript:",
	}
	
	result := input
	for _, dangerous := range dangerous {
		result = strings.ReplaceAll(strings.ToLower(result), strings.ToLower(dangerous), "")
	}
	
	// SQL injection protection (basic)
	sqlKeywords := []string{
		"union", "select", "insert", "update", "delete", "drop", "create",
		"alter", "exec", "execute", "declare", "cast", "convert",
	}
	
	for _, keyword := range sqlKeywords {
		// This is a basic check - in production, use proper parameterized queries
		if strings.Contains(strings.ToLower(result), keyword) {
			sm.logger.Warn("Potential SQL injection attempt detected",
				zap.String("input", input),
				zap.String("keyword", keyword),
			)
		}
	}
	
	return result
}

// RateLimiter implements rate limiting
type RateLimiter struct {
	logger logger.Logger
	clients map[string][]time.Time
	requestsPerMinute int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(logger logger.Logger, requestsPerMinute int) *RateLimiter {
	return &RateLimiter{
		logger: logger,
		clients: make(map[string][]time.Time),
		requestsPerMinute: requestsPerMinute,
	}
}

// Limit implements rate limiting middleware
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := rl.getClientIP(c)
		now := time.Now()
		
		// Clean old requests
		if times, exists := rl.clients[clientIP]; exists {
			var validTimes []time.Time
			for _, t := range times {
				if now.Sub(t) < time.Minute {
					validTimes = append(validTimes, t)
				}
			}
			rl.clients[clientIP] = validTimes
		}
		
		// Check rate limit
		if times, exists := rl.clients[clientIP]; exists && len(times) >= rl.requestsPerMinute {
			rl.logger.Warn("Rate limit exceeded",
				zap.String("client_ip", clientIP),
				zap.String("path", c.Request.URL.Path),
			)
			
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
				"retry_after": 60,
			})
			c.Abort()
			return
		}
		
		// Add current request
		rl.clients[clientIP] = append(rl.clients[clientIP], now)
		
		// Add rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", rl.requestsPerMinute))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", rl.requestsPerMinute-len(rl.clients[clientIP])))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", now.Add(time.Minute).Unix()))
		
		c.Next()
	}
}

// getClientIP gets the real client IP address
func (rl *RateLimiter) getClientIP(c *gin.Context) string {
	// Check for forwarded headers
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		return strings.Split(ip, ",")[0]
	}
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := c.GetHeader("X-Client-IP"); ip != "" {
		return ip
	}
	
	return c.ClientIP()
}

// RequestLogger logs HTTP requests
type RequestLogger struct {
	logger logger.Logger
}

// NewRequestLogger creates a new request logger
func NewRequestLogger(logger logger.Logger) *RequestLogger {
	return &RequestLogger{
		logger: logger,
	}
}

// Log implements request logging middleware
func (rl *RequestLogger) Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		
		// Process request
		c.Next()
		
		// Calculate duration
		duration := time.Since(start)
		
		// Get request ID
		requestID, _ := c.Get("request_id")
		
		// Log request
		rl.logger.LogHTTPRequest(
			c.Request.Method,
			path,
			c.ClientIP(),
			c.Writer.Status(),
			duration,
			c.Request.UserAgent(),
		)
		
		// Log additional details for errors
		if c.Writer.Status() >= 400 {
			rl.logger.Error("HTTP request failed",
				zap.String("request_id", requestID.(string)),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.String("query", raw),
				zap.Int("status", c.Writer.Status()),
				zap.Duration("duration", duration),
				zap.String("client_ip", c.ClientIP()),
				zap.String("user_agent", c.Request.UserAgent()),
			)
		}
	}
}

// ErrorHandler handles panics and errors
func (sm *SecurityMiddleware) ErrorHandler() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			sm.logger.Error("Panic recovered",
				zap.String("error", err),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)
		}
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code": "INTERNAL_ERROR",
		})
	})
} 