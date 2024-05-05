package main

import (
	"fmt"
	"os"

	"github.com/gatecheckdev/configkit"
	"github.com/spf13/cobra"
)

type Config struct {
	Name string
}

var MetaConfig = struct {
	Name configkit.MetaField
}{
	Name: configkit.MetaField{
		FieldName:    "Name",
		EnvKey:       "APP_NAME",
		DefaultValue: "some default",
		CobraSetupFunc: func(f configkit.MetaField, cmd *cobra.Command) {
			cmd.Flags().String("name", "", "a name flag")
		},
	},
}

func main() {
	config := Config{}
	configkit.ApplyValues(&config, MetaConfig)
	fmt.Printf("%+v\n", config)
	os.Setenv("APP_NAME", "some value from env")
	config = Config{}
	configkit.ApplyValues(&config, MetaConfig)
	fmt.Printf("%+v\n", config)
}
