package tenants

import (
	. "conductor.com.br/pier/hooks/pkg/error"
	"github.com/jinzhu/gorm"
)

type Error struct {
	OrmHandler
}

func (e Error) HandleTenantSetupForIssuer(err error) {
	if err == nil {
		HttpError{
			Status:       400,
			Description:  "tenant is already set up for this issuer",
			DisplayStack: true,
		}.WithChild(err).Throw()
	} else if !gorm.IsRecordNotFoundError(err) {
		e.HandleDatabaseError(err)
	}
}

func (e Error) HandleUnavailableConnection(err error) {
	if err != nil {
		HttpError{
			Status:       400,
			Description:  "database connection is not available",
			DisplayStack: true,
		}.WithChild(err).Throw()
	}
}

func (e Error) HandleUpdateError(err error) {
	if err != nil {
		HttpError{
			Status:       500,
			Description:  "could not update tenant data",
			DisplayStack: true,
		}.WithChild(err).Throw()
	}
}
