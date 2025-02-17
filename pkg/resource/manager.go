package resource

import (
	"reflect"
	"strings"

	"github.com/n3crone/gapi-platform/pkg/state"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

// ResourceManager handles the creation and management of API resources.
// It provides a centralized way to create and configure resources with
// their associated CRUD operations.
type ResourceManager struct {
	DB     *gorm.DB
	logger *zerolog.Logger
}

// NewResourceManager creates a new instance of ResourceManager with the provided
// database connection. This manager will be used to create and configure
// API resources with default or custom configurations.
//
// Parameters:
//   - db: A GORM database instance for database operations
//
// Returns:
//   - *ResourceManager: A new resource manager instance
func NewResourceManager(db *gorm.DB, logger *zerolog.Logger) *ResourceManager {
	return &ResourceManager{DB: db, logger: logger}
}

// CreateResource creates a new API resource with the given model and optional
// custom configurations. It automatically sets up default CRUD operations
// and allows customization through functional options.
//
// Parameters:
//   - model: The data model struct for the resource
//   - customConfig: Optional functional options to customize resource configuration
//
// Features:
//   - Automatically generates API endpoint paths based on model name
//   - Sets up default CRUD operations with standard providers and processors
//   - Allows custom configuration of any aspect of the resource
//
// Returns:
//   - *Resource: A configured resource instance ready for route registration
func (rm *ResourceManager) CreateResource(model interface{}, customConfig ...func(*ResourceConfig)) *Resource {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	defaultPath := "/" + strings.ToLower(modelType.Name()) + "s"

	// Initialize default resource configuration with all CRUD operations
	config := ResourceConfig{
		Model: model,
		Path:  defaultPath,
		Operations: map[Operation]*OperationConfig{
			OperationCreate: {
				Provider:  &state.DefaultProvider{DB: rm.DB},
				Processor: &state.DefaultProcessor{DB: rm.DB},
				Enabled:   true,
			},
			OperationUpdate: {
				Provider:  &state.DefaultProvider{DB: rm.DB},
				Processor: &state.DefaultProcessor{DB: rm.DB},
				Enabled:   true,
			},
			OperationGetItem: {
				Provider:  &state.DefaultProvider{DB: rm.DB},
				Processor: &state.DefaultProcessor{DB: rm.DB},
				Enabled:   true,
			},
			OperationGetList: {
				Provider:  &state.DefaultProvider{DB: rm.DB},
				Processor: &state.DefaultProcessor{DB: rm.DB},
				Enabled:   true,
			},
			OperationDelete: {
				Provider:  &state.DefaultProvider{DB: rm.DB},
				Processor: &state.DefaultProcessor{DB: rm.DB},
				Enabled:   true,
			},
		},
	}

	// Apply any custom configurations provided
	for _, customizer := range customConfig {
		customizer(&config)
	}

	return &Resource{
		manager: rm,
		config:  config,
	}
}
