package resource

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterRoutes(t *testing.T) {
	t.Run("All operations enabled", func(t *testing.T) {
		// Setup
		app := fiber.New()
		resource := createTestResource("/api/test", map[Operation]bool{
			OperationGetList: true,
			OperationCreate:  true,
			OperationGetItem: true,
			OperationUpdate:  true,
			OperationDelete:  true,
		})

		resource.RegisterRoutes(app)

		// Test GET /api/test (List)
		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/api/test", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// Test POST /api/test (Create)
		resp, err = app.Test(httptest.NewRequest(http.MethodPost, "/api/test", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// Test GET /api/test/123 (GetItem)
		resp, err = app.Test(httptest.NewRequest(http.MethodGet, "/api/test/123", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// Test PUT /api/test/123 (Update)
		resp, err = app.Test(httptest.NewRequest(http.MethodPut, "/api/test/123", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		// Test DELETE /api/test/123 (Delete)
		resp, err = app.Test(httptest.NewRequest(http.MethodDelete, "/api/test/123", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
	})

	t.Run("All operations disabled", func(t *testing.T) {
		app := fiber.New()
		resource := createTestResource("/api/test", map[Operation]bool{
			OperationGetList: false,
			OperationCreate:  false,
			OperationGetItem: false,
			OperationUpdate:  false,
			OperationDelete:  false,
		})

		resource.RegisterRoutes(app)

		// Test disabled endpoints return 404
		endpoints := []struct {
			method string
			path   string
		}{
			{http.MethodGet, "/api/test"},
			{http.MethodPost, "/api/test"},
			{http.MethodGet, "/api/test/123"},
			{http.MethodPut, "/api/test/123"},
			{http.MethodDelete, "/api/test/123"},
		}

		for _, e := range endpoints {
			resp, err := app.Test(httptest.NewRequest(e.method, e.path, nil))
			require.NoError(t, err)
			assert.Equal(t, fiber.StatusNotFound, resp.StatusCode,
				"Expected 404 for %s %s", e.method, e.path)
		}
	})

	t.Run("Provider returns data", func(t *testing.T) {
		app := fiber.New()
		resource := createTestResource("/api/test", map[Operation]bool{
			OperationGetList: true,
		})

		testData := []map[string]interface{}{
			{"id": "11", "name": "test item"},
		}
		resource.config.Operations[OperationGetList].Provider = &mockProvider{
			response: testData,
		}
		resource.config.Operations[OperationGetList].Processor = &mockProcessor{
			response: testData,
		}

		resource.RegisterRoutes(app)

		resp, err := app.Test(httptest.NewRequest(http.MethodGet, "/api/test", nil))
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var result []map[string]interface{}
		body, _ := io.ReadAll(resp.Body)
		require.NoError(t, json.Unmarshal(body, &result))
		assert.Equal(t, testData, result)
	})

	t.Run("Operation responses match expected types", func(t *testing.T) {
		app := fiber.New()
		resource := createTestResource("/api/test", map[Operation]bool{
			OperationGetList: true,
			OperationCreate:  true,
			OperationGetItem: true,
			OperationUpdate:  true,
			OperationDelete:  true,
		})

		resource.RegisterRoutes(app)

		tests := []struct {
			name       string
			method     string
			path       string
			operation  Operation
			wantStatus int
			assertBody func(t *testing.T, body []byte)
		}{
			{
				name:       "GetList returns array",
				method:     http.MethodGet,
				path:       "/api/test",
				operation:  OperationGetList,
				wantStatus: fiber.StatusOK,
				assertBody: func(t *testing.T, body []byte) {
					var result []map[string]interface{}
					require.NoError(t, json.Unmarshal(body, &result))
					assert.Len(t, result, 1)
					assert.Equal(t, "1", result[0]["id"])
				},
			},
			{
				name:       "Create returns object",
				method:     http.MethodPost,
				path:       "/api/test",
				operation:  OperationCreate,
				wantStatus: fiber.StatusOK,
				assertBody: func(t *testing.T, body []byte) {
					var result map[string]interface{}
					require.NoError(t, json.Unmarshal(body, &result))
					assert.Equal(t, "2", result["id"])
				},
			},
			{
				name:       "GetItem returns object",
				method:     http.MethodGet,
				path:       "/api/test/123",
				operation:  OperationGetItem,
				wantStatus: fiber.StatusOK,
				assertBody: func(t *testing.T, body []byte) {
					var result map[string]interface{}
					require.NoError(t, json.Unmarshal(body, &result))
					assert.Equal(t, "3", result["id"])
				},
			},
			{
				name:       "Delete returns no content",
				method:     http.MethodDelete,
				path:       "/api/test/123",
				operation:  OperationDelete,
				wantStatus: fiber.StatusNoContent,
				assertBody: func(t *testing.T, body []byte) {
					assert.Empty(t, body)
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				resp, err := app.Test(httptest.NewRequest(tt.method, tt.path, nil))
				require.NoError(t, err)
				assert.Equal(t, tt.wantStatus, resp.StatusCode)

				if tt.assertBody != nil {
					body, _ := io.ReadAll(resp.Body)
					tt.assertBody(t, body)
				}
			})
		}
	})
}

// Enhanced mock implementations
type mockProvider struct {
	response interface{}
	err      error
}

func (m *mockProvider) Provide(c *fiber.Ctx) (interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

type mockProcessor struct {
	response interface{}
	err      error
}

func (m *mockProcessor) Process(c *fiber.Ctx, data interface{}) (interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

// createTestResource creates a Resource instance with specified operations enabled/disabled
func createTestResource(path string, operations map[Operation]bool) *Resource {
	config := ResourceConfig{
		Path:       path,
		Operations: make(map[Operation]*OperationConfig),
	}

	// Setup operations with default responses
	defaultResponses := map[Operation]interface{}{
		OperationGetList: []map[string]interface{}{{"id": "1", "name": "test"}},
		OperationCreate:  map[string]interface{}{"id": "2", "name": "created"},
		OperationGetItem: map[string]interface{}{"id": "3", "name": "item"},
		OperationUpdate:  map[string]interface{}{"id": "4", "name": "updated"},
		OperationDelete:  nil, // Delete should return 204
	}

	for op, enabled := range operations {
		config.Operations[op] = &OperationConfig{
			Enabled:  enabled,
			Provider: &mockProvider{},
			Processor: &mockProcessor{
				response: defaultResponses[op],
			},
		}
	}

	return &Resource{
		config: config,
	}
}
