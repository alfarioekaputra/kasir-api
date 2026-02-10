package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DatabaseType string

const (
	PostgreSQL DatabaseType = "postgres"
	MySQL      DatabaseType = "mysql"
)

type Config struct {
	Type             DatabaseType
	ConnectionString string
	MaxOpenConns     int
	MaxIdleConns     int
}

// InitDB initializes database connection with support for PostgreSQL and MySQL
func InitDB(config Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	// Select database driver based on type
	switch config.Type {
	case PostgreSQL:
		dialector = postgres.Open(config.ConnectionString)
	case MySQL:
		dialector = mysql.Open(config.ConnectionString)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}

	// Open database connection
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB to set connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Set connection pool settings
	maxOpenConns := config.MaxOpenConns
	if maxOpenConns == 0 {
		maxOpenConns = 25
	}
	sqlDB.SetMaxOpenConns(maxOpenConns)

	maxIdleConns := config.MaxIdleConns
	if maxIdleConns == 0 {
		maxIdleConns = 5
	}
	sqlDB.SetMaxIdleConns(maxIdleConns)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Database connected successfully (Type: %s)", config.Type)
	return db, nil
}

// AutoMigrate runs auto-migration for all models
func AutoMigrate(db *gorm.DB, models ...interface{}) error {
	return db.AutoMigrate(models...)
}
