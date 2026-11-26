package helpers

import (
	"conductor.com.br/pier/hooks/pkg/conf/global"
	utils2 "conductor.com.br/pier/hooks/pkg/helpers/primitives"
	"github.com/astaxie/beego/context"
	"math"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
)

// Param 分页参数
type Param struct {
	DB      *gorm.DB
	Page    int
	Limit   int
	OrderBy []string
	ShowSQL bool
}

// Paginator 分页返回
type Paginator struct {
	TotalRecord int         `json:"totalRecord"`
	TotalPage   int         `json:"totalPage"`
	Records     interface{} `json:"records"`
	Offset      int         `json:"offset"`
	Limit       int         `json:"limit"`
	Page        int         `json:"page"`
	PrevPage    int         `json:"prevPage"`
	NextPage    int         `json:"nextPage"`
}

type Pagination struct {
	Page  int
	Limit int
	Sort  []string
}

// Paging 分页
func Paging(p *Param, result interface{}) *Paginator {
	db := p.DB

	if p.ShowSQL {
		db = db.Debug()
	}
	if p.Page < 1 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}
	if len(p.OrderBy) > 0 {
		for _, o := range p.OrderBy {
			db = db.Order(toSnakeCase(o))
		}
	}

	done := make(chan bool, 1)
	var paginator Paginator
	var count int
	var offset int

	go countRecords(db, result, done, &count)

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.Limit
	}

	db.Limit(p.Limit).Offset(offset).Find(result)
	<-done

	paginator.TotalRecord = count
	paginator.Records = result
	paginator.Page = p.Page

	paginator.Offset = offset
	paginator.Limit = p.Limit
	paginator.TotalPage = int(math.Ceil(float64(count) / float64(p.Limit)))

	if p.Page > 1 {
		paginator.PrevPage = p.Page - 1
	} else {
		paginator.PrevPage = p.Page
	}

	if p.Page == paginator.TotalPage {
		paginator.NextPage = p.Page
	} else {
		paginator.NextPage = p.Page + 1
	}
	return &paginator
}

func countRecords(db *gorm.DB, anyType interface{}, done chan bool, count *int) {
	db.Model(anyType).Count(count)
	done <- true
}

func toSnakeCase(txt string) string {
	return utils2.ToSnakeCase(txt)
}

func GetPagination(ctx *context.Context) Pagination {

	page := coalesceQueryWithDefault(ctx.Input.Query("page"), global.Viper.GetString("pagination.pageDefault"))
	limit := coalesceQueryWithDefault(ctx.Input.Query("limit"), global.Viper.GetString("pagination.limitDefault"))

	tmpString := ctx.Input.Query("sort")
	var sortValues []string

	if tmpString == "" {
		sortValues = []string{global.Viper.GetString("pagination.sortDefault")}
	} else {

		if strings.Contains(tmpString, ",") {
			sortValues = strings.Split(tmpString, ",")
		} else {
			sortValues = []string{tmpString}
		}
	}

	return Pagination{

		Page:  page,
		Limit: limit,
		Sort:  sortValues,
	}

}

func coalesceQueryWithDefault(query string, defaultValue string) int {
	var value = query
	if query == "" {
		value = defaultValue
	}
	result, _ := strconv.Atoi(value)
	return result
}

func (p *Paginator) HasContent() bool {
	return p.Page <= p.TotalPage
}
