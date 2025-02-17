package core

import (
	"fmt"

	"github.com/n3crone/gapi-platform/pkg/resource"

	"github.com/gofiber/fiber/v2"
)

// RegisterResource registers a new API resource with the application.
// It takes a Registrable resource interface and:
// - Creates a new resource instance with the resource manager
// - Registers all CRUD routes for the resource with the Fiber app
func (a *App) RegisterResource(resource resource.Registrable) {
	resourceType := fmt.Sprintf("%T", resource)

	a.log.Info().
		Str("resource_type", resourceType).
		Msg("Starting resource registration")

	newResource := resource.CreateResource(a.rm)
	config := newResource.Config()

	// Create detailed operation info for logging
	opDetails := make(map[string]map[string]string)
	for op, cfg := range config.Operations {
		if cfg.Enabled {
			opDetails[string(op)] = map[string]string{
				"enabled":   "true",
				"provider":  fmt.Sprintf("%T", cfg.Provider),
				"processor": fmt.Sprintf("%T", cfg.Processor),
			}
		}
	}

	a.log.Debug().
		Str("resource_type", resourceType).
		Str("resource_path", config.Path).
		Interface("resource_model", config.Model).
		Interface("resource_operations", getOperationNames(config.Operations)).
		Interface("operation_details", opDetails).
		Msg("Resource created with configuration")

	newResource.RegisterRoutes(a.Fiber)

	a.log.Info().
		Str("resource_type", resourceType).
		Msg("Resource routes registered successfully")
}

// RegisterHealthRoutes registers the core application health check route.
func (s *App) RegisterHealthRoute() {
	s.log.Info().Msg("Registering core application routes")
	s.Fiber.Get("/health", s.healthHandler)
}

func getOperationNames(ops map[resource.Operation]*resource.OperationConfig) []string {
	names := make([]string, 0, len(ops))
	for op, config := range ops {
		if config.Enabled {
			names = append(names, string(op))
		}
	}
	return names
}

// healthHandler is an HTTP handler that responds to health check requests.
func (s *App) healthHandler(c *fiber.Ctx) error {
	s.log.Debug().
		Str("ip", c.IP()).
		Str("method", c.Method()).
		Str("path", c.Path()).
		Msg("Health check requested")

	health := s.Db.Health()

	s.log.Info().
		Interface("status", health).
		Msg("Health check completed")

	return c.JSON(health)
}
