package sirius

import "testing"

type ExtensionConfigAccessorTest struct {
	t   *testing.T
	f   func(string, ExtensionConfig) interface{}
	exp valueMatches
	def interface{}
}

type valueMatches map[string]interface{}

func testExtensionConfigAccessor(test ExtensionConfigAccessorTest) {
	cfg := testingExtensionConfig()

	for field := range cfg {
		a := test.f(field, cfg)

		e, match := test.exp[field]

		if match {
			if a != e {
				test.t.Fatalf("Expected (%s) to resolve into %v, got %v", field, e, a)
			}
			return
		}

		if a != test.def {
			test.t.Fatalf("Expected default value '%v' for (%s), got %v", test.def, field, a)
		}
	}
}

func testingExtensionConfig() ExtensionConfig {
	return ExtensionConfig{
		"int_0":        0,
		"int_1":        1,
		"float_0.0":    0.0,
		"float_1.1":    1.1,
		"bool_true":    true,
		"bool_false":   false,
		"string_hello": "hello",
		"string_empty": "",
		"list":         []string{"Hit", "Me", "Up"},
	}

}

func TestExtensionConfig_Boolean(t *testing.T) {
	testExtensionConfigAccessor(ExtensionConfigAccessorTest{
		t: t,
		f: func(k string, cfg ExtensionConfig) interface{} {
			return cfg.Boolean(k)
		},
		exp: valueMatches{
			"int_0":      false,
			"int_1":      true,
			"bool_true":  true,
			"bool_false": false,
		},
		def: false,
	})
}

func TestExtensionConfig_Integer(t *testing.T) {
	testExtensionConfigAccessor(ExtensionConfigAccessorTest{
		t: t,
		f: func(k string, cfg ExtensionConfig) interface{} {
			return cfg.Integer(k, 999)
		},
		exp: valueMatches{
			"int_0": 0,
			"int_1": 1,
		},
		def: 999,
	})
}

func TestExtensionConfig_Float(t *testing.T) {
	testExtensionConfigAccessor(ExtensionConfigAccessorTest{
		t: t,
		f: func(k string, cfg ExtensionConfig) interface{} {
			return cfg.Float(k, 999.99)
		},
		exp: valueMatches{
			"float_0.0": 0.0,
			"float_1.1": 1.1,
		},
		def: 999.99,
	})
}

func TestExtensionConfig_String(t *testing.T) {
	testExtensionConfigAccessor(ExtensionConfigAccessorTest{
		t: t,
		f: func(k string, cfg ExtensionConfig) interface{} {
			return cfg.String(k, "Darth Vader")
		},
		exp: valueMatches{
			"string_hello": "hello",
			"string_empty": "",
		},
		def: "Darth Vader",
	})
}

func TestExtensionConfig_Read(t *testing.T) {
	testExtensionConfigAccessor(ExtensionConfigAccessorTest{
		t: t,
		f: func(k string, cfg ExtensionConfig) interface{} {
			return cfg.Read(k, nil)
		},
		exp: valueMatches{
			"int_0":        0,
			"int_1":        1,
			"float_0.0":    0.0,
			"float_1.1":    1.1,
			"bool_true":    true,
			"bool_false":   false,
			"string_hello": "hello",
			"string_empty": "",
		},
		def: nil,
	})
}

func TestExtensionConfig_List(t *testing.T) {
	cfg := testingExtensionConfig()

	exp := map[string]bool{
		"list": true,
	}

	for field := range cfg {
		a := cfg.List(field, nil)

		e, match := exp[field]

		if match {
			if a[0] != "Hit" || a[1] != "Me" || a[2] != "Up" {
				t.Fatalf("Expected (%s) to resolve into %v, got %v", field, e, a)
			}
			return
		}

		if a != nil {
			t.Fatalf("Expected default value '%v' for (%s), got %v", nil, field, a)
		}
	}
}
