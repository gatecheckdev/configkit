package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gatecheckdev/configkit"
	"github.com/spf13/cobra"
)

var Config = struct {
	FieldOne   configkit.Field[string]
	FieldTwo   configkit.Field[int]
	FieldThree configkit.Field[bool]
}{
	FieldOne: configkit.Field[string]{
		FlagP:        new(string),
		DefaultValue: "default value",
		FlagName:     "field-one",
		FlagInit:     "default value",
		Usage:        "the first value as a string",
		EnvKey:       "APP_FIELD_ONE",
		EnvToValueFunc: func(s string) string {
			return s
		},
		CobraSetupFunc: func(f *configkit.Field[string], cmd *cobra.Command) {
			cmd.Flags().String(f.FlagName, f.FlagInit, f.Usage)
		},
	},
	FieldTwo: configkit.Field[int]{
		FlagP:        new(int),
		DefaultValue: 1,
		FlagName:     "field-two",
		FlagInit:     1,
		Usage:        "the second value as a int",
		EnvKey:       "APP_FIELD_TWO",
		EnvToValueFunc: func(s string) int {
			v, _ := strconv.Atoi(s)
			return v
		},
		CobraSetupFunc: func(f *configkit.Field[int], cmd *cobra.Command) {
			cmd.Flags().Int(f.FlagName, f.FlagInit, f.Usage)
		},
	},
	FieldThree: configkit.Field[bool]{
		FlagP:        new(bool),
		DefaultValue: false,
		FlagName:     "field-three",
		FlagInit:     false,
		Usage:        "the first value as a string",
		EnvKey:       "APP_FIELD_THREE",
		EnvToValueFunc: func(s string) bool {
			b, _ := strconv.ParseBool(s)
			return b
		},
	},
}

func main() {
	fmt.Printf(
		"Field 1: %s\nField 2: %d\nField 3: %v\n\n",
		Config.FieldOne.Value(),
		Config.FieldTwo.Value(),
		Config.FieldThree.Value(),
	)
	os.Setenv("APP_FIELD_TWO", "10")
	fmt.Printf(
		"Field 1: %s\nField 2: %d\nField 3: %v\n",
		Config.FieldOne.Value(),
		Config.FieldTwo.Value(),
		Config.FieldThree.Value(),
	)
	/*
		go run ./example/one

		Field 1: default value
		Field 2: 1
		Field 3: false

		Field 1: default value
		Field 2: 10
		Field 3: false
	*/
}
