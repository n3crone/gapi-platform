package state

import (
	"net/http/httptest"
	"testing"

	"github.com/n3crone/gapi-platform/testutils"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestProvider(_ *testing.T) (*DefaultProvider, *testutils.MockDB, *fiber.App) {
	mockDB := &testutils.MockDB{
		Records: []interface{}{
			&TestModel{ID: 1, Name: "Test 1"},
			&TestModel{ID: 2, Name: "Test 2"},
		},
	}

	provider := &DefaultProvider{DB: mockDB}
	app := fiber.New()

	return provider, mockDB, app
}

func TestProvide(t *testing.T) {
	t.Run("Get single record by ID successfully", func(t *testing.T) {
		provider, _, app := setupTestProvider(t)

		app.Get("/:id", func(c *fiber.Ctx) error {
			c.Locals("model", &TestModel{})
			data, err := provider.Provide(c)
			require.NoError(t, err)

			model, ok := data.(*TestModel)
			require.True(t, ok)
			assert.Equal(t, uint(1), model.ID)
			assert.Equal(t, "Test 1", model.Name)

			return nil
		})

		req := httptest.NewRequest("GET", "/1", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Get all records successfully", func(t *testing.T) {
		provider, _, app := setupTestProvider(t)

		app.Get("/", func(c *fiber.Ctx) error {
			c.Locals("model", &TestModel{})
			data, err := provider.Provide(c)
			require.NoError(t, err)

			models, ok := data.(*[]TestModel)
			require.True(t, ok)
			assert.Len(t, *models, 2)

			return nil
		})

		req := httptest.NewRequest("GET", "/", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Record not found error", func(t *testing.T) {
		provider, mockDB, app := setupTestProvider(t)
		mockDB.FindByIDError = gorm.ErrRecordNotFound

		app.Get("/:id", func(c *fiber.Ctx) error {
			c.Locals("model", &TestModel{})
			_, err := provider.Provide(c)
			assert.Error(t, err)

			fiberErr, ok := err.(*fiber.Error)
			require.True(t, ok)
			assert.Equal(t, fiber.StatusNotFound, fiberErr.Code)

			return err
		})

		req := httptest.NewRequest("GET", "/999", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	})

	t.Run("Database error", func(t *testing.T) {
		provider, mockDB, app := setupTestProvider(t)
		mockDB.FindAllError = gorm.ErrInvalidTransaction

		app.Get("/", func(c *fiber.Ctx) error {
			c.Locals("model", &TestModel{})
			_, err := provider.Provide(c)
			assert.Error(t, err)

			fiberErr, ok := err.(*fiber.Error)
			require.True(t, ok)
			assert.Equal(t, fiber.StatusInternalServerError, fiberErr.Code)

			return err
		})

		req := httptest.NewRequest("GET", "/", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Invalid model in context", func(t *testing.T) {
		provider, _, app := setupTestProvider(t)

		app.Get("/", func(c *fiber.Ctx) error {
			// Don't set model in context
			_, err := provider.Provide(c)
			assert.Error(t, err)

			fiberErr, ok := err.(*fiber.Error)
			require.True(t, ok)
			assert.Equal(t, fiber.StatusBadRequest, fiberErr.Code)

			return err
		})

		req := httptest.NewRequest("GET", "/", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}
