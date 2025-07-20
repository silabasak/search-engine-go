package logger

import (
	"context"
	"time"

	"search-engine-service/internal/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger interface for structured logging
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	
	WithContext(ctx context.Context) Logger
	WithFields(fields ...zap.Field) Logger
	
	// HTTP request logging
	LogHTTPRequest(method, path, remoteAddr string, statusCode int, duration time.Duration, userAgent string)
	
	// Business logic logging
	LogSearchQuery(query string, contentType string, page, limit int, duration time.Duration)
	LogProviderFetch(provider string, count int, duration time.Duration, err error)
	LogContentScore(contentID uint, score float64, breakdown map[string]float64)
}

// logger implements the Logger interface
type logger struct {
	zap *zap.Logger
}

// NewLogger creates a new structured logger
func NewLogger(cfg *config.Config) (Logger, error) {
	var zapConfig zap.Config
	
	if cfg.Environment == "production" {
		zapConfig = zap.NewProductionConfig()
		zapConfig.EncoderConfig.TimeKey = "timestamp"
		zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	
	zapConfig.OutputPaths = []string{"stdout"}
	zapConfig.ErrorOutputPaths = []string{"stderr"}
	
	zapLogger, err := zapConfig.Build()
	if err != nil {
		return nil, err
	}
	
	return &logger{zap: zapLogger}, nil
}

// Debug logs a debug message
func (l *logger) Debug(msg string, fields ...zap.Field) {
	l.zap.Debug(msg, fields...)
}

// Info logs an info message
func (l *logger) Info(msg string, fields ...zap.Field) {
	l.zap.Info(msg, fields...)
}

// Warn logs a warning message
func (l *logger) Warn(msg string, fields ...zap.Field) {
	l.zap.Warn(msg, fields...)
}

// Error logs an error message
func (l *logger) Error(msg string, fields ...zap.Field) {
	l.zap.Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func (l *logger) Fatal(msg string, fields ...zap.Field) {
	l.zap.Fatal(msg, fields...)
}

// WithContext creates a logger with context
func (l *logger) WithContext(ctx context.Context) Logger {
	if ctx == nil {
		return l
	}
	
	// Extract request ID from context if available
	if requestID, ok := ctx.Value("request_id").(string); ok {
		return &logger{zap: l.zap.With(zap.String("request_id", requestID))}
	}
	
	return l
}

// WithFields creates a logger with additional fields
func (l *logger) WithFields(fields ...zap.Field) Logger {
	return &logger{zap: l.zap.With(fields...)}
}

// LogHTTPRequest logs HTTP request details
func (l *logger) LogHTTPRequest(method, path, remoteAddr string, statusCode int, duration time.Duration, userAgent string) {
	level := zap.InfoLevel
	if statusCode >= 400 {
		level = zap.WarnLevel
	}
	if statusCode >= 500 {
		level = zap.ErrorLevel
	}
	
	fields := []zap.Field{
		zap.String("method", method),
		zap.String("path", path),
		zap.String("remote_addr", remoteAddr),
		zap.Int("status_code", statusCode),
		zap.Duration("duration", duration),
		zap.String("user_agent", userAgent),
	}
	
	l.zap.Check(level, "HTTP Request").Write(fields...)
}

// LogSearchQuery logs search query details
func (l *logger) LogSearchQuery(query string, contentType string, page, limit int, duration time.Duration) {
	l.Info("Search query executed",
		zap.String("query", query),
		zap.String("content_type", contentType),
		zap.Int("page", page),
		zap.Int("limit", limit),
		zap.Duration("duration", duration),
	)
}

// LogProviderFetch logs provider fetch details
func (l *logger) LogProviderFetch(provider string, count int, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("provider", provider),
		zap.Int("count", count),
		zap.Duration("duration", duration),
	}
	
	if err != nil {
		fields = append(fields, zap.Error(err))
		l.Error("Provider fetch failed", fields...)
	} else {
		l.Info("Provider fetch completed", fields...)
	}
}

// LogContentScore logs content scoring details
func (l *logger) LogContentScore(contentID uint, score float64, breakdown map[string]float64) {
	fields := []zap.Field{
		zap.Uint("content_id", contentID),
		zap.Float64("score", score),
	}
	
	for key, value := range breakdown {
		fields = append(fields, zap.Float64(key, value))
	}
	
	l.Debug("Content scored", fields...)
}

// Sync flushes any buffered log entries
func (l *logger) Sync() error {
	return l.zap.Sync()
}

// Helper functions for common logging patterns
func (l *logger) LogDatabaseQuery(query string, duration time.Duration, err error) {
	fields := []zap.Field{
		zap.String("query", query),
		zap.Duration("duration", duration),
	}
	
	if err != nil {
		fields = append(fields, zap.Error(err))
		l.Error("Database query failed", fields...)
	} else {
		l.Debug("Database query executed", fields...)
	}
}

func (l *logger) LogCacheHit(key string) {
	l.Debug("Cache hit", zap.String("key", key))
}

func (l *logger) LogCacheMiss(key string) {
	l.Debug("Cache miss", zap.String("key", key))
}

func (l *logger) LogRateLimitExceeded(clientIP string) {
	l.Warn("Rate limit exceeded", zap.String("client_ip", clientIP))
}

func (l *logger) LogSecurityEvent(event string, details map[string]interface{}) {
	fields := make([]zap.Field, 0, len(details))
	for key, value := range details {
		fields = append(fields, zap.Any(key, value))
	}
	
	l.Warn("Security event", append(fields, zap.String("event", event))...)
} 