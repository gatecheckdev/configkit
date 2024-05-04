package configkit

import (
	"errors"
	"os"
	"reflect"

	"github.com/spf13/cobra"
)

type Field struct {
	FlagP              any
	DefaultValue       any
	FlagName           string
	MapFieldName       string
	Usage              string
	EnvKey             string
	NoEnv              bool
	StringOverrideFunc func(Field) string
	CobraSetupFunc     func(Field, *cobra.Command)
	EnvToValueFunc     func(string) any
	Metadata           map[string]string
}

// Value returns the value based on order of precedence
//
// Order of precedence:
// 1. Flag Value
// 2. Environment Variable Value (only if stringToValueFunc is defined)
// 3. Explicit Default Value
func (f Field) Value() any {
	// Check for a non-zero FlagP value as the first priority
	if f.FlagP != nil {
		val := reflect.ValueOf(f.FlagP)
		if val.Kind() == reflect.Pointer {
			val = val.Elem()
		}
		if !val.IsZero() {
			return val.Interface()
		}
	}

	// Check for an environment variable
	if !f.NoEnv {
		return f.DefaultValue
	}

	value := os.Getenv(f.EnvKey)
	if value != "" {
		if f.EnvToValueFunc == nil {
			panic("An Environment Variable was used but EnvToValueFunc was not defined.")
		}
		return f.EnvToValueFunc(value)
	}

	return f.DefaultValue
}

// FieldValue asserts the type T to the field's value
//
// This function will panic if the type assertion fails. Use SafeFieldValue if the type
// is not known at run time.
func FieldValue[T any](field Field) T {
	assertedValue, ok := field.Value().(T)
	if !ok {
		panic("invalid type assertion")
	}
	return assertedValue
}

func SafeFieldValue[T any](field Field) (T, error) {
	assertedValue, ok := field.Value().(T)
	if !ok {
		return assertedValue, errors.New("invalid type assertion")
	}
	return assertedValue, nil
}

// configFieldsFrom extracts export ConfigKit Fields
//
// Panic Conditions:
//
//   - If v is not a struct or a pointer to a struct
//   - If there is a key conflict
func ConfigFieldMap(v any) map[string]Field {
	// Check if v is a pointer and get its element if it is
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	// Ensure v is a struct
	if val.Kind() != reflect.Struct {
		panic(errors.New("ConfigKit Field: input is not a struct"))
	}

	results := map[string]Field{}
	structs := []reflect.Value{val}

	for len(structs) > 0 {
		val, structs = structs[len(structs)-1], structs[:len(structs)-1]

		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if !field.CanInterface() {
				continue
			}
			// Dereference field if it is a pointer
			if field.Kind() == reflect.Pointer {
				field = field.Elem()
			}

			if field.Kind() != reflect.Struct {
				continue
			}

			configField, ok := field.Interface().(Field)
			// It's a struct but not a config field struct, skip it
			if !ok {
				structs = append(structs, field)
				continue
			}
			name := configField.MapFieldName
			if name == "" {
				name = val.Type().Field(i).Name
			}
			if _, exists := results[name]; exists {
				panic("map key already exists, use MapFieldName to reassign")
			}
			results[name] = configField
		}
	}

	return results
}
