package resource

import (
	"github.com/gofiber/fiber/v2"
)

// Resource represents a RESTful API resource that can be registered with
// the application. It encapsulates the resource configuration and provides
// methods for handling HTTP requests.
type Resource struct {
	manager *ResourceManager
	config  ResourceConfig
}

// RegisterRoutes sets up all enabled CRUD operation routes for the resource
// using the provided Fiber application instance.
//
// Parameters:
//   - app: The Fiber application instance for route registration
//
// The following routes are registered if enabled in the configuration:
// - POST   /{path}      -> Create operation
// - PUT    /{path}/:id  -> Update operation
// - DELETE /{path}/:id  -> Delete operation
// - GET    /{path}/:id  -> Get item operation
// - GET    /{path}      -> Get list operation
func (r *Resource) RegisterRoutes(router fiber.Router) {
	path := r.config.Path

	if op, exists := r.config.Operations[OperationGetList]; exists && op.Enabled {
		router.Get(path, r.handleOperation(OperationGetList))
	}

	if op, exists := r.config.Operations[OperationCreate]; exists && op.Enabled {
		router.Post(path, r.handleOperation(OperationCreate))
	}

	if op, exists := r.config.Operations[OperationGetItem]; exists && op.Enabled {
		router.Get(path+"/:id", r.handleOperation(OperationGetItem))
	}

	if op, exists := r.config.Operations[OperationUpdate]; exists && op.Enabled {
		router.Put(path+"/:id", r.handleOperation(OperationUpdate))
	}

	if op, exists := r.config.Operations[OperationDelete]; exists && op.Enabled {
		router.Delete(path+"/:id", r.handleOperation(OperationDelete))
	}
}

// handleOperation creates a Fiber handler function for the specified operation.
// It implements the standard request processing pipeline:
// 1. Validates operation availability
// 2. Sets model context
// 3. Gets initial state from Provider
// 4. Processes state with Processor
// 5. Returns result to client
//
// Parameters:
//   - op: The Operation type to handle (create, update, delete, etc.)
//
// Returns:
//   - fiber.Handler: A handler function that processes the operation
//
// Error Handling:
//   - Returns 404 if operation is not found or disabled
//   - Returns 204 if operation succeeds but has no content
//   - Returns provider/processor errors as-is
func (r *Resource) handleOperation(op Operation) fiber.Handler {
	return func(c *fiber.Ctx) error {
		operationConfig, exists := r.config.Operations[op]
		if !exists || !operationConfig.Enabled {
			return fiber.NewError(fiber.StatusNotFound, "Operation not found")
		}

		// Set model in context
		c.Locals("model", r.config.Model)

		// Get data from provider
		data, err := operationConfig.Provider.Provide(c)
		if err != nil {
			return err
		}

		// Process data
		result, err := operationConfig.Processor.Process(c, data)
		if err != nil {
			return err
		}

		if result == nil {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.JSON(result)
	}
}

// Config returns the resource configuration.
// This method provides read-only access to the resource's configuration.
func (r *Resource) Config() ResourceConfig {
	return r.config
}
