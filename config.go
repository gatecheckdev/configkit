package configkit

import (
	"errors"
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/cobra"
)

// Field represents a configurable option that can be used in various parts of an application,
// such as command line interface tools, configuration management, and environment variable parsing.
// It is a generic type that supports string, bool, and int types.
//
// Fields:
//   - FlagP: Pointer to the variable that stores the value of the command line flag.
//   - DefaultValue: The default value for the field when not explicitly provided by an input source.
//   - FlagName: The name of the command line flag as used in CLI applications, typically with packages like cobra.
//   - Usage: Description of the flag's purpose, displayed in help messages.
//   - EnvKey: The environment variable key associated with the field. If set, the application can look for this
//     environment variable to configure the value.
//   - NoEnv: A boolean indicating whether the environment variable should be ignored. If true, the application
//     does not use the environment variable even if it is set.
//   - StringOverrideFunc: Optional function to provide a custom string representation of the field. This can be
//     used to control how the field is displayed in logs or user interfaces.
//   - CobraSetupFunc: Function to integrate this field into a cobra.Command. This function typically sets up
//     bindings and defaults for the command based on the field's properties.
//   - EnvToValueFunc: Function to convert an environment variable string value to the field's type. This is used
//     when parsing environment variables to ensure type safety and correct parsing logic.
//   - Metadata: A map of string keys to string values for storing additional metadata about the field. This can
//     be used to store attributes like validation rules, group identifiers, or other contextual data.
//
// Example:
//
//	field := Field[int]{
//	    FlagP: &someIntVar,
//	    DefaultValue: 10,
//	    FlagName: "verbosity",
//	    Usage: "Adjust the verbosity level of the application",
//	    EnvKey: "APP_VERBOSITY",
//	    NoEnv: false,
//	    StringOverrideFunc: func(f *Field[int]) string {
//	        return fmt.Sprintf("Current verbosity level: %d", *f.FlagP)
//	    },
//	    CobraSetupFunc: func(f *Field[int], cmd *cobra.Command) {
//	        cmd.Flags().IntVar(f.FlagP, f.FlagName, f.DefaultValue, f.Usage)
//	    },
//	    EnvToValueFunc: func(val string) int {
//	        parsed, _ := strconv.Atoi(val)
//	        return parsed
//	    },
//	    Metadata: map[string]string{"group": "logging"},
//	}
//
// This example demonstrates setting up a Field for an integer type to manage verbosity levels in a CLI application,
// using the cobra library for command line parsing.
type Field[T string | bool | int] struct {
	FlagP              *T
	DefaultValue       T
	FlagName           string
	MapFieldName       string
	Usage              string
	EnvKey             string
	NoEnv              bool
	StringOverrideFunc func(*Field[T]) string
	CobraSetupFunc     func(*Field[T], *cobra.Command)
	EnvToValueFunc     func(string) T
	Metadata           map[string]string
}

// SetupCobra runs the CobraSetupFunc if not set to nil
//
// This can be used to setup flags at runtime
func (f *Field[T]) SetupCobra(cmd *cobra.Command) {
	if f.CobraSetupFunc != nil {
		f.CobraSetupFunc(f, cmd)
	}
}

// Value returns the value based on order of precedence
//
// Order of precedence:
// 1. Flag Value
// 2. Environment Variable Value (only if stringToValueFunc is defined)
// 3. Explicit Default Value
func (f *Field[T]) Value() T {
	// 1. Flag Value
	if f.FlagP != nil && *f.FlagP != *new(T) {
		return *f.FlagP
	}
	// 2. Environment Variable Value
	if v := os.Getenv(f.EnvKey); v != "" {
		if f.EnvToValueFunc != nil || !f.NoEnv {
			return f.EnvToValueFunc(v)
		}
	}

	// 3. Explicit Default Value
	return f.DefaultValue
}

