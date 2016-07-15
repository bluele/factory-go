package factory

import (
	"reflect"
)

func getAttrName(sf reflect.StructField, tagName string) string {
	tag := sf.Tag.Get(tagName)
	if tag != "" {
		return tag
	}
	return sf.Name
}
