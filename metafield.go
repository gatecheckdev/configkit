package configkit

import (
	"errors"
	"os"
	"reflect"

	"github.com/spf13/cobra"
)

type MetaField struct {
	FieldName      string
	EnvKey         string
	DefaultValue   any
	FlagValueP     any
	Metadata       map[string]string
	CobraSetupFunc func(f MetaField, cmd *cobra.Command)
	EnvToValueFunc func(s string) any
	ValueRules     []func(m MetaField) any
}

func (m MetaField) defaultRules() []func(m MetaField) any {
	flagRule := func(m MetaField) any {
		if m.FlagValueP == nil {
			return nil
		}
		val := reflect.ValueOf(m.FlagValueP).Elem()
		if val.IsZero() {
			return nil
		}
		return val.Interface()
	}

	envRule := func(m MetaField) any {
		if m.EnvKey == "" {
			return nil
		}
		value := os.Getenv(m.EnvKey)
		if value == "" {
			return nil
		}
		envToValueFunc := m.EnvToValueFunc
		if envToValueFunc == nil {
			envToValueFunc = func(s string) any {
				return value
			}
		}
		return envToValueFunc(value)
	}

	defaultValueRule := func(m MetaField) any {
		return m.DefaultValue
	}

	return []func(m MetaField) any{
		flagRule,
		envRule,
		defaultValueRule,
	}
}

func (m MetaField) Value() any {
	rules := m.ValueRules
	var result any
	if len(rules) == 0 {
		rules = m.defaultRules()
	}
	for _, ruleFunc := range rules {
		result = ruleFunc(m)
		if result != nil {
			return result
		}
	}

	return nil
}

func AllMetaFields(v any) []MetaField {
	// Check if v is a pointer and get its element if it is
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	// Ensure v is a struct
	if val.Kind() != reflect.Struct {
		panic(errors.New("ConfigKit AllMetaFields: input is not a struct"))
	}

	results := []MetaField{}
	structs := []reflect.Value{val}

	for len(structs) > 0 {
		val, structs = structs[len(structs)-1], structs[:len(structs)-1]

		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			// Dereference field if it is a pointer
			if field.Kind() == reflect.Pointer {
				field = field.Elem()
			}

			if field.Kind() != reflect.Struct || !field.CanInterface() {
				continue
			}

			// It's a struct but not a config field struct, skip it
			switch configField := field.Interface().(type) {
			case MetaField:
				results = append(results, configField)
			default:
				structs = append(structs, field)

			}
		}
	}

	return results
}
