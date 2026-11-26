package helpers

import (
	"conductor.com.br/pier/hooks/pkg/conf/global"
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	error2 "conductor.com.br/pier/hooks/pkg/error"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"conductor.com.br/pier/hooks/pkg/helpers/constants"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/validation"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
	"strings"
)

func GetTenant(r *http.Request) (tenantKey *tenants.TenantKey) {
	tenantInterface := r.Context().Value(constants.ConstTenantHeaderContext)

	if tenantInterface == nil {
		// TODO Handle error
		return nil
	}

	if value, ok := tenantInterface.(string); ok {
		tenantKey := tenants.CreateTenantKeyFromIdentifier(value)
		tenantKey.IssuerId = GetIssuerId(r)
		return tenantKey
	} else {
		// TODO Handle error
		return nil
	}
}

func GetPagination(r *http.Request) helpers.Pagination {
	page := coalesceQueryWithDefault(r.URL.Query().Get("page"), global.Viper.GetString("pagination.pageDefault"))
	limit := coalesceQueryWithDefault(r.URL.Query().Get("limit"), global.Viper.GetString("pagination.limitDefault"))

	tmpString := r.URL.Query().Get("sort")
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

	return helpers.Pagination{

		Page:  page,
		Limit: limit,
		Sort:  sortValues,
	}
}

func GetIssuerId(r *http.Request) int {
	issuerIdInterface := r.Context().Value(constants.ConstHeaderIssuerId)

	if issuerIdInterface == nil {
		// TODO Handle error
		return 0
	}

	if value, ok := issuerIdInterface.(int); ok {
		return value
	} else {
		// TODO Handle error
		return 0
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

// TODO Handle errors
func BindQuery(r *http.Request, destination interface{}) {
	if bytes, e := json.Marshal(r.URL.Query()); e != nil {
		HandlerError{}.JsonDecode().WithChild(e).Throw()
	} else {
		newQuery := normalizeQuery(bytes)
		if e := json.Unmarshal([]byte(newQuery), destination); e != nil {
			HandlerError{}.JsonDecode().WithChild(e).Throw()
		}
	}
}

func normalizeQuery(query []byte) string {
	text := fmt.Sprint(string(query))
	text = strings.Replace(text, "[", "", -1)
	text = strings.Replace(text, "]", "", -1)
	return text
}

func BindJson(r *http.Request, destination interface{}) {
	e := render.DecodeJSON(r.Body, destination)
	if e != nil {
		HandlerError{}.JsonDecode().WithChild(e).Throw()
	}
	validateJsonFields(destination)
}

func validateJsonFields(input interface{}) {
	validator := validation.Validation{}
	if isValid, err := validator.Valid(input); err != nil {
		HandlerError{}.JsonValidation().WithChild(err).Throw()
	} else if !isValid || validator.HasErrors() {
		e := HandlerError{}.JsonValidation()
		for _, valError := range validator.Errors {
			e = e.WithChild(error2.UnhandledError{Description: fmt.Sprintf("%s: %s", valError.Key, valError.Message)})
		}
		e.Throw()
	}
}
