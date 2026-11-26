package webhooks

import (
	. "conductor.com.br/pier/hooks/pkg/error"
	"github.com/jinzhu/gorm"
)

type Error struct {
	OrmHandler
}

func (e Error) HandleHttpRequest(err error) {
	if err != nil {
		HttpError{Status: 400, Description: "error when calling webHook"}.WithChild(err).Throw()
	}
}

func (e Error) HandleKeyVaultNotFound(err error) {
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			HttpError{Status: 400, Description: "keyVaultId is invalid"}.WithChild(err).Throw()
		} else {
			HttpError{Status: 500, Description: "error fetching keyVault on webhook trigger"}.WithChild(err).Throw()
		}
	}
}

func (e Error) HandlePayloadEncryption(err error) {
	if err != nil {
		HttpError{Status: 500, Description: "error encrypting payload on webhook trigger"}.WithChild(err).Throw()
	}
}

func (e Error) HandleAesEncryption(err error) {
	if err != nil {
		HttpError{Status: 500, Description: "error encrypting aes key"}.WithChild(err).Throw()
	}
}
