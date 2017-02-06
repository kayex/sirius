package sirius

type Extension interface {
	Run(Message) (error, MessageAction)
}

type CfgExtension interface {
	RunWithConfig(Message, ExtensionConfig)
}

type cfg_int int
type cfg_bool bool
type cfg_float float64
type cfg_list []string

type EID string

type ExtensionLoader interface {
	Load(EID) Extension
}

type ExtensionConfig map[string]interface{}

type InvalidConfig struct {
	Key string
	msg string
}

func (ic InvalidConfig) Error() string {
	return ic.msg
}

func (cfg ExtensionConfig) Read(key string, def interface{}) interface{} {
	if val, ok := cfg[key]; ok {
		return val
	}

	return def
}

func (cfg ExtensionConfig) Integer(key string, def int) int {
	if val, ok := cfg[key]; ok {
		switch b := val.(type) {
		case int:
			return b
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
		switch l := val.(type) {
		case []string:
			return l
		}

	}

	return def
}
