package state

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// DefaultProvider implements the StateProvider interface for GORM database operations.
// It provides data retrieval functionality for REST API resources by:
// - Handling both single and collection queries
// - Supporting dynamic model types
// - Managing database connections via GORM
type DefaultProvider struct {
	DB GormDB
}

type GormDB interface {
	First(dest interface{}, conds ...interface{}) *gorm.DB
	Find(dest interface{}, conds ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	Save(value interface{}) *gorm.DB
	Delete(value interface{}, conds ...interface{}) *gorm.DB
}

// Provide implements StateProvider.Provide() for GORM-based data retrieval.
// It determines the appropriate query type based on URL parameters:
// - GET /{resource}/:id -> Single item lookup
// - GET /{resource}     -> Collection lookup
//
// Parameters:
//   - c: *fiber.Ctx containing the request context and model information
//
// Returns:
//   - interface{}: Retrieved item(s) or nil if not found
//   - error: HTTP-aware error with appropriate status code
func (p *DefaultProvider) Provide(c *fiber.Ctx) (interface{}, error) {
	modelType, err := validateModel(c)
	if err != nil {
		return nil, err
	}

	if id := c.Params("id"); id != "" {
		return p.findById(id, modelType)
	}

	return p.findAll(modelType)
}

// findById retrieves a single record by ID
func (p *DefaultProvider) findById(id string, modelType interface{}) (interface{}, error) {
	result := p.DB.First(modelType, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, fiber.NewError(fiber.StatusNotFound, "record not found")
		}
		return nil, fiber.NewError(fiber.StatusInternalServerError, "database error")
	}

	return modelType, nil
}

// findAll retrieves all records of the given model type
func (p *DefaultProvider) findAll(modelType interface{}) (interface{}, error) {
	modelValue := reflect.ValueOf(modelType)
	results := reflect.New(reflect.SliceOf(modelValue.Type().Elem())).Interface()

	result := p.DB.Find(results)
	if result.Error != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to fetch records")
	}

	return results, nil
}
