package helpers

import (
	"github.com/gofrs/uuid"
	"github.com/jinzhu/gorm"
)

type GormHooks struct {
}

// BeforeCreate will set a UUID rather than numeric ID.
func (o *GormHooks) BeforeCreate(scope *gorm.Scope) error {
	uuid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	return scope.SetColumn("ID", uuid.String())
}
