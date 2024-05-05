package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gatecheckdev/configkit"
	"github.com/spf13/cobra"
)

type Config struct {
	FieldOne   string
	FieldTwo   int
	FieldThree bool
	FieldFour  int
}

var metaConfig = struct {
	FieldOne   configkit.MetaField
	FieldTwo   configkit.MetaField
	FieldThree configkit.MetaField
	SubConfig  struct {
		SomeOtherField string
		FieldFour      configkit.MetaField
	}
}{
	FieldOne: configkit.MetaField{
		FlagValueP:   new(string),
		DefaultValue: "default value",
		FieldName:    "FieldOne",
		EnvKey:       "APP_FIELD_ONE",
		Metadata: map[string]string{
			"flag_name":  "field-one",
			"flag_usage": "a field one example cobra flag",
		},
		EnvToValueFunc: func(s string) any {
			return s
		},
		CobraSetupFunc: func(f configkit.MetaField, cmd *cobra.Command) {
			defaultValue, _ := f.DefaultValue.(string)
			cmd.Flags().StringVar(f.FlagValueP.(*string), f.Metadata["flag_name"], defaultValue, f.Metadata["flag_usage"])
		}},
	FieldTwo: configkit.MetaField{
		FlagValueP:   new(int),
		DefaultValue: 1,
		FieldName:    "FieldTwo",
		Metadata: map[string]string{
			"flag_name":  "field-two",
			"flag_usage": "the second value as a int",
		},
		EnvKey: "APP_FIELD_TWO",
		EnvToValueFunc: func(s string) any {
			v, _ := strconv.Atoi(s)
			return v
		},
		CobraSetupFunc: func(f configkit.MetaField, cmd *cobra.Command) {
			cmd.Flags().IntVar(f.FlagValueP.(*int), f.Metadata["flag_name"], 0, f.Metadata["flag_usage"])
		}},

	FieldThree: configkit.MetaField{
		FlagValueP:   new(bool),
		DefaultValue: false,
		FieldName:    "FieldThree",
		Metadata: map[string]string{
			"flag_name":  "field-three",
			"flag_usage": "the first value as a string",
		},
		EnvKey: "APP_FIELD_THREE",
		CobraSetupFunc: func(f configkit.MetaField, cmd *cobra.Command) {
			cmd.Flags().BoolVar(f.FlagValueP.(*bool), f.Metadata["flag_name"], false, f.Metadata["flag_usage"])
		},
		EnvToValueFunc: func(s string) any {
			b, _ := strconv.ParseBool(s)
			return b
		}},
	SubConfig: struct {
		SomeOtherField string
		FieldFour      configkit.MetaField
	}{
		SomeOtherField: "some other value",
		FieldFour: configkit.MetaField{
			FlagValueP:   new(int),
			DefaultValue: 2,
			FieldName:    "FieldFour",
			Metadata: map[string]string{
				"flag_name":  "field-four",
				"flag_usage": "an example of a sub field in a nested config",
			},
			EnvKey: "APP_FIELD_FOUR",
			EnvToValueFunc: func(s string) any {
				v, _ := strconv.Atoi(s)
				return v
			},
			CobraSetupFunc: func(f configkit.MetaField, cmd *cobra.Command) {
				cmd.Flags().IntVar(f.FlagValueP.(*int), f.Metadata["flag_name"], f.DefaultValue.(int), f.Metadata["flag_usage"])
			},
		},
	},
}

func main() {
	configFields := configkit.AllMetaFields(metaConfig)

	for i, field := range configFields {
		fmt.Printf("%-2d field name: %-20s value: %v\n", i, field.FieldName, field.Value())
	}

	config := Config{}
	err := configkit.ApplyValues(&config, metaConfig)
	if err != nil {
		panic(err)
	}
	fmt.Println(strings.Repeat("-", 50))
	fmt.Printf("%-2d field name: %-20s value: %v\n", 0, "FieldOne", config.FieldOne)
	fmt.Printf("%-2d field name: %-20s value: %v\n", 1, "FieldTwo", config.FieldTwo)
	fmt.Printf("%-2d field name: %-20s value: %v\n", 2, "FieldThree", config.FieldThree)
	fmt.Printf("%-2d field name: %-20s value: %v\n", 3, "FieldFour", config.FieldFour)

	// go run ./examples/one
	//
	// 0  field name: FieldOne             value: default value
	// 1  field name: FieldTwo             value: 1
	// 2  field name: FieldThree           value: false
	// 3  field name: FieldFour            value: 2
	// --------------------------------------------------
	// 0  field name: FieldOne             value: default value
	// 1  field name: FieldTwo             value: 1
	// 2  field name: FieldThree           value: false
	// 3  field name: FieldFour            value: 2
}
