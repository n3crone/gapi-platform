package database

import (
	"fmt"
	"reflect"

	"github.com/rs/zerolog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB defines the interface for database operations.
// It provides methods for health checking, connection management,
// and access to the underlying ORM instance.
type DB interface {
	// Health returns the current status of the database connection
	// along with various performance metrics and statistics.
	Health() map[string]string

	// Close properly terminates the database connection and releases
	// any associated resources.
	Close() error

	// GetOrm provides access to the underlying GORM database instance
	// for database operations.
	GetOrm() *gorm.DB

	// AutoMigrate automatically migrates database schema for given models
	AutoMigrate(models ...interface{}) error
}

// service implements the DB interface and manages the database connection.
// It encapsulates the GORM ORM instance and provides additional
// functionality for connection management and monitoring.
type service struct {
	orm    *gorm.DB
	logger zerolog.Logger
}

// GetOrm returns the GORM database instance for database operations.
func (s *service) GetOrm() *gorm.DB {
	s.logger.Debug().Msg("Getting GORM database instance")
	return s.orm
}

// New creates or returns an existing database connection.
// It implements a singleton pattern to ensure only one database
// connection is maintained throughout the application lifecycle.
// The function:
// - Reuses an existing connection if available
// - Creates a new connection with optimal pool settings
// - Configures connection pooling for performance
// Returns a DB interface for database operations
func New(dsn string, logger zerolog.Logger) (DB, error) {
	logger.Debug().Msg("Initializing database connection")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to connect to database")
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to get database instance")
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(0)

	logger.Info().Msg("Database connection pool configured")

	svc := &service{
		orm:    db,
		logger: logger,
	}

	// Verify connection with health check
	health := svc.Health()
	if health["status"] != "up" {
		logger.Error().
			Interface("health", health).
			Msg("Database health check failed after connection")
		return nil, fmt.Errorf("database health check failed: %s", health["error"])
	}

	logger.Info().
		Interface("health", health).
		Msg("Database connection established successfully")

	return svc, nil
}

// Close gracefully shuts down the database connection.
// It ensures all resources are properly released and logs
// the disconnection event.

func (s *service) Close() error {
	s.logger.Info().Msg("Initiating database connection shutdown")

	sqlDB, err := s.orm.DB()
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Failed to get database instance during shutdown")
		return fmt.Errorf("failed to get database instance: %v", err)
	}

	err = sqlDB.Close()
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Error during database connection shutdown")
		return err
	}

	s.logger.Info().Msg("Database connection closed successfully")
	return nil
}

func (s *service) AutoMigrate(models ...interface{}) error {
	s.logger.Info().Msg("Starting database auto-migration")

	for _, model := range models {
		modelType := reflect.TypeOf(model).String()
		s.logger.Debug().
			Str("model", modelType).
			Msg("Migrating model schema")

		if err := s.orm.AutoMigrate(model); err != nil {
			s.logger.Error().
				Err(err).
				Str("model", modelType).
				Msg("Failed to migrate model schema")
			return fmt.Errorf("failed to migrate %s: %v", modelType, err)
		}
	}

	s.logger.Info().Msg("Database migration completed successfully")
	return nil
}
