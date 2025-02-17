package resource

// Registrable defines the interface that API resources must implement
// to be registered with the application. It provides a contract for
// creating resource instances with proper configuration.
//
// This interface enables:
// - Automatic resource registration with the application
// - Consistent resource configuration across different resource types
// - Flexible resource customization through the ResourceManager
type Registrable interface {
	// CreateResource initializes a new resource instance with the given
	// resource manager. This method should configure the resource with
	// appropriate CRUD operations and custom handlers if needed.
	//
	// Parameters:
	//   - rm: The ResourceManager instance to use for resource creation
	//
	// Returns:
	//   - *Resource: A fully configured resource ready for route registration
	CreateResource(*ResourceManager) *Resource
}
