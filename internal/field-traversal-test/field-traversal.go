package main

import (
	"errors"
	"fmt"
	"reflect"
)

func main() {
	var Config = TestConfig{
		Name: "Base Config",
		FieldOne: Field[string]{
			Name:  "field-one-1-base-config",
			Value: "some value 1",
		},
		FieldTwo: Field[bool]{
			Name:  "field-two-2-base-config",
			Value: true,
		},
		FieldThree: Field[int]{
			Name:  "field-three-3-base-config",
			Value: 32,
		},
		ConfigA: TestConfigA{
			Name: "Config A",
			FieldFour: Field[string]{
				Name:  "field-four-4-config-a",
				Value: "some other value for field 4",
			},
			FieldFive: Field[bool]{
				Name:  "field-five-5-config-a",
				Value: false,
			},
			FieldSix: Field[int]{
				Name:  "field-six-6-config-a",
				Value: 32089,
			},
		},
		ConfigB: TestConfigB{
			Name: "Config B",
			FieldSeven: Field[string]{
				Name:  "field-seven-7-config-b",
				Value: "Some 7th value",
			},
			FieldEight: Field[bool]{
				Name:  "field-eight-8-config-b",
				Value: true,
			},
			FieldNine: Field[int]{
				Name:  "field-nine-9-config-b",
				Value: 1,
			},
		},
	}

	stack := getConfigFieldsStack(Config)
	recrs := getConfigFieldsRecursive(Config)

	for i, field := range stack {
		name := field.FieldByName("Name").Interface()
		value := field.FieldByName("Value").Interface()
		fmt.Printf("%2d Field Name: %-30v Value: %v\n", i, name, value)
	}
	fmt.Println("-----------------------------")
	for i, field := range recrs {
		name := field.FieldByName("Name").Interface()
		value := field.FieldByName("Value").Interface()
		fmt.Printf("%2d Field Name: %-30v Value: %v\n", i, name, value)
	}
}

type TestConfig struct {
	Name       string
	FieldOne   Field[string]
	FieldTwo   Field[bool]
	FieldThree Field[int]
	ConfigA    TestConfigA
	ConfigB    TestConfigB
}

type TestConfigA struct {
	Name      string
	FieldFour Field[string]
	FieldFive Field[bool]
	FieldSix  Field[int]
}

type TestConfigB struct {
	Name       string
	FieldSeven Field[string]
	FieldEight Field[bool]
	FieldNine  Field[int]
}

type Field[T string | int | bool] struct {
	Name  string
	Value T
}

func getConfigFieldsStack(v any) []reflect.Value {
	// Check if v is a pointer and get its element if it is
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	// Ensure v is a struct
	if val.Kind() != reflect.Struct {
		panic(errors.New("ConfigKit Field: input is not a struct"))
	}

	configFields := []reflect.Value{}
	nonConfigFields := []reflect.Value{val}
	for {
		if len(nonConfigFields) == 0 {
			break
		}
		val, nonConfigFields = nonConfigFields[len(nonConfigFields)-1], nonConfigFields[:len(nonConfigFields)-1]

		for i := 0; i < val.NumField(); i++ {
			field := val.Field(i)
			if !field.CanInterface() {
				continue
			}
			switch field.Interface().(type) {
			case Field[string], Field[bool], Field[int]:
				configFields = append(configFields, field)
			default:
				// reasign the field pointer to the field itself
				if field.Kind() == reflect.Pointer {
					field = field.Elem()
				}
				// Skip a non struct sub field
				if field.Kind() != reflect.Struct {
					continue
				}
				nonConfigFields = append(nonConfigFields, field)
			}
		}
	}
	return configFields
}

func getConfigFieldsRecursive(v any) []reflect.Value {
	// Check if v is a pointer and get its element if it is
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}

	// Ensure v is a struct
	if val.Kind() != reflect.Struct {
		panic(errors.New("ConfigKit Field: input is not a struct"))
	}
	configFields := []reflect.Value{}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.CanInterface() {
			continue
		}
		switch field.Interface().(type) {
		case Field[string], Field[bool], Field[int]:
			configFields = append(configFields, field)
		default:
			// reasign the field pointer to the field itself
			if field.Kind() == reflect.Pointer {
				field = field.Elem()
			}
			// Skip a non struct sub field
			if field.Kind() != reflect.Struct {
				continue
			}
			configFields = append(configFields, getConfigFieldsRecursive(field.Interface())...)
		}
	}
	return configFields
}
