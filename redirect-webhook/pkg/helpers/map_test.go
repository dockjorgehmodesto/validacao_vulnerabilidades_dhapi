package helpers

import (
	"reflect"
	"testing"
)

type CustomStruct struct {
	A string
	B string
	C int
}

type CustomAnnotatedStruct struct {
	A string `model:",omitempty"`
	B string `model:"renamedb"`
	C int    `model:"renamedc"`
}

func TestParseStructToMap(t *testing.T) {

	type args struct {
		structure interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{"empty result", args{nil}, nil},
		{"valid result with no tags", args{CustomStruct{"a", "b", 1}}, map[string]interface{}{"A": "a", "B": "b", "C": 1}},
		{"valid result with tags", args{CustomAnnotatedStruct{"", "b", 1}}, map[string]interface{}{"renamedb": "b", "renamedc": 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseStructToMap(tt.args.structure); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseStructToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
