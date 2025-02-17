package state

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
)

// ValidateModel ensures the model type exists and is valid.
// It checks that:
// - Model exists in the Fiber context
// - Model is a pointer to a struct
//
// Parameters:
//   - c: Fiber context containing the model in locals
//
// Returns:
//   - interface{}: Validated model type
//   - error: Validation error with HTTP status code
func validateModel(c *fiber.Ctx) (interface{}, error) {
	modelType := c.Locals("model")
	if modelType == nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "model not found in context")
	}

	modelValue := reflect.ValueOf(modelType)
	if modelValue.Kind() != reflect.Ptr || modelValue.Elem().Kind() != reflect.Struct {
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid model type")
	}

	return modelType, nil
}
