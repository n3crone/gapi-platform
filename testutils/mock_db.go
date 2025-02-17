package testutils

import (
	"reflect"

	"gorm.io/gorm"
)

type MockDB struct {
	FindByIDError error
	FindAllError  error
	CreateError   error
	UpdateError   error
	DeleteError   error
	Records       []interface{}
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	if m.CreateError != nil {
		return &gorm.DB{Error: m.CreateError}
	}
	// Simulate ID assignment
	reflect.ValueOf(value).Elem().FieldByName("ID").SetUint(1)
	return &gorm.DB{}
}

func (m *MockDB) Save(value interface{}) *gorm.DB {
	if m.UpdateError != nil {
		return &gorm.DB{Error: m.UpdateError}
	}
	return &gorm.DB{}
}

func (m *MockDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	if m.DeleteError != nil {
		return &gorm.DB{Error: m.DeleteError}
	}
	return &gorm.DB{}
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	if m.FindByIDError != nil {
		return &gorm.DB{Error: m.FindByIDError}
	}

	if len(m.Records) > 0 {
		val := reflect.ValueOf(dest).Elem()
		record := reflect.ValueOf(m.Records[0])
		copyFields(val, record.Elem())
	}

	return &gorm.DB{Error: nil}
}

func (m *MockDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	if m.FindAllError != nil {
		return &gorm.DB{Error: m.FindAllError}
	}

	val := reflect.ValueOf(dest).Elem()
	for _, record := range m.Records {
		newElem := reflect.New(val.Type().Elem()).Elem()
		copyFields(newElem, reflect.ValueOf(record).Elem())
		val.Set(reflect.Append(val, newElem))
	}

	return &gorm.DB{Error: nil}
}

func copyFields(dest, src reflect.Value) {
	for i := 0; i < src.NumField(); i++ {
		destField := dest.FieldByName(src.Type().Field(i).Name)
		if destField.IsValid() && destField.CanSet() {
			destField.Set(src.Field(i))
		}
	}
}
