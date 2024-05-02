package configkit

import (
	"os"

	"github.com/spf13/cobra"
)

type Field[T string | bool | int] struct {
	FlagP          *T
	DefaultValue   T
	FlagName       string
	FlagInit       T
	Usage          string
	EnvKey         string
	CobraSetupFunc func(*Field[T], *cobra.Command)
	EnvToValueFunc func(string) T
}

func (f *Field[T]) SetupCobra(cmd *cobra.Command) {
	f.CobraSetupFunc(f, cmd)
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
		if f.EnvToValueFunc != nil {
			return f.EnvToValueFunc(v)
		}
	}

	// 3. Explicit Default Value
	return f.DefaultValue
}
