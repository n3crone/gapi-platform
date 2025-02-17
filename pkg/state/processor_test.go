package state

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/n3crone/gapi-platform/testutils"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type TestModel struct {
	ID   uint   `json:"id" gorm:"primarykey"`
	Name string `json:"name"`
}

func setupTestProcessor(_ *testing.T) (*DefaultProcessor, *testutils.MockDB, *fiber.App) {
	mockDB := &testutils.MockDB{
		Records: []interface{}{
			&TestModel{ID: 1, Name: "Test 1"},
			&TestModel{ID: 2, Name: "Test 2"},
		},
	}

	processor := &DefaultProcessor{DB: mockDB}
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if ferr, ok := err.(*fiber.Error); ok {
				return c.Status(ferr.Code).JSON(fiber.Map{
					"error": ferr.Message,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	return processor, mockDB, app
}

func TestProcess(t *testing.T) {
	t.Run("Create operation successfully", func(t *testing.T) {
		processor, _, app := setupTestProcessor(t)

		app.Post("/test", func(c *fiber.Ctx) error {
			c.Locals("model", &TestModel{})
			data, err := processor.Process(c, nil)
			require.NoError(t, err)

			model, ok := data.(*TestModel)
			require.True(t, ok)
			assert.NotZero(t, model.ID)
			assert.Equal(t, "test item", model.Name)

			return c.JSON(data)
		})

		payload := `{"name":"test item"}`
		req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Update operation successfully", func(t *testing.T) {
		processor, _, app := setupTestProcessor(t)

		app.Put("/test/:id", func(c *fiber.Ctx) error {
			c.Locals("model", &TestModel{})
			data, err := processor.Process(c, &TestModel{ID: 1})
			require.NoError(t, err)

			model, ok := data.(*TestModel)
			require.True(t, ok)
			assert.Equal(t, uint(1), model.ID)
			assert.Equal(t, "updated name", model.Name)

			return c.JSON(data)
		})

		payload := `{"name":"updated name"}`
		req := httptest.NewRequest("PUT", "/test/1", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})

	t.Run("Delete operation successfully", func(t *testing.T) {
		processor, _, app := setupTestProcessor(t)

		app.Delete("/test/:id", func(c *fiber.Ctx) error {
			c.Locals("model", &TestModel{})
			result, err := processor.Process(c, &TestModel{ID: 1})
			require.NoError(t, err)
			assert.Nil(t, result)

			return c.SendStatus(fiber.StatusNoContent)
		})

		req := httptest.NewRequest("DELETE", "/test/1", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})

	t.Run("Database errors", func(t *testing.T) {
		processor, mockDB, app := setupTestProcessor(t)
		mockDB.CreateError = gorm.ErrInvalidTransaction
		mockDB.UpdateError = gorm.ErrRecordNotFound
		mockDB.DeleteError = gorm.ErrInvalidData

		tests := []struct {
			name     string
			method   string
			path     string
			payload  string
			wantCode int
		}{
			{
				name:     "Create error",
				method:   "POST",
				path:     "/test",
				payload:  `{"name":"test"}`,
				wantCode: fiber.StatusInternalServerError,
			},
			{
				name:     "Update error",
				method:   "PUT",
				path:     "/test/1",
				payload:  `{"name":"test"}`,
				wantCode: fiber.StatusInternalServerError,
			},
			{
				name:     "Delete error",
				method:   "DELETE",
				path:     "/test/1",
				wantCode: fiber.StatusInternalServerError,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var req *http.Request
				if tt.payload != "" {
					req = httptest.NewRequest(tt.method, tt.path, bytes.NewBufferString(tt.payload))
					req.Header.Set("Content-Type", "application/json")
				} else {
					req = httptest.NewRequest(tt.method, tt.path, nil)
				}

				app.Add(tt.method, tt.path, func(c *fiber.Ctx) error {
					c.Locals("model", &TestModel{})
					_, err := processor.Process(c, &TestModel{ID: 1})
					return err
				})

				resp, err := app.Test(req)
				require.NoError(t, err)
				assert.Equal(t, tt.wantCode, resp.StatusCode)
			})
		}
	})

	t.Run("Invalid model in context", func(t *testing.T) {
		processor, _, app := setupTestProcessor(t)

		app.Post("/test", func(c *fiber.Ctx) error {
			// Don't set model in context
			_, err := processor.Process(c, nil)
			return err
		})

		req := httptest.NewRequest("POST", "/test", bytes.NewBufferString(`{"name":"test"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})
}
