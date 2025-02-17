package state

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type DefaultProcessor struct {
	DB GormDB
}

// Process implements StateProcessor.Process() for GORM-based data manipulation.
// It handles different HTTP methods:
// - POST   -> Create new record
// - PUT    -> Update existing record
// - DELETE -> Remove record
// - GET    -> Validates/transforms output
//
// Parameters:
//   - c: *fiber.Ctx containing the request context
//   - data: Current state data from provider
//
// Returns:
//   - interface{}: Processed result or nil for deletion
//   - error: HTTP-aware error with appropriate status code
func (p *DefaultProcessor) Process(c *fiber.Ctx, data interface{}) (interface{}, error) {
	modelType, err := validateModel(c)
	if err != nil {
		return nil, err
	}

	switch c.Method() {
	case "POST":
		return p.handleCreate(c, modelType)
	case "PUT":
		return p.handleUpdate(c, modelType, data)
	case "DELETE":
		return p.handleDelete(data)
	default:
		return data, nil
	}
}

func (p *DefaultProcessor) handleCreate(c *fiber.Ctx, modelType interface{}) (interface{}, error) {
	newInstance := reflect.New(reflect.ValueOf(modelType).Type().Elem()).Interface()

	if err := c.BodyParser(newInstance); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	result := p.DB.Create(newInstance)
	if result.Error != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to create record")
	}

	return newInstance, nil
}

func (p *DefaultProcessor) handleUpdate(c *fiber.Ctx, modelType interface{}, existing interface{}) (interface{}, error) {
	if existing == nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "record not found")
	}

	// Create new instance for updated data
	newInstance := reflect.New(reflect.ValueOf(modelType).Type().Elem()).Interface()

	if err := c.BodyParser(newInstance); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	// Copy ID from existing record to ensure we update the correct record
	existingValue := reflect.ValueOf(existing).Elem()
	newValue := reflect.ValueOf(newInstance).Elem()
	if idField := existingValue.FieldByName("ID"); idField.IsValid() {
		newValue.FieldByName("ID").Set(idField)
	}

	result := p.DB.Save(newInstance)
	if result.Error != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to update record")
	}

	return newInstance, nil
}

func (p *DefaultProcessor) handleDelete(data interface{}) (interface{}, error) {
	if data == nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "no data to delete")
	}

	result := p.DB.Delete(data)
	if result.Error != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to delete record")
	}

	return nil, nil
}
