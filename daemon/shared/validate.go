package shared

// See https://github.com/go-playground/validator as potential future
// alternative, but only if needed.

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// ValidateStruct validates the struct fields based on a small set of valid `valdidate` tags
func ValidateStruct(s any) (err error) {
	return validateStruct(s, "")
}
func validateStruct(s any, parent string) (err error) {
	var me *MultiErr

	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	if v.Kind() != reflect.Struct {
		err = ErrParamIsNotKindOfStruct
		goto end
	}

	me = NewMultiErr()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)
		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}
		if err := validateField(field, fieldValue, tag, parent); err != nil {
			me.Add(err)
		}
	}
end:
	return nil
}

// Helper function to validate a single field
func validateField(field reflect.StructField, fieldValue reflect.Value, tag string, parent string) error {
	tags := strings.Split(tag, ",")
	me := NewMultiErr()
	for _, t := range tags {
		switch t {
		case "required":
			if fieldValue.IsZero() {
				var err error
				if parent == "" {
					err = fmt.Errorf("field=%s", field.Name)
				} else {
					err = fmt.Errorf("field=%s.%s", parent, field.Name)
				}
				me.Add(errors.Join(ErrNonZeroFieldValueRequired, err))
			}
		}
		if field.Type.Kind() == reflect.Struct {
			continue
		}
		err := validateStruct(fieldValue.Interface(), fmt.Sprintf("%s.%s", parent, field.Name))
		if err != nil {
			me.Add(err)
		}
	}
	return me.Err()
}

// Helper function to check if a value is zero value
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	case reflect.Struct:
		zero := reflect.New(v.Type())
		elem := zero.Elem()
		zeroStruct := elem.Interface()
		v.IsZero()
		return v.Interface() == zeroStruct
	default:
		return false
	}
}
