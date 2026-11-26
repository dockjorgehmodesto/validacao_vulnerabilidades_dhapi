package helpers

import "gopkg.in/jeevatkm/go-model.v1"

//ParseStructToMap
// Parses a structure to a map by using the go-model library
// To customize how the output is made, make sure to use the library's tags.
//
// Example:
//	type StructureThatNeedsToConvertToAMap struct {
//		AbcProp         string `model:"map_key_abc_prop,omitempty"`
//		DefProp         string `model:"map_key_def_prop,omitempty"`
//		NumericProperty int    `model:"intProperty,omitEmpty"`
//		AnotherNumber   int    `model:",omitEmpty"`
//	}
//	newStruct := StructureThatNeedsToConvertToAMap{"", "def", 0, 15}
//	result := ParseStructToMap(newStruct)
//
// The result map will be the following
//	{
//		"map_key_def_prop": "def",
//		"AnotherNumber": 15
//	}
func ParseStructToMap(structure interface{}) map[string]interface{} {
	result, err := model.Map(structure)
	if err != nil {
		return nil
	}
	return result
}
