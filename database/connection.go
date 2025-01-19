package database

import (
	"fmt"
	"log"
	"myproject/config"
	"os"
	"time"

	"gorm.io/driver/postgres" // or mysql, sqlite, etc.
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	*gorm.DB
}

// NewDatabase creates a new database instance
func NewDatabase(config *config.DBConfig) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		// config.Password,
		config.DBName,
		// config.SSLMode,
	)

	// Custom logger configuration
	loggerConfig := logger.Config{
		SlowThreshold:             time.Second, // Log slow queries
		LogLevel:                  logger.Info, // Log level
		IgnoreRecordNotFoundError: true,        // Ignore record not found errors
		Colorful:                  true,        // Enable color
	}

	// Database configuration
	gormConfig := &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			loggerConfig,
		),
		PrepareStmt:            true, // Cache prepared statements
		SkipDefaultTransaction: true, // Skip default transaction
	}

	log.Printf("Attempting to connect with: %s", dsn)

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)           // Set maximum number of idle connections
	sqlDB.SetMaxOpenConns(100)          // Set maximum number of open connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Set maximum lifetime of connections

	return &Database{db}, nil
}
