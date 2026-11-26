package tenants

import (
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestManager_NewInstance(t *testing.T) {

	t.Run("is not null", func(t *testing.T) {
		manager := NewTenantManager()
		assert.NotNil(t, manager)
	})

	t.Run("contains a map", func(t *testing.T) {
		manager := NewTenantManager()
		tenantMap := manager.GetTenantMap()
		assert.NotNil(t, tenantMap)
	})

	t.Run("map is empty", func(t *testing.T) {
		manager := NewTenantManager()
		tenantMap := manager.GetTenantMap()
		assert.Empty(t, tenantMap)
	})
}

func TestManager_Get(t *testing.T) {
	t.Run("should return nil and false when key is not present", func(t *testing.T) {
		manager := NewTenantManager()
		db, b := manager.Get("invalidkey")

		assert.False(t, b)
		assert.Nil(t, db)
	})

	t.Run("should receive the connection and no errors", func(t *testing.T) {
		manager := NewTenantManager()
		tenant := TenantMockSuccess{Key: MockTenantKey}
		addError := manager.Add(tenant)
		assert.Nil(t, addError)

		db, exists := manager.Get(MockTenantKey)
		assert.True(t, exists)
		assert.NotNil(t, db)
	})
}

func TestManager_GetTenantMap(t *testing.T) {
	t.Run("returns all values", func(t *testing.T) {
		manager := NewTenantManager()
		tenant := TenantMockSuccess{}
		addError := manager.Add(tenant)
		assert.Nil(t, addError)

		tenantMap := manager.GetTenantMap()
		assert.NotEmpty(t, tenantMap)
		assert.Equal(t, 1, len(tenantMap))
	})

	t.Run("changing the returned map does not influence on the actual map", func(t *testing.T) {
		manager := NewTenantManager()
		tenant := TenantMockSuccess{Key: MockTenantKey}
		addError := manager.Add(tenant)
		assert.Nil(t, addError)

		mockDb := gorm.DB{}
		tenantMap := manager.GetTenantMap()
		tenantMap[MockTenantKey] = &mockDb
		db, exists := manager.Get(MockTenantKey)
		assert.True(t, exists)
		assert.NotNil(t, db)
		assert.NotEqual(t, db, mockDb)

		delete(tenantMap, MockTenantKey)
		_, exists = manager.Get(MockTenantKey)
		assert.True(t, exists)
	})
}

func TestManager_Add(t *testing.T) {
	t.Run("adding invalid tenant or connection", func(t *testing.T) {
		manager := NewTenantManager()
		tenant := TenantMockError{}
		addError := manager.Add(tenant)
		assert.NotNil(t, addError)
	})

	t.Run("adding valid tenant or connection", func(t *testing.T) {
		manager := NewTenantManager()
		tenant := TenantMockSuccess{}
		addError := manager.Add(tenant)
		assert.Nil(t, addError)
	})
}

func TestManager_AddByConnection(t *testing.T) {
	t.Run("adding a connection increments the manager map", func(t *testing.T) {
		manager := NewTenantManager()
		tenant := TenantMockSuccess{}
		db, e := tenant.OpenConnection()
		assert.Nil(t, e)
		assert.NotNil(t, db)
		assert.Equal(t, 0, len(manager.GetTenantMap()))
		manager.AddByConnection(MockTenantKey, db)

		assert.Equal(t, 1, len(manager.GetTenantMap()))
	})

	t.Run("key is available after adding it to manager map", func(t *testing.T) {
		manager := NewTenantManager()
		tenant := TenantMockSuccess{}
		db, e := tenant.OpenConnection()
		assert.Nil(t, e)
		assert.NotNil(t, db)
		manager.AddByConnection(MockTenantKey, db)

		connection, exists := manager.Get(MockTenantKey)
		assert.True(t, exists)
		assert.NotNil(t, connection)
	})

	t.Run("writing same key overrides the map", func(t *testing.T) {
		manager := NewTenantManager()
		tenant := TenantMockSuccess{}
		db, e := tenant.OpenConnection()
		assert.Nil(t, e)
		assert.NotNil(t, db)
		manager.AddByConnection(MockTenantKey, db)

		assert.Equal(t, 1, len(manager.GetTenantMap()))

		mockDb := &gorm.DB{}
		manager.AddByConnection(MockTenantKey, mockDb)
		assert.Equal(t, 1, len(manager.GetTenantMap()))

		connection, exists := manager.Get(MockTenantKey)
		assert.True(t, exists)
		assert.Equal(t, connection, mockDb)
	})
}
