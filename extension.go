package sirius

type EID string
type ExtensionConfig map[string]interface{}

type Extension interface {
	Run(Message, ExtensionConfig) (MessageAction, error)
}

type ExtensionLoader interface {
	Load(EID) (Extension, error)
}

// Read fetches a value of any type for key.
// Returns def if key is not set.
func (cfg ExtensionConfig) Read(key string, def interface{}) interface{} {
	if val, ok := cfg[key]; ok {
		return val
	}
	return def
}

// String fetches a string value for key.
// Returns def if key is not set.
func (cfg ExtensionConfig) String(key string, def string) string {
	if val, ok := cfg[key]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return def
}

// Integer fetches an integer value for key.
// Returns def if key is not set.
func (cfg ExtensionConfig) Integer(key string, def int) int {
	if val, ok := cfg[key]; ok {
		if i, ok := val.(int); ok {
			return i
		}
	}
	return def
}

// Boolean fetches a boolean value for key.
// Returns false if key is not set.
func (cfg ExtensionConfig) Boolean(key string) bool {
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
	return false
}

// Float fetches a float value for key.
// Returns def if key is not set.
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

// List fetches a list value for key.
// Returns an empty list if key is not set.
func (cfg ExtensionConfig) List(key string) []string {
	var list []string

	if val, ok := cfg[key]; ok {
		switch l := val.(type) {
		case []interface{}:
			for _, lv := range l {
				if s, ok := lv.(string); ok {
					list = append(list, s)
				}
			}
			return list
		case []string:
			return l
		}
	}
	return []string{}
}
