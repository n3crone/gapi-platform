package resource

import "github.com/gofiber/fiber/v2"

// ResourceConfig defines the configuration for an API resource.
// It specifies the data model, available operations, and base path
// for the resource endpoints.
type ResourceConfig struct {
	Model      interface{}                    // The data model struct for this resource
	Operations map[Operation]*OperationConfig // Available CRUD operations and their configurations
	Path       string                         // Base URL path for the resource
}

// Operation represents a CRUD operation type.
// It is used as a key in the Operations map to configure
// different aspects of each operation.
type Operation string

// Standard CRUD operations supported by the resource system
const (
	OperationCreate  Operation = "create"   // Create new resource instance (POST)
	OperationUpdate  Operation = "update"   // Update existing resource (PUT)
	OperationDelete  Operation = "delete"   // Delete resource instance (DELETE)
	OperationGetItem Operation = "get_item" // Retrieve single resource (GET with ID)
	OperationGetList Operation = "get_list" // Retrieve list of resources (GET)
)

// OperationConfig defines the behavior of a specific CRUD operation
// by configuring its state management and processing pipeline.
type OperationConfig struct {
	Provider  StateProvider  // Responsible for fetching data from database
	Processor StateProcessor // Handles state transformation and business logic
	Enabled   bool           // Whether this operation is available
}

// StateProvider defines the interface for preparing initial state
// before operation processing. It retrieves data from database
type StateProvider interface {
	Provide(c *fiber.Ctx) (interface{}, error)
}

// StateProcessor defines the interface for processing operation state.
// It handles business logic, data transformation, and validation
// before completing the operation.
type StateProcessor interface {
	Process(c *fiber.Ctx, data interface{}) (interface{}, error)
}
