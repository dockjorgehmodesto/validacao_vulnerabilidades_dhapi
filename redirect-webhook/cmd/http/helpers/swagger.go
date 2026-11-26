package helpers

// import (
// 	"conductor.com.br/pier/hooks/pkg/conf/global"
// 	json2 "encoding/json"
// 	"fmt"
// 	"github.com/getkin/kin-openapi/openapi2"
// 	"github.com/getkin/kin-openapi/openapi2conv"
// 	"github.com/getkin/kin-openapi/openapi3"
// 	"github.com/gobuffalo/packr"
// 	"strings"
// )

// type OpenApiSpec struct {
// 	openapi3.Swagger
// 	SpecVersion string `json:"swagger"`
// }

// func FetchSwagger(cleanInternals bool) []byte {
// 	box := packr.NewBox("../docs")
// 	swaggerJsonFile, err := box.Find("swagger.json")
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	swagger, err := readSwaggerJson(swaggerJsonFile)
// 	if err != nil {
// 		fmt.Println(err)
// 		return []byte("")
// 	}

// 	if cleanInternals {
// 		// Clean operations
// 		cleanInternalOperations(&swagger)

// 		// Clean paths
// 		cleanEmptyPaths(&swagger)
// 	}

// 	swagger.Info.Title = global.Viper.GetString("app.name")
// 	swagger.Info.Version = global.VersionInfo.Current
// 	swagger.Info.Description = global.Viper.GetString("app.description")
// 	swagger.Schemes = []string{global.Viper.GetString("swagger.scheme")}
// 	swagger.Host = fmt.Sprintf("%s:%s", global.Viper.GetString("swagger.ip"), global.Viper.GetString("swagger.port"))

// 	json, err := convertToSwaggerV3(swagger)
// 	if err != nil {
// 		fmt.Println(err)
// 		return []byte("")
// 	}
// 	return json
// }

// func convertToSwaggerV3(swagger openapi2.Swagger) ([]byte, error) {
// 	v3Swagger, err2 := openapi2conv.ToV3Swagger(&swagger)
// 	if err2 != nil {
// 		return []byte(""), err2
// 	}

// 	openapi3spec := &OpenApiSpec{
// 		*v3Swagger,
// 		"3.0.0",
// 	}
// 	openapi3spec.OpenAPI = "3.0.0"
// 	marshal, err2 := json2.Marshal(openapi3spec)
// 	if err2 != nil {
// 		return nil, err2
// 	}

// 	replaced := strings.Replace(string(marshal), "#/definitions/", "#/components/schemas/", -1)
// 	return []byte(replaced), nil
// }

// func readSwaggerJson(input []byte) (openapi2.Swagger, error) {
// 	swagger := openapi2.Swagger{}

// 	err := json2.Unmarshal(input, &swagger)
// 	if err != nil {
// 		fmt.Println(err)
// 		return openapi2.Swagger{}, err
// 	}
// 	return swagger, err
// }

// func cleanEmptyPaths(swagger *openapi2.Swagger) {
// 	for key, value := range swagger.Paths {
// 		shouldClean := true
// 		for _, operation := range value.Operations() {
// 			if operation != nil && shouldClean == true {
// 				shouldClean = false
// 			}
// 		}

// 		if shouldClean {
// 			delete(swagger.Paths, key)
// 		}
// 	}
// }

// func cleanInternalOperations(swagger *openapi2.Swagger) {
// 	for _, value := range swagger.Paths {
// 		for method, operation := range value.Operations() {
// 			if isInternalPath(operation) {
// 				value.SetOperation(method, nil)
// 			}
// 		}
// 	}
// }

// func isInternalPath(operation *openapi2.Operation) bool {
// 	contains := false
// 	for _, tag := range operation.Tags {
// 		if contains == false {
// 			contains = tag == "Internals"
// 		}
// 	}
// 	return contains
// }
