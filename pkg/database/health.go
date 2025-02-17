package database

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

// Health performs a comprehensive health check of the database connection.
// It returns a map containing:
// - Connection status (up/down)
// - Error messages if any
// - Connection pool statistics
// - Performance metrics and warnings
// The function uses a context with timeout to ensure health checks complete quickly.
func (s *service) Health() map[string]string {
	s.logger.Debug().Msg("Starting database health check")
	stats := make(map[string]string)

	// Get the underlying *sql.DB instance
	sqlDB, err := s.orm.DB()
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("Failed to get database instance during health check")

		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("failed to get database instance: %v", err)
		return stats
	}

	// Create a context with timeout for the health check
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// Ping the database
	s.logger.Debug().Msg("Pinging database for health check")
	err = sqlDB.PingContext(ctx)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("timeout", "1s").
			Msg("Database ping failed during health check")

		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats
	dbStats := sqlDB.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Log detailed connection statistics
	s.logger.Debug().
		Int("open_connections", dbStats.OpenConnections).
		Int("in_use", dbStats.InUse).
		Int("idle", dbStats.Idle).
		Int64("wait_count", dbStats.WaitCount).
		Str("wait_duration", dbStats.WaitDuration.String()).
		Int64("max_idle_closed", dbStats.MaxIdleClosed).
		Int64("max_lifetime_closed", dbStats.MaxLifetimeClosed).
		Msg("Database connection pool statistics")

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 {
		s.logger.Warn().
			Int("open_connections", dbStats.OpenConnections).
			Int("threshold", 40).
			Msg("High number of open connections detected")
		stats["message"] = "The database is experiencing heavy load."
	}
	if dbStats.WaitCount > 1000 {
		s.logger.Warn().
			Int64("wait_count", dbStats.WaitCount).
			Int("threshold", 1000).
			Msg("High number of connection wait events detected")
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}
	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		s.logger.Warn().
			Int64("max_idle_closed", dbStats.MaxIdleClosed).
			Int("open_connections", dbStats.OpenConnections).
			Msg("High number of idle connections being closed")
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}
	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		s.logger.Warn().
			Int64("max_lifetime_closed", dbStats.MaxLifetimeClosed).
			Int("open_connections", dbStats.OpenConnections).
			Msg("High number of connections being closed due to max lifetime")
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	s.logger.Info().
		Str("status", stats["status"]).
		Str("message", stats["message"]).
		Msg("Database health check completed")

	return stats
}
