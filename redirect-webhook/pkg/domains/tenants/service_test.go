package tenants

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/labstack/gommon/random"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

var mockMigrationSuccess = MigrationsManagerMockSuccess{}

var mockMigrationError = MigrationsManagerMockError{}

func TestService_FetchOnMap(t *testing.T) {
	t.Run("tenant is not found in memory, then retrieves from database and adds to map", func(t *testing.T) {
		repo := NewRepositoryMock()
		repo.On("FindByIssuerIdAndType", rand.Int(), random.String(10)).Return(TenantMockSuccess{Key: MockTenantKey}, nil)

		service := NewService(repo, mockMigrationSuccess)
		db, e := service.FetchOnMap(&TenantKey{Key: MockTenantKey})
		assert.NotNil(t, db)
		assert.Nil(t, e)
	})

	t.Run("tenant already exists on memory, then repository should not be called", func(t *testing.T) {
		repo := NewRepositoryMock()

		// Forcing an error on repository call, as it should not be called
		repo.On("FindByIssuerIdAndType", rand.Int(), random.String(10)).Return(TenantMockSuccess{Key: MockTenantKey}, nil)
		service := NewService(NewRepositoryMock(), mockMigrationSuccess)

		// By starting up a tenant, it will be added in the service's tenant manager
		startupError := service.Startup(TenantMockSuccess{Key: MockTenantKey}, true)
		assert.Nil(t, startupError)

		db, fetchErr := service.FetchOnMap(&TenantKey{Key: MockTenantKey})
		assert.Nil(t, fetchErr)
		assert.NotNil(t, db)
	})

	t.Run("tenant is not found neither in memory nor in the database", func(t *testing.T) {
		repo := NewRepositoryMock()
		repo.On("FindByIssuerIdAndType", rand.Int(), random.String(10)).Return(TenantMockSuccess{Key: MockTenantKey}, errors.New("mock"))

		service := NewService(repo, mockMigrationSuccess)
		service.FlushManager()
		db, e := service.FetchOnMap(&TenantKey{Key: random.String(15)})
		assert.NotNil(t, e)
		assert.Nil(t, db)
	})
}

func TestService_GetMainDatabaseConnection(t *testing.T) {
	t.Run("return an existing connection", func(t *testing.T) {
		service := NewService(NewRepositoryMock(), mockMigrationSuccess)

		startupError := service.Startup(TenantMockSuccess{Key: TENANT_PRINCIPAL_ID}, true)
		assert.Nil(t, startupError)

		db, connectionExistsInMemory := service.GetMainDatabaseConnection()
		assert.NotNil(t, db)
		assert.True(t, connectionExistsInMemory)
	})

	t.Run("return an invalid connection", func(t *testing.T) {
		service := NewService(NewRepositoryMock(), mockMigrationSuccess)

		service.FlushManager()
		db, connectionExistsInMemory := service.GetMainDatabaseConnection()
		assert.Nil(t, db)
		assert.False(t, connectionExistsInMemory)
	})
}

func TestService_Startup(t *testing.T) {
	t.Run("success", func(t *testing.T) {

		service := NewService(NewRepositoryMock(), mockMigrationSuccess)

		startupError := service.Startup(TenantMockSuccess{Key: MockTenantKey}, true)
		assert.Nil(t, startupError)
	})

	t.Run("error when pinging", func(t *testing.T) {

		service := NewService(NewRepositoryMock(), mockMigrationSuccess)

		startupError := service.Startup(TenantMockError{}, true)
		assert.NotNil(t, startupError)
	})

	t.Run("error when running migration", func(t *testing.T) {

		service := NewService(NewRepositoryMock(), mockMigrationError)

		startupError := service.Startup(TenantMockSuccess{Key: MockTenantKey}, true)
		assert.NotNil(t, startupError)
	})

}

