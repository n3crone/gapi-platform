package core

import (
	"fmt"
	"os"

	"github.com/n3crone/gapi-platform/pkg/database"
	"github.com/n3crone/gapi-platform/pkg/resource"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

// App represents the main application structure that combines Fiber web framework
// with resource management and database connectivity.
type App struct {
	Fiber *fiber.App                // Embedded Fiber application instance
	Db    database.DB               // Database connection interface
	rm    *resource.ResourceManager // Resource manager for handling API resources
	log   zerolog.Logger            // Application logger
}

type Config struct {
	FiberConfig *fiber.Config // Fiber configuration settings
	DatabaseUri string        // Database connection URI
	LogLevel    zerolog.Level // Log level for the application
	LogFormat   string        // Log format for the application
}

// New creates and initializes a new App instance with the provided configuration.
// It sets up the core components of the application:
//   - Configures structured logging with the specified level and format
//   - Establishes a database connection using the provided URI
//   - Initializes a Fiber web server with custom or default configuration
//   - Sets up a resource manager for API endpoint handling
//
// Example usage:
//
//	app, err := core.New(core.Config{
//		DatabaseUri: "postgres://user:pass@localhost:5432/dbname",
//		LogLevel:    zerolog.InfoLevel,
//		LogFormat:   "console",
//		Fiber: &fiber.Config{
//			AppName: "my-api",
//		},
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Returns:
//   - *App: The initialized application instance
//   - error: Any error that occurred during initialization
func New(config Config) (*App, error) {
	if config.DatabaseUri == "" {
		return nil, fmt.Errorf("DatabaseUri is required")
	}

	logger := configureLogger(config.LogLevel, config.LogFormat)
	logger.Debug().
		Interface("config", config).
		Msg("Initializing application with configuration")

	fiberConfig := fiber.Config{
		AppName: "gapi-platform",
	}
	if config.FiberConfig != nil {
		fiberConfig = *config.FiberConfig
		logger.Debug().
			Interface("fiber_config", fiberConfig).
			Msg("Using custom Fiber configuration")
	}

	logger.Info().Msg("Establishing database connection")
	db, err := database.New(config.DatabaseUri, logger)
	if err != nil {
		logger.Fatal().
			Err(err).
			Str("database_uri", config.DatabaseUri).
			Msg("Failed to connect to database")
		return nil, err
	}

	logger.Info().Msg("Initializing resource manager")
	rm := resource.NewResourceManager(db.GetOrm(), &logger)

	app := &App{
		Fiber: fiber.New(fiberConfig),
		Db:    db,
		rm:    rm,
		log:   logger,
	}

	logger.Info().
		Str("app_name", fiberConfig.AppName).
		Msg("Application initialized successfully")

	return app, nil
}

// Migrate runs database migrations for the provided models.
// It automatically creates or updates database tables based on the model structures.
//
// Example usage:
//
//	err := app.Migrate(
//		&User{},
//		&Product{},
//		&Order{},
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// Parameters:
//   - models: Variable number of struct instances that represent database models
//
// Returns:
//   - error: Any error that occurred during migration, nil on success
func (a *App) Migrate(models ...interface{}) error {
	a.log.Info().Msg("Running database migrations")
	return a.Db.AutoMigrate(models...)
}

// configureLogger sets up the zerolog logger with the specified level and format.
// If level is not provided (0), it defaults to Debug level.
// Format can be either "json" or "console" (pretty print).
//
// Parameters:
//   - level: The minimum log level to output (Debug, Info, Warn, Error, Fatal)
//   - format: The output format ("json" or "console")
//
// Returns:
//   - zerolog.Logger: Configured logger instance
func configureLogger(level zerolog.Level, format string) zerolog.Logger {
	if level == 0 {
		level = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var logger zerolog.Logger
	if format == "json" {
		logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	} else {
		output := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 15:04:05",
		}
		logger = zerolog.New(output).With().Timestamp().Logger()
	}

	logger.Debug().
		Str("log_level", level.String()).
		Str("log_format", format).
		Msg("Logger configured")

	return logger
}
