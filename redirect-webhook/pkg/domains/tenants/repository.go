package tenants

import (
	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"github.com/jinzhu/gorm"
)

var mainTenantKey = NewMainTenantKey()

type Repository interface {
	Create(t TenantContract) (Tenant, error)
	Delete(id string) error
	Update(id string, t *Tenant) error
	FindById(id string) (Tenant, error)
	FindByIssuerIdAndType(issuerId int, accessType string) (TenantContract, error)
	FindAll() ([]Tenant, error)
	FindAllByPagination(page int, limit int, orderBy []string) ([]Tenant, error)
	Ping(tenantKey *TenantKey) (SQLResult, error)
}

type SQLResult struct {
	Result string
}

type DatabaseRepository struct{}

func NewDatabaseRepository() Repository {
	return DatabaseRepository{}
}

func (d DatabaseRepository) FindByIssuerIdAndType(issuerId int, accessType string) (TenantContract, error) {
	tenant := Tenant{}
	err := GetConnection(mainTenantKey, func(db *gorm.DB) (ex error) {
		return db.Where("issuer_id = ? and type = ?", issuerId, accessType).First(&tenant).Error
	})
	return tenant, err
}

func (d DatabaseRepository) FindAllByPagination(page int, limit int, orderBy []string) ([]Tenant, error) {
	tenantSlice := make([]Tenant, 0)
	err := GetConnection(mainTenantKey, func(db *gorm.DB) (ex error) {
		helpers.Paging(&helpers.Param{
			DB:      db,
			Page:    page,
			Limit:   limit,
			OrderBy: orderBy,
			ShowSQL: global.ShowSQL,
		}, &tenantSlice)

		return nil
	})

	return tenantSlice, err
}

func (d DatabaseRepository) Create(t TenantContract) (Tenant, error) {
	tenantToSave := t.(Tenant)
	err := GetConnection(mainTenantKey, func(db *gorm.DB) (ex error) {
		return db.Create(&tenantToSave).Error
	})
	return tenantToSave, err
}

func (d DatabaseRepository) Delete(id string) error {
	return GetConnection(mainTenantKey, func(db *gorm.DB) (ex error) {
		return db.Where("id = ? ", id).Delete(&Tenant{}).Error
	})
}

func (d DatabaseRepository) Update(id string, t *Tenant) error {
	return GetConnection(mainTenantKey, func(db *gorm.DB) (ex error) {
		return db.Where("id = ? ", id).Update(t).Error
	})
}

func (d DatabaseRepository) FindById(id string) (Tenant, error) {
	var t = Tenant{}
	err := GetConnection(mainTenantKey, func(db *gorm.DB) (ex error) {
		return db.Where("id = ? ", id).Find(&t).Error
	})

	return t, err
}

func (d DatabaseRepository) FindAll() ([]Tenant, error) {
	var tenantSlice = make([]Tenant, 0)
	err := GetConnection(mainTenantKey, func(db *gorm.DB) (ex error) {
		return db.Find(&tenantSlice).Error
	})

	return tenantSlice, err
}

func (d DatabaseRepository) Ping(tenantKey *TenantKey) (SQLResult, error) {
	var result SQLResult
	err := GetConnection(tenantKey, func(db *gorm.DB) error {
		err := db.Raw("SELECT 1").Scan(&result).Error
		return err
	})
	return result, err
}
