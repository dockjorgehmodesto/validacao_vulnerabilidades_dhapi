package tenants

import (
	"errors"
)

type MigrationsManagerMockSuccess struct{}

func (MigrationsManagerMockSuccess) Run(key TenantKey, enabled bool, up bool, isIssuerMigration bool) error {
	return nil
}

type MigrationsManagerMockError struct{}

func (MigrationsManagerMockError) Run(key TenantKey, enabled bool, up bool, isIssuerMigration bool) error {
	return errors.New("mock")
}
