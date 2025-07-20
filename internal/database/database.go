package database

import (
	"fmt"
	"log"

	"search-engine-service/internal/config"
	"search-engine-service/internal/database/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database represents the database connection
type Database struct {
	DB *gorm.DB
}

var DB *gorm.DB

// NewDatabase creates a new database instance
func NewDatabase(cfg *config.Config) (*Database, error) {
	db, err := Init(cfg.Database)
	if err != nil {
		return nil, err
	}
	
	return &Database{DB: db}, nil
}

// Init initializes the database connection
func Init(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	DB = db

	// Auto migrate the schema
	if err := AutoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	log.Println("Database initialized successfully")
	return db, nil
}

// AutoMigrate runs database migrations
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Content{},
	)
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
} 