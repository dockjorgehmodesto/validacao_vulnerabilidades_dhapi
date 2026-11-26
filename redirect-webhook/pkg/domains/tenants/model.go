package tenants

import (
	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.uber.org/zap"
	gorm3 "gopkg.in/DataDog/dd-trace-go.v1/contrib/jinzhu/gorm"
	"time"
)

const TENANT_PRINCIPAL_ID = "TENANT_PRINCIPAL"

const (
	TenantTypeReadWrite = "READ_AND_WRITE"
	TenantTypeRead      = "READ"
)

type TenantKey struct {
	Key      string `json:"key"`
	IssuerId int    `json:"issuerId"`
}

func (t TenantKey) ToString() string {
	bytes, _ := json.Marshal(t)
	return string(bytes)
}

type TenantContract interface {
	GetTenantKey() TenantKey
	FormatUriConnection() string
	OpenConnection() (*gorm.DB, error)
	PingConnection() error
	GetIssuerId() uint
	FormatTenantKey() string
}

type Tenant struct {
	helpers.GormHooks

	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt" sql:"index"`
	IssuerID  uint       `json:"issuer_id" gorm:"column:issuer_id"`
	Dialect   string     `json:"dialect"`
	Uri       string     `json:"uri"`
	User      string     `json:"user"`
	Password  string     `json:"password"`
	Type      string     `json:"type"`
}

func (t Tenant) GetTenantKey() TenantKey {

	return TenantKey{
		Key:      t.FormatTenantKey(),
		IssuerId: int(t.IssuerID),
	}
}

func (t Tenant) FormatTenantKey() string {
	return fmt.Sprintf("%d_%s", t.IssuerID, t.Type)
}

func (t Tenant) FormatUriConnection() string {
	return fmt.Sprintf("%s user=%s password=%s", t.Uri, t.User, t.Password)
}

func (t Tenant) OpenConnection() (*gorm.DB, error) {
	dbConnection, err := gorm.Open(t.Dialect, t.FormatUriConnection())
	if err != nil {
		logging.Logger.Error("[TENANT]", zap.Error(err))
		return nil, err
	}
	dbConnection = gorm3.WithCallbacks(dbConnection)

	t.setupConnection(dbConnection)
	return dbConnection, nil
}

func (t Tenant) setupConnection(db *gorm.DB) {

	debug := global.Viper.GetString("log.level")
	if debug == "debug" {
		db.LogMode(global.ShowSQL)
	}

	db.DB().SetMaxIdleConns(global.Viper.GetInt("database.migrationPrincipal.minimumIdle"))
	db.DB().SetMaxOpenConns(global.Viper.GetInt("database.migrationPrincipal.maximumPoolSize"))
	db.DB().SetConnMaxLifetime(time.Minute * 10)
}

func (t Tenant) PingConnection() error {
	db, e := t.OpenConnection()
	if e != nil {
		return e
	}
	defer db.Close()

	var result SQLResult
	return db.Raw("SELECT 1").Scan(&result).Error
}

func (t Tenant) GetIssuerId() uint {
	return t.IssuerID
}

func CreateTenantKeyString(t *Tenant) string {

	if t.ID == "" {
		return TENANT_PRINCIPAL_ID
	} else {
		return t.FormatTenantKey()
	}

}

func NewMainTenantKey() *TenantKey {

	return &TenantKey{Key: CreateTenantKeyString(&Tenant{IssuerID: 0, Type: "READ"})}

}

func CreateTenantKeyFromIdentifier(id string) *TenantKey {

	return &TenantKey{
		Key: id,
	}

}