func TestService_Create(t *testing.T) {
	tenantToPersist := Tenant{
		ID:        "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IssuerID:  0,
		Dialect:   random.String(10),
		Uri:       random.String(10),
		User:      random.String(10),
		Password:  random.String(10),
		Type:      random.String(10),
	}

	mockTenantSuccess := TenantMockSuccess{Key: MockTenantKey}
	mockTenantError := TenantMockError{}

	t.Run("success", func(t *testing.T) {
		repo := NewRepositoryMock()

		repo.On("FindByIssuerIdAndType", tenantToPersist, "").Return(mockTenantSuccess, gorm.ErrRecordNotFound)
		repo.On("Create", tenantToPersist).Return(tenantToPersist, nil)

		service := NewService(repo, mockMigrationSuccess)
		createdTenant := service.Create(mockTenantSuccess)
		assert.NotNil(t, createdTenant)
	})

	t.Run("tenant already setup for issuer", func(t *testing.T) {
		repo := NewRepositoryMock()

		// The first validation should fail for this test case
		repo.On("FindByIssuerIdAndType", tenantToPersist, "").Return(nil, errors.New("mock error"))
		repo.On("Create", tenantToPersist).Return(mockTenantSuccess, nil)

		service := NewService(repo, mockMigrationSuccess)
		assert.Panics(t, func() {
			// A valid tenant is passed as argument, as we don't want to crash unexpectedly
			service.Create(mockTenantSuccess)
		})
	})

	t.Run("error when pinging connection", func(t *testing.T) {
		repo := NewRepositoryMock()

		// Validations will pass
		repo.On("FindByIssuerIdAndType", tenantToPersist, "").Return(mockTenantSuccess, nil)
		repo.On("Create", tenantToPersist).Return(mockTenantSuccess, nil)

		service := NewService(repo, mockMigrationSuccess)
		assert.Panics(t, func() {
			// Passing an tenant that will throw an error when pinging
			service.Create(mockTenantError)
		})
	})

	t.Run("error when saving to database", func(t *testing.T) {
		repo := NewRepositoryMock()

		// Validations will pass
		repo.On("FindByIssuerIdAndType", tenantToPersist, "").Return(mockTenantSuccess, nil)

		// Throwing an error when passing data to the repository
		repo.On("Create", tenantToPersist).Return(nil, errors.New("err"))

		service := NewService(repo, mockMigrationSuccess)
		assert.Panics(t, func() {
			// A valid tenant is passed as argument, as we don't want to crash unexpectedly
			service.Create(mockTenantSuccess)
		})
	})
}

func TestService_ListAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := NewRepositoryMock()

		repo.On("FindAll").Return(make([]Tenant, 0), nil)
		service := NewService(repo, mockMigrationSuccess)
		allTenants, e := service.ListAll()
		assert.Nil(t, e)
		assert.NotNil(t, allTenants)
		assert.Empty(t, allTenants)
	})

	t.Run("error", func(t *testing.T) {
		repo := NewRepositoryMock()

		repo.On("FindAll").Return(make([]Tenant, 0), errors.New(""))
		service := NewService(repo, mockMigrationSuccess)
		allTenants, e := service.ListAll()
		assert.Empty(t, allTenants)
		assert.NotNil(t, e)
	})
}

func TestService_ListAllByPagination(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := NewRepositoryMock()

		repo.On("ListAllByPagination", 0, 0, make([]string, 0)).Return(make([]Tenant, 0), nil)
		service := NewService(repo, mockMigrationSuccess)
		allTenants, e := service.ListAllByPagination(0, 0, make([]string, 0))
		assert.Nil(t, e)
		assert.NotNil(t, allTenants)
		assert.Empty(t, allTenants)
	})

	t.Run("error", func(t *testing.T) {
		repo := NewRepositoryMock()

		repo.On("ListAllByPagination", 0, 0, make([]string, 0)).Return(make([]Tenant, 0), errors.New(""))
		service := NewService(repo, mockMigrationSuccess)
		allTenants, e := service.ListAllByPagination(0, 0, make([]string, 0))
		assert.Empty(t, allTenants)
		assert.NotNil(t, e)
	})
}

func TestService_FindByIdd(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := NewRepositoryMock()

		repo.On("FindById", random.String(10)).Return(Tenant{}, nil)
		service := NewService(repo, mockMigrationSuccess)

		tenant := service.FindById(random.String(10))
		assert.NotNil(t, tenant)
	})

	t.Run("not found", func(t *testing.T) {
		repo := NewRepositoryMock()

		repo.On("FindById", random.String(10)).Return(nil, errors.New("mock"))
		service := NewService(repo, mockMigrationSuccess)

		assert.Panics(t, func() {
			service.FindById(random.String(10))
		})
	})
}

