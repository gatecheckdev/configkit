package main

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/gatecheckdev/configkit"
	"github.com/spf13/cobra"
)

var Config = struct {
	FieldOne   configkit.Field
	FieldTwo   configkit.Field
	FieldThree configkit.Field
	SubConfig  struct {
		SomeOtherField string
		FieldFour      configkit.Field
	}
}{
	FieldOne: configkit.Field{
		FlagP:        new(string),
		DefaultValue: "default value",
		FlagName:     "field-one",
		Usage:        "the first value as a string",
		EnvKey:       "APP_FIELD_ONE",
		MapFieldName: "field1",
		EnvToValueFunc: func(s string) any {
			return s
		},
		CobraSetupFunc: func(f configkit.Field, cmd *cobra.Command) {
			defaultValue, _ := f.DefaultValue.(string)
			cmd.Flags().String(f.FlagName, defaultValue, f.Usage)
		}},
	FieldTwo: configkit.Field{
		FlagP:        new(int),
		DefaultValue: 1,
		FlagName:     "field-two",
		Usage:        "the second value as a int",
		EnvKey:       "APP_FIELD_TWO",
		MapFieldName: "field2",
		EnvToValueFunc: func(s string) any {
			v, _ := strconv.Atoi(s)
			return v
		},
		CobraSetupFunc: func(f configkit.Field, cmd *cobra.Command) {
			cmd.Flags().Int(f.FlagName, 0, f.Usage)
		}},

	FieldThree: configkit.Field{
		FlagP:        new(bool),
		DefaultValue: false,
		FlagName:     "field-three",
		Usage:        "the first value as a string",
		EnvKey:       "APP_FIELD_THREE",
		MapFieldName: "field3",
		EnvToValueFunc: func(s string) any {
			b, _ := strconv.ParseBool(s)
			return b
		}},
	SubConfig: struct {
		SomeOtherField string
		FieldFour      configkit.Field
	}{
		SomeOtherField: "some other value",
		FieldFour: configkit.Field{
			FlagP:        new(int),
			DefaultValue: 2,
			FlagName:     "field-four",
			Usage:        "an example of a sub field in a nested config",
			EnvKey:       "APP_FIELD_FOUR",
			MapFieldName: "field4",
			EnvToValueFunc: func(s string) any {
				v, _ := strconv.Atoi(s)
				return v
			},
		},
	},
}

func main() {
	configMap := configkit.ConfigFieldMap(Config)
	keys := make([]string, 0, len(configMap))
	for k := range configMap {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	for _, key := range keys {
		fmt.Printf("key: %-8s value: %v\n", key, configMap[key].Value())
	}

	// go run ./examples/one
	//
	// key: field1   value: default value
	// key: field2   value: 1
	// key: field3   value: false
	// key: field4   value: 2
}
