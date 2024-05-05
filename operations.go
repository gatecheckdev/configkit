package configkit

import (
	"errors"
	"reflect"
)

func ApplyValues(config any, metaConfig any) error {
	configVal := reflect.ValueOf(config)
	if reflect.ValueOf(config).Kind() != reflect.Pointer {
		return errors.New("config is not a pointer")
	}
	configVal = configVal.Elem()
	metaFields := AllMetaFields(metaConfig)
	for _, metaField := range metaFields {
		configField := configVal.FieldByName(metaField.FieldName)
		configField.Set(reflect.ValueOf(metaField.Value()))
	}
	return nil
}