// String returns a string representation of the Field instance. It provides a default
// string format but also supports a customizable representation through the StringOverrideFunc.
//
// If StringOverrideFunc is set, it will call this function, passing the current Field instance,
// and will return the result of this function as the Field's string representation. This allows
// for dynamic and context-specific string representations of Field values, which can be particularly
// useful for logs, user interfaces, or any other outputs where Field data needs to be human-readable.
//
// If StringOverrideFunc is not set, the method defaults to returning a formatted string of the Field's
// value using fmt.Sprintf with the "%+v" verb, which includes both field names and values for structs.
//
// Returns:
//   - string: The string representation of the Field instance, either from the custom function or the default format.
//
// Example:
//
//	verbosityField := Field[int]{
//	    DefaultValue: 3,
//	    StringOverrideFunc: func(f *Field[int]) string {
//	        return fmt.Sprintf("Verbosity Level: %d", *f.FlagP)
//	    },
//	}
//	fmt.Println(verbosityField.String())  // Output: "Verbosity Level: 3"
//
// This example demonstrates the flexibility of the String method to utilize a custom string
// representation if provided. It showcases setting a custom function that adjusts the output
// string based on the field's current value.
func (f *Field[T]) String() string {
	if f.StringOverrideFunc != nil {
		return f.StringOverrideFunc(f)
	}
	return fmt.Sprintf("%+v", f.Value())
}

// ConfigFieldMap traverses a struct and constructs a map where each key is a field name
// and each value is the corresponding field's value. The function is designed to handle structs
// containing fields of types Field[string], Field[bool], or Field[int], which are generic types
// defined to hold specific configurations. It also supports nested structs by recursively processing
// each field that is itself a struct.
//
// Parameters:
//   - v: A struct or a pointer to a struct. The function will panic if 'v' is not a struct or a pointer to a struct.
//
// Returns:
//   - map[string]any: A map where each key represents the name of a field within the struct, and
//     each value is the field's value. If a field is of a generic type (Field[T]),
//     the value is obtained from the Value() method of that field.
//
// Panics:
//   - The function panics if the input 'v' is not a struct or a pointer to a struct, ensuring that the
//     caller must provide a valid struct input.
//   - It also panics if two fields result in the same map key unless explicitly handled by setting
//     MapFieldName within the Field types to unique values.
//
// Usage:
//
//	This function is useful in scenarios where configuration data stored in structs needs to be accessed
//	programmatically in a generic manner, such as when loading configurations from a file or environment
//	variables into a structured format in Go.
//
// Example:
//
//	type AppConfig struct {
//	    DatabaseURL Field[string]  `config:"dbUrl"`
//	    LogLevel    Field[string]  `config:"logLevel"`
//	    DebugMode   Field[bool]    `config:"debug"`
//	}
//
//	config := AppConfig{
//	    DatabaseURL: Field[string]{Value: "localhost"},
//	    LogLevel:    Field[string]{Value: "info"},
//	    DebugMode:   Field[bool]{Value: true},
//	}
//
//	configMap := ConfigFieldMap(config)
//	fmt.Println(configMap["DatabaseURL"]) // Output: localhost
//	fmt.Println(configMap["DebugMode"])   // Output: true
//
// This example shows how ConfigFieldMap can be used to convert a struct into a map for easy access to configuration values.
// Each field of the struct should be tagged or named appropriately to ensure that keys in the resulting map are unique and
// meaningful, avoiding collisions and panics.
func ConfigFieldMap(v any) map[string]any {
	// Check if v is a pointer and get its element if it is
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	// Ensure v is a struct
	if val.Kind() != reflect.Struct {
		panic(errors.New("ConfigKit Field: input is not a struct"))
	}

	results := map[string]any{}

	// interate over all the fields and find configFields
	// the interface is to handle the generic nature of the ConfigField type
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.CanInterface() {
			continue
		}
		fieldInterface := field.Interface()

		name := ""
		var value any
		switch configField := fieldInterface.(type) {
		case Field[string]:
			name = configField.MapFieldName
			value = configField.Value()
		case Field[bool]:
			name = configField.MapFieldName
			value = configField.Value()
		case Field[int]:
			name = configField.MapFieldName
			value = configField.Value()
		default:
			subVal := reflect.ValueOf(fieldInterface)
			if subVal.Kind() == reflect.Pointer {
				subVal = subVal.Elem()
			}

			if subVal.Kind() != reflect.Struct {
				continue
			}
			for key, value := range ConfigFieldMap(subVal.Interface()) {
				if _, ok := results[name]; ok {
					panic(errors.New("ConfigKit Field: map key already exists. use MapFieldName to override the default value"))
				}
				results[key] = value
			}
			continue
		}

		if name == "" {
			name = val.Type().Field(i).Name
		}

		if _, ok := results[name]; ok {
			panic(errors.New("ConfigKit Field: map key already exists. use MapFieldName to override the default value"))
		}

		results[name] = value
	}

	return results
}
