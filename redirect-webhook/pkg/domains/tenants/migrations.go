package tenants

import (
	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/helpers/logging"
	"fmt"
	"github.com/gobuffalo/packr"
	"github.com/jinzhu/gorm"
	migrate "github.com/rubenv/sql-migrate"
	"go.uber.org/zap"
)

var runIssuerMigration = true

type DatabaseMigrationManager interface {
	Run(key TenantKey, enabled bool, up bool, isIssuerMigration bool) error
}

type MigrationsManager struct{}

func NewMigrationsManager() DatabaseMigrationManager {
	return MigrationsManager{}
}

func (MigrationsManager) Run(key TenantKey, enabled bool, up bool, isIssuerMigration bool) error {
	dialect := global.Viper.GetString("database.dialect")
	migrationType := "[MIGRATIONS-PRINCIPAL]"
	if isIssuerMigration {
		migrationType = "[MIGRATIONS-ISSUER]"
		if !runIssuerMigration {
			return nil
		}
	}

	if enabled {

		err := GetConnection(&key, func(db *gorm.DB) (ex error) {

			var migrationsDir migrate.PackrMigrationSource
			migrationsDir = migrate.PackrMigrationSource{
				Box: packr.NewBox("../../db/migrations"),
			}
			if isIssuerMigration {
				migrationsDir = migrate.PackrMigrationSource{
					Box: packr.NewBox("../../db/migrations_issuers"),
				}
			}

			direction := migrate.Up
			if !up {
				direction = migrate.Down
			}

			_, err := migrate.Exec(db.DB(), dialect, migrationsDir, direction)

			if err != nil {
				logging.Logger.Error(fmt.Sprintf("%s Error executing migration files", migrationType), zap.Error(err))
				return err
			}
			return nil

		})
		if err != nil {
			return err
		}
		if isIssuerMigration {
			runIssuerMigration = false
		}

	} else {

		logging.Logger.Info(fmt.Sprintf("%s Disabled", migrationType))
	}
	return nil
}
