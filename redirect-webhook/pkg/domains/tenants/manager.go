package tenants

import (
	"conductor.com.br/pier/hooks/pkg/conf/global"
	"github.com/jinzhu/gorm"
	"sync"
)

var ApplicationTenantManager = NewTenantManager()

type TenantManager interface {
	Get(key string) (*gorm.DB, bool)
	GetTenantMap() map[string]*gorm.DB
	Add(tenant TenantContract) error
	AddByConnection(key string, db *gorm.DB)
	Flush()
}

func NewTenantManager() TenantManager {
	return &TenantManagerImpl{connections: make(map[string]*gorm.DB), mutex: sync.Mutex{}}
}

type TenantManagerImpl struct {
	connections map[string]*gorm.DB
	mutex       sync.Mutex
}

func (t *TenantManagerImpl) Get(key string) (*gorm.DB, bool) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if value, ok := t.connections[key]; ok {
		return value, true
	}
	return nil, false
}

func (t *TenantManagerImpl) GetTenantMap() map[string]*gorm.DB {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	resultMap := make(map[string]*gorm.DB)
	for k, v := range t.connections {
		resultMap[k] = v
	}
	return resultMap
}

func (t *TenantManagerImpl) Add(tenant TenantContract) error {
	if global.Viper.GetBool("database.reuseConnectionsOnMap") {
		if len(t.connections) > 1 {
			for tenantKey, db := range t.connections {
				if tenantKey != TENANT_PRINCIPAL_ID {
					t.AddByConnection(tenant.GetTenantKey().Key, db)
					return nil
				}
			}
		}
	}

	db, e := tenant.OpenConnection()
	if e != nil {
		return e
	}
	t.AddByConnection(tenant.GetTenantKey().Key, db)
	return nil
}

func (t *TenantManagerImpl) AddByConnection(key string, db *gorm.DB) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.connections[key] = db
}

func (t *TenantManagerImpl) Flush() {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.connections = nil
	t.connections = make(map[string]*gorm.DB)
}
