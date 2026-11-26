package error

import (
	"conductor.com.br/pier/hooks/pkg/helpers"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
)

type OrmHandler struct {
	HttpError
	IsSearchError bool
}

func (o OrmHandler) HandleSearch(e error) {
	if e != nil {
		if gorm.IsRecordNotFoundError(e) {
			o.Status = http.StatusNoContent
			o.Description = "no result for your search"
			o.Throw()
		} else {
			o.HandleDatabaseError(e)
		}
	}
}

func (o OrmHandler) HandleSearchWithMessage(e error, message string) {
	if e != nil {
		if gorm.IsRecordNotFoundError(e) {
			o.Status = http.StatusNoContent
			o.Description = fmt.Sprintf("record not found for %s", message)
			o.Throw()
		} else {
			o.HandleDatabaseError(e)
		}
	}
}

func (o OrmHandler) HandlePagination(pagination helpers.Paginator, e error) {
	if e != nil {
		o.HandleDatabaseError(e)
	} else if pagination.TotalRecord <= 0 || !pagination.HasContent() {
		o.Status = http.StatusNoContent
		o.Description = "no result for your search"
		o.Throw()
	}
}
