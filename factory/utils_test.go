package factory

import (
	"reflect"
	"testing"
)

func TestIsExportedField(t *testing.T) {
	type TestStruct struct {
		ExportedField   string
		unexportedField string
	}

	rt := reflect.TypeOf(TestStruct{})
	if field, _ := rt.FieldByName("ExportedField"); !isExported(field) {
		t.Error("Failed to identify exported field")
	}
	if field, _ := rt.FieldByName("unexportedField"); isExported(field) {
		t.Error("Failed to identify non exported field")
	}
}
