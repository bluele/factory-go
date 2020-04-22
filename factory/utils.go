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

func setValueWithAttrPath(inst *reflect.Value, tp reflect.Type, attr string, v interface{}) bool {
	attrs := strings.Split(attr, ".")
	if len(attrs) <= 1 {
		return false
	}
	current := inst
	currentTp := tp
	isSet := true
	for _, attr := range attrs {
		rt, rv := indirectPtrValue(currentTp, *current)
		currentTp = rt
		current = &rv

		if currentTp.Kind() != reflect.Struct {
			isSet = false
			break
		}
		ftp, ok := currentTp.FieldByName(attr)
		if !ok {
			isSet = false
			break
		}

		field := current.FieldByName(attr)
		if field == emptyValue {
			isSet = false
			break
		}

		if ftp.Type.Kind() == reflect.Ptr && field.IsNil() {
			field.Set(reflect.New(ftp.Type.Elem()))
		}
		current = &field
		currentTp = ftp.Type
	}
	if isSet {
		current.Set(reflect.ValueOf(v))
	}
	return isSet
}

func indirectPtrValue(rt reflect.Type, rv reflect.Value) (reflect.Type, reflect.Value) {
	for rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}
	return rt, rv
}
