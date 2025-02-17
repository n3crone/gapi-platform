package resource

import (
	"testing"

	"github.com/n3crone/gapi-platform/pkg/state"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// Test models
type TestModel struct {
	ID   uint   `gorm:"primarykey"`
	Name string `json:"name"`
}

type ComplexTestModel struct {
	ID        uint   `gorm:"primarykey"`
	Title     string `json:"title"`
	Reference string `json:"reference"`
}

// Mock DB for testing
type mockDB struct {
	*gorm.DB
}

func newMockDB() *mockDB {
	return &mockDB{}
}

func setupTestEnvironment(t *testing.T) (*ResourceManager, *mockDB, *zerolog.Logger) {
	db := newMockDB()
	logger := zerolog.New(nil)
	rm := NewResourceManager(db.DB, &logger)
	require.NotNil(t, rm, "ResourceManager should be created successfully")
	return rm, db, &logger
}

func TestNewResourceManager(t *testing.T) {
	t.Run("Successfully creates ResourceManager", func(t *testing.T) {
		db := newMockDB()
		logger := zerolog.New(nil)
		rm := NewResourceManager(db.DB, &logger)

		assert.NotNil(t, rm)
		assert.Equal(t, db.DB, rm.DB)
		assert.Equal(t, &logger, rm.logger)
	})
}

func TestCreateResource(t *testing.T) {
	t.Run("Creates resource with pointer model", func(t *testing.T) {
		rm, _, _ := setupTestEnvironment(t)
		model := &TestModel{}
		resource := rm.CreateResource(model)

		assert.NotNil(t, resource)
		assert.Equal(t, "/testmodels", resource.config.Path)
		assert.Equal(t, model, resource.config.Model)
	})

	t.Run("Creates resource with non-pointer model", func(t *testing.T) {
		rm, _, _ := setupTestEnvironment(t)
		model := TestModel{}
		resource := rm.CreateResource(model)

		assert.NotNil(t, resource)
		assert.Equal(t, "/testmodels", resource.config.Path)
	})

	t.Run("Verifies all CRUD operations are configured", func(t *testing.T) {
		rm, _, _ := setupTestEnvironment(t)
		resource := rm.CreateResource(&TestModel{})

		operations := []Operation{
			OperationCreate,
			OperationUpdate,
			OperationGetItem,
			OperationGetList,
			OperationDelete,
		}

		for _, op := range operations {
			opConfig, exists := resource.config.Operations[op]
			assert.True(t, exists, "Operation %v should exist", op)
			assert.True(t, opConfig.Enabled, "Operation %v should be enabled", op)
			assert.NotNil(t, opConfig.Provider, "Provider should be set for operation %v", op)
			assert.NotNil(t, opConfig.Processor, "Processor should be set for operation %v", op)
		}
	})

	t.Run("Applies custom configuration", func(t *testing.T) {
		rm, _, _ := setupTestEnvironment(t)
		customPath := "/custom-path"

		resource := rm.CreateResource(&TestModel{}, func(rc *ResourceConfig) {
			rc.Path = customPath
			rc.Operations[OperationCreate].Enabled = false
		})

		assert.Equal(t, customPath, resource.config.Path)
		assert.False(t, resource.config.Operations[OperationCreate].Enabled)
	})

	t.Run("Handles complex model names correctly", func(t *testing.T) {
		rm, _, _ := setupTestEnvironment(t)
		resource := rm.CreateResource(&ComplexTestModel{})

		assert.Equal(t, "/complextestmodels", resource.config.Path)
	})
}

func TestResourceManagerIntegration(t *testing.T) {
	t.Run("Full resource configuration workflow", func(t *testing.T) {
		rm, _, _ := setupTestEnvironment(t)

		customProvider := &state.DefaultProvider{}
		customProcessor := &state.DefaultProcessor{}

		resource := rm.CreateResource(&TestModel{}, func(rc *ResourceConfig) {
			rc.Path = "/custom-test-models"
			rc.Operations[OperationCreate].Provider = customProvider
			rc.Operations[OperationCreate].Processor = customProcessor
		})

		assert.NotNil(t, resource)
		assert.Equal(t, "/custom-test-models", resource.config.Path)
		assert.Equal(t, customProvider, resource.config.Operations[OperationCreate].Provider)
		assert.Equal(t, customProcessor, resource.config.Operations[OperationCreate].Processor)
	})
}

func BenchmarkCreateResource(b *testing.B) {
	db := newMockDB()
	logger := zerolog.New(nil)
	rm := NewResourceManager(db.DB, &logger)
	model := &TestModel{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rm.CreateResource(model)
	}
}

func BenchmarkCreateResourceWithCustomConfig(b *testing.B) {
	db := newMockDB()
	logger := zerolog.New(nil)
	rm := NewResourceManager(db.DB, &logger)
	model := &TestModel{}
	customConfig := func(rc *ResourceConfig) {
		rc.Path = "/custom-path"
		rc.Operations[OperationCreate].Enabled = false
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rm.CreateResource(model, customConfig)
	}
}
