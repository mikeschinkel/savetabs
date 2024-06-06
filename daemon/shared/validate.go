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
