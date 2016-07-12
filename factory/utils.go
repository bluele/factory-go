package factory

import (
	"reflect"
	"strings"
)

func getAttrName(sf reflect.StructField, tagName string) string {
	tag := sf.Tag.Get(tagName)
	if tag != "" {
		return tag
	}
	return sf.Name
}

func isExported(sf reflect.StructField) bool {
	if sf.Name[0] != strings.ToUpper(string(sf.Name[0]))[0] {
		return false
	}
	return true
}
