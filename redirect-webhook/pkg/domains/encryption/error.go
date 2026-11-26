package encryption

import (
	. "conductor.com.br/pier/hooks/pkg/error"
	"github.com/jinzhu/gorm"
)

type Error struct {
	OrmHandler
}

func (e Error) HandleInvalidPublicKey(err error) {
	if err != nil {
		HttpError{}.BadRequest("invalid public key").WithChild(err).Throw()
	}
}

func (e Error) HandleExpiredPublicKey(hasExpired bool) {
	if hasExpired {
		HttpError{}.BadRequest("issuer's public key has expired").Throw()
	}
}

func (e Error) HandleAesKey(err error) {
	if err != nil {
		HttpError{}.BadRequest("could not create aes key").WithChild(err).Throw()
	}
}

func (e Error) HandleDisabledVault(enabled bool) {
	if !enabled {
		HttpError{}.BadRequest("key vault is disabled").Throw()
	}
}

func (e Error) HandleAesEncryption(err error) {
	if err != nil {
		HttpError{}.BadRequest("could not encrypt aes key").WithChild(err).Throw()
	}
}

func (e Error) HandleInvalidTargetIssuer(err error) {
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			HttpError{}.BadRequest("target issuer ID is not valid").WithChild(err).Throw()
		} else {
			HttpError{}.HandleDatabaseError(err)
		}
	}
}
