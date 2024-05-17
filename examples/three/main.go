package main

import (
	"fmt"
	"os"
	"reflect"

	"github.com/gatecheckdev/configkit"
	"github.com/spf13/cobra"
)

type metaConfig struct {
	Name configkit.MetaField
}

func newMetadata(flagName, flagUsage string) map[string]string {
	m := map[string]string{
		"flag_name":  flagName,
		"flag_usage": flagUsage,
	}
	return m
}

func defaultRules() []func(m configkit.MetaField) any {
	flagRule := func(m configkit.MetaField) any {
		if m.FlagValueP == nil {
			return nil
		}
		val := reflect.ValueOf(m.FlagValueP).Elem()
		if val.IsZero() {
			return nil
		}
		return val.Interface()
	}

	envRule := func(m configkit.MetaField) any {
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

	defaultValueRule := func(m configkit.MetaField) any {
		return m.DefaultValue
	}

	return []func(m configkit.MetaField) any{
		flagRule,
		envRule,
		defaultValueRule,
	}
}
func NewDefaultConfig() metaConfig {
	return metaConfig{
		Name: configkit.MetaField{
			FieldName:    "Name",
			EnvKey:       "APP_NAME",
			DefaultValue: "Bro",
			ValueRules:   defaultRules(),
			FlagValueP:   new(string),
			Metadata:     newMetadata("name", "the name of the person to greet"),
			CobraSetupFunc: func(f configkit.MetaField, cmd *cobra.Command) {
				cmd.Flags().StringVar(f.FlagValueP.(*string), f.Metadata["flag_name"], "", f.Metadata["flag_usage"])
			},
		},
	}
}

func NewCommand(config metaConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app",
		Short: "a simple CLI app",
	}

	howdyCmd := &cobra.Command{
		Use:   "howdy",
		Short: "a friendly sub command",
		Run:   runHowdyFunc(config),
	}

	config.Name.SetupCobra(howdyCmd)

	cmd.AddCommand(howdyCmd)

	return cmd
}

func runHowdyFunc(config metaConfig) func(*cobra.Command, []string) {
	return func(c *cobra.Command, s []string) {
		fmt.Println("Howdy! Name:", config.Name.Value().(string))
	}
}

func main() {

	err := NewCommand(NewDefaultConfig()).Execute()
	if err != nil {
		panic(err)
	}
}
