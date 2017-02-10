package sirius

type EID string
type ExtensionConfig map[string]interface{}

type Extension interface {
	Run(Message, ExtensionConfig) (MessageAction, error)
}

type ExtensionLoader interface {
	Load(EID) (Extension, error)
}

func (cfg ExtensionConfig) Read(key string, def interface{}) interface{} {
	if val, ok := cfg[key]; ok {
		return val
	}

	return def
}

func (cfg ExtensionConfig) String(key string, def string) string {
	if val, ok := cfg[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}

	return def
}

func (cfg ExtensionConfig) Integer(key string, def int) int {
	if val, ok := cfg[key]; ok {
		if i, ok := val.(int); ok {
			return i
		}
	}

	return def
}

func (cfg ExtensionConfig) Boolean(key string, def bool) bool {
	if val, ok := cfg[key]; ok {
		switch b := val.(type) {
		case bool:
			return b
		case int:
			// Require explicit 0 or 1
			if b == 0 {
				return false
			} else if b == 1 {
				return true
			}
		}
	}

	return def
}

func (cfg ExtensionConfig) Float(key string, def float64) float64 {
	if val, ok := cfg[key]; ok {
		switch f := val.(type) {
		case float32:
			return float64(f)
		case float64:
			return f
		}
	}

	return def
}

func (cfg ExtensionConfig) List(key string, def []string) []string {
	if val, ok := cfg[key]; ok {
		if l, ok := val.([]string); ok {
			return l
		}
	}

	return def
}
