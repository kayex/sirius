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
	cfg["list"] = []string{"Hit", "Me", "Up"}

	match := map[string][]string{
		"list": {"Hit", "Me", "Up"},
	}

	for field := range cfg {
		a := cfg.List(field)

		exp, ok := match[field]

		if !ok {
			if len(a) != 0 {
				t.Fatalf("Expected default value []string{} for (%s), got %v with len %v", field, a, len(a))
			}
		}

		if len(a) != len(exp) {
			t.Fatalf("Expected list of length %v, got %v", len(exp), len(a))
		}

		for i, e := range exp {
			if a[i] != e {
				t.Fatalf("Expected (%s) to resolve into %v, got %v", field, e, a)
			}
		}
	}
}