func TestService_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := NewRepositoryMock()

		repo.On("FindById", random.String(10)).Return(Tenant{}, nil)
		repo.On("Update", random.String(10), Tenant{}).Return(nil)

		service := NewService(repo, mockMigrationSuccess)
		updatedTenant := service.Update(random.String(10), &Tenant{})
		assert.NotNil(t, updatedTenant)
	})

	t.Run("invalid identifier", func(t *testing.T) {
		repo := NewRepositoryMock()

		repo.On("FindById", random.String(10)).Return(nil, errors.New("mock"))
		repo.On("Update", random.String(10), Tenant{}).Return(nil)

		service := NewService(repo, mockMigrationSuccess)

		assert.Panics(t, func() {
			service.Update(random.String(10), &Tenant{})
		})
	})

	t.Run("database error when updating", func(t *testing.T) {
		repo := NewRepositoryMock()

		repo.On("FindById", random.String(10)).Return(Tenant{}, nil)
		repo.On("Update", random.String(10), Tenant{}).Return(errors.New("mock"))

		service := NewService(repo, mockMigrationSuccess)

		assert.Panics(t, func() {
			service.Update(random.String(10), &Tenant{})
		})
	})

}

func TestService_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := NewRepositoryMock()
		repo.On("FindById", random.String(10)).Return(Tenant{}, nil)
		repo.On("Delete", random.String(10)).Return(nil)

		service := NewService(repo, mockMigrationSuccess)
		assert.NotPanics(t, func() {
			service.Delete(random.String(10))
		})
	})

	t.Run("does not find tenant by its id", func(t *testing.T) {
		repo := NewRepositoryMock()
		repo.On("FindById", random.String(10)).Return(nil, errors.New("mock"))
		repo.On("Delete", random.String(10)).Return(nil)

		service := NewService(repo, mockMigrationSuccess)
		assert.Panics(t, func() {
			service.Delete(random.String(10))
		})
	})

	t.Run("database error when deleting", func(t *testing.T) {
		repo := NewRepositoryMock()
		repo.On("FindById", random.String(10)).Return(Tenant{}, nil)
		repo.On("Delete", random.String(10)).Return(errors.New("mock"))

		service := NewService(repo, mockMigrationSuccess)
		assert.Panics(t, func() {
			service.Delete(random.String(10))
		})
	})
}

func TestService_FindOneByIssuerId(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := NewRepositoryMock()
		repo.On("FindByIssuerIdAndType", rand.Int(), random.String(10)).Return(Tenant{}, nil)

		service := NewService(repo, mockMigrationSuccess)
		tenantContract, e := service.FindOneByIssuerId(rand.Int())
		assert.Nil(t, e)
		assert.NotNil(t, tenantContract)
	})

	t.Run("error", func(t *testing.T) {
		repo := NewRepositoryMock()
		repo.On("FindByIssuerIdAndType", rand.Int(), random.String(10)).Return(Tenant{}, errors.New("mock"))

		service := NewService(repo, mockMigrationSuccess)
		_, e := service.FindOneByIssuerId(rand.Int())
		assert.NotNil(t, e)
	})
}

func TestService_Ping(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := NewService(NewRepositoryMock(), mockMigrationSuccess)
		err := service.Ping(TenantMockSuccess{Key: MockTenantKey})
		assert.Nil(t, err)
	})

	t.Run("error", func(t *testing.T) {
		service := NewService(NewRepositoryMock(), mockMigrationSuccess)
		err := service.Ping(TenantMockError{})
		assert.NotNil(t, err)
	})
}

func TestService_PingMainTenant(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := NewRepositoryMock()
		repo.On("Ping", MockTenantKey).Return(SQLResult{}, nil)

		service := NewService(repo, mockMigrationSuccess)
		result, e := service.PingMainTenant()
		assert.NotNil(t, result)
		assert.Nil(t, e)
	})

	t.Run("error", func(t *testing.T) {
		repo := NewRepositoryMock()
		repo.On("Ping", MockTenantKey).Return(SQLResult{}, errors.New("mock"))

		service := NewService(repo, mockMigrationSuccess)
		_, e := service.PingMainTenant()
		assert.NotNil(t, e)
	})

}
