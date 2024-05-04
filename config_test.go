package configkit

import "testing"

type TestConfig struct {
	Name       string
	FieldOne   Field
	FieldTwo   Field
	FieldThree Field
	ConfigA    TestConfigA
	ConfigB    TestConfigB
}

type TestConfigA struct {
	Name      string
	FieldFour Field
	FieldFive Field
	FieldSix  Field
}

type TestConfigB struct {
	Name       string
	FieldSeven Field
	FieldEight Field
	FieldNine  Field
}

func TestGetConfigFields(t *testing.T) {
	var Config = TestConfig{
		Name:       "Base Config",
		FieldOne:   Field{},
		FieldTwo:   Field{},
		FieldThree: Field{},
		ConfigA: TestConfigA{
			Name:      "Config A",
			FieldFour: Field{},
			FieldFive: Field{},
			FieldSix:  Field{},
		},
		ConfigB: TestConfigB{
			Name:       "Config B",
			FieldSeven: Field{},
			FieldEight: Field{},
			FieldNine:  Field{},
		},
	}

	configMap := ConfigFieldMap(Config)

	for key, field := range configMap {
		t.Logf("%-15s %v\n", key, field.Value())
	}
}
