package tenants

import (
	"errors"
	"github.com/jinzhu/gorm"
)

const MockTenantKey = "MOCKTENANTKEY"

type TenantMockError struct{}

func (t TenantMockError) GetTenantKey() TenantKey {
	return TenantKey{}
}
func (t TenantMockError) FormatUriConnection() string {
	return ""
}
func (t TenantMockError) OpenConnection() (*gorm.DB, error) {
	return nil, errors.New("database error mock")
}
func (t TenantMockError) PingConnection() error {
	return errors.New("ping error mock")
}

func (t TenantMockError) GetIssuerId() uint {
	return 0
}

func (t TenantMockError) FormatTenantKey() string {
	return ""
}

type TenantMockSuccess struct {
	Key string
}

func (t TenantMockSuccess) GetTenantKey() TenantKey {
	return TenantKey{Key: t.Key}
}
func (t TenantMockSuccess) FormatUriConnection() string {
	return ""
}
func (t TenantMockSuccess) OpenConnection() (*gorm.DB, error) {
	return &gorm.DB{}, nil
}
func (t TenantMockSuccess) PingConnection() error {
	return nil
}

func (t TenantMockSuccess) GetIssuerId() uint {
	return 0
}

func (t TenantMockSuccess) FormatTenantKey() string {
	return ""
}
