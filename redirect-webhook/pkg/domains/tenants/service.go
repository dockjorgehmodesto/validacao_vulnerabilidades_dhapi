package tenants

import (
	"fmt"
	"time"

	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"conductor.com.br/pier/hooks/pkg/helpers/splunk"

	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var TenantService Service = nil

func GetConnection(tenantKey *TenantKey, f func(db *gorm.DB) (ex error)) (err error) {
	db, err := TenantService.FetchOnMap(tenantKey)
	if err != nil {
		return err
	}
	return f(db)
}

type Service interface {
	StartupMainTenant()
	StartupIssuerTenants()
	Create(tenant TenantContract) TenantContract
	ListAll() ([]Tenant, error)
	ListAllByPagination(page int, limit int, sort []string) ([]Tenant, error)
	FindById(id string) Tenant
	Replicate(tenantId string, issuerId int) TenantContract
	Update(id string, t *Tenant) Tenant
	Delete(id string)
	FindOneByIssuerId(issuerId int) (TenantContract, error)
	Ping(tenant TenantContract) (err error)
	PingMainTenant() (SQLResult, error)
	FetchOnMap(key *TenantKey) (*gorm.DB, error)
	FetchOnMapByIssuerId(issuerId int) (string, *gorm.DB, error)
	Startup(tenant TenantContract, isIssuerMigration bool) error
	CreateMainDatabaseConnection() (*gorm.DB, error)
	SetupMainDatabase()
	GetMainDatabaseConnection() (*gorm.DB, bool)
	FlushManager()
	LogOpenConnections()
	GetIssuerTenantKey() (*TenantKey, error)
}

func NewService(repository Repository, migrationsManager DatabaseMigrationManager) Service {
	return ServiceImpl{
		repository:        repository,
		migrationsManager: migrationsManager,
		tenantManager:     ApplicationTenantManager,
	}
}

type ServiceImpl struct {
	repository        Repository
	tenantManager     TenantManager
	migrationsManager DatabaseMigrationManager
}

func (s ServiceImpl) CreateMainDatabaseConnection() (*gorm.DB, error) {
	logging.Logger.Info("[TENANT] Creating main tenantKey")

	mainTenant := Tenant{
		IssuerID: 0,
		Dialect:  global.Viper.GetString("database.dialect"),
		Uri:      global.Viper.GetString("database.connectionUri"),
		User:     global.Viper.GetString("database.user"),
		Password: global.Viper.GetString("database.password"),
	}

	return mainTenant.OpenConnection()
}

func (s ServiceImpl) SetupMainDatabase() {
	database, err := s.CreateMainDatabaseConnection()
	if err != nil {
		logging.Logger.Fatal("[TENANT] could not create main database connection", zap.Error(err))
	}

	s.tenantManager.AddByConnection(TENANT_PRINCIPAL_ID, database)
}

func (s ServiceImpl) GetMainDatabaseConnection() (*gorm.DB, bool) {
	return s.tenantManager.Get(TENANT_PRINCIPAL_ID)
}

func (s ServiceImpl) FetchOnMap(key *TenantKey) (*gorm.DB, error) {
	db, exists := s.tenantManager.Get(key.Key)
	if exists {
		return db, nil
	}

	tenant, e := s.FindOneByIssuerId(key.IssuerId)
	if e != nil {
		return nil, errors.Wrap(e, "could not find tenant for issuer")
	}

	_ = s.Startup(tenant, true)
	db, _ = s.tenantManager.Get(key.Key)
	return db, nil
}

func (s ServiceImpl) FetchOnMapByIssuerId(issuerId int) (string, *gorm.DB, error) {
	sample := Tenant{
		IssuerID: uint(issuerId),
		Type:     "READ_AND_WRITE",
	}
	db, exists := s.tenantManager.Get(sample.FormatTenantKey())
	if exists {
		return sample.FormatTenantKey(), db, nil
	}

	tenant, e := s.FindOneByIssuerId(issuerId)
	if e != nil {
		return "", nil, errors.Wrap(e, "could not find tenant for issuer")
	}

	_ = s.Startup(tenant, true)
	db, _ = s.tenantManager.Get(tenant.FormatTenantKey())
	return tenant.FormatTenantKey(), db, nil
}

func (s ServiceImpl) StartupMainTenant() {
	logging.Logger.Info("[TENANT] Starting main tenant")
	s.SetupMainDatabase()
	err := s.migrationsManager.Run(*mainTenantKey, true, true, false)
	if err != nil {
		logging.Logger.Fatal("[TENANT] could not execute migration", zap.Error(err))
	}

	_, e := s.ListAll()
	if e != nil {
		logging.Logger.Fatal("[TENANT] could not list tenants on startup", zap.Error(err))
	}

}

func (s ServiceImpl) StartupIssuerTenants() {
	logging.Logger.Info("[TENANT] Starting issuer tenants")
	tenantSlice, _ := s.ListAll()
	for _, tenant := range tenantSlice {
		s.Startup(tenant, true)
	}
}

func (s ServiceImpl) Startup(tenant TenantContract, isIssuerMigration bool) error {
	id := tenant.FormatTenantKey()
	logging.Logger.Info(fmt.Sprintf("[TENANT] Initializing tenantKey: %s", id))

	e := tenant.PingConnection()
	if e != nil {
		return e
	}

	e = s.tenantManager.Add(tenant)
	if e != nil {
		logging.Logger.Warn(fmt.Sprintf("[TENANT] Could not open connection for tenantKey %s", id), zap.Error(e))
		return e
	}

	tenantKey := *CreateTenantKeyFromIdentifier(id)
	return s.migrationsManager.Run(tenantKey, true, true, isIssuerMigration)
}

func (s ServiceImpl) Create(tenant TenantContract) TenantContract {
	_, e := s.FindOneByIssuerId(int(tenant.GetIssuerId()))
	Error{}.HandleTenantSetupForIssuer(e)

	e = tenant.PingConnection()
	Error{}.HandleUnavailableConnection(e)

	savedTenant, e := s.repository.Create(tenant)
	Error{}.HandleDatabaseErrorWithMessage(e, "error creating tenant")
	return savedTenant
}

func (s ServiceImpl) ListAll() ([]Tenant, error) {
	return s.repository.FindAll()
}

func (s ServiceImpl) ListAllByPagination(page int, limit int, sort []string) ([]Tenant, error) {
	return s.repository.FindAllByPagination(page, limit, sort)
}

func (s ServiceImpl) FindById(id string) Tenant {
	tenant, err := s.repository.FindById(id)
	Error{}.HandleSearchWithMessage(err, "could not find tenant with given id")
	return tenant
}

func (s ServiceImpl) Replicate(tenantId string, issuerId int) TenantContract {
	tenant, err := s.repository.FindById(tenantId)
	Error{}.HandleSearchWithMessage(err, "could not find tenant with given id")

	tenant.ID = ""
	tenant.IssuerID = uint(issuerId)
	tenant.CreatedAt = time.Now()
	tenant.UpdatedAt = time.Now()
	return s.Create(tenant)
}

func (s ServiceImpl) Update(id string, t *Tenant) Tenant {
	s.FindById(id)
	updateErr := s.repository.Update(id, t)
	Error{}.HandleUpdateError(updateErr)

	return s.FindById(id)
}

func (s ServiceImpl) Delete(id string) {
	s.FindById(id)
	e := s.repository.Delete(id)
	Error{}.HandleDatabaseErrorWithMessage(e, "could not delete tenant")
}

func (s ServiceImpl) FindOneByIssuerId(issuerId int) (TenantContract, error) {
	return s.repository.FindByIssuerIdAndType(issuerId, TenantTypeReadWrite)
}

func (s ServiceImpl) Ping(tenant TenantContract) error {
	return tenant.PingConnection()
}

func (s ServiceImpl) PingMainTenant() (SQLResult, error) {
	return s.repository.Ping(mainTenantKey)
}

func (s ServiceImpl) FlushManager() {
	s.tenantManager.Flush()
}

func (s ServiceImpl) LogOpenConnections() {
	uniqueIdentifier, _ := uuid.NewV4()
	loggedIssuerDb := false
	for key, db := range s.tenantManager.GetTenantMap() {
		if key == "TENANT_PRINCIPAL" || !loggedIssuerDb {
			(&splunk.DatabaseConnectionLog{}).InitLog(db.DB().Stats(), key, uniqueIdentifier.String()).Send("open connections")
			if key != "TENANT_PRINCIPAL" {
				loggedIssuerDb = true
			}
		}
	}
}

func (s ServiceImpl) GetIssuerTenantKey() (*TenantKey, error) {
	for key, _ := range s.tenantManager.GetTenantMap() {
		if key != "TENANT_PRINCIPAL" {
			return &TenantKey{Key: key}, nil
		}
	}
	tenants, err := s.ListAll()
	if err != nil {
		return nil, err
	}
	for _, tenant := range tenants {
		err := s.Startup(tenant, true)
		if err != nil {
			return nil, err
		}
		return s.GetIssuerTenantKey()
	}
	return nil, errors.New("could not retrieve issuer tenantKey")
}
