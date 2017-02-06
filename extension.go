package sirius

type Extension interface {
	Run(Message) (error, MessageAction)
}

type CfgExtension interface {
	Run(Message, ExtensionConfig) (error, MessageAction)
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

func (cfg ExtensionConfig) read(key string, def interface{}) interface{} {
	if val, ok := cfg[key]; ok {
		return val
	}

	return def
}

func (cfg ExtensionConfig) integer(key string, def int) cfg_int {
	if val, ok := cfg[key]; ok {
		switch b := val.(type) {
		case int:
			return b
		}
	}

	return def
}

func (cfg ExtensionConfig) boolean(key string, def bool) cfg_bool {
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

func (cfg ExtensionConfig) float(key string, def float64) cfg_float {
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

func (cfg ExtensionConfig) list(key string, def []string) cfg_list {
	if val, ok := cfg[key]; ok {
		switch l := val.(type) {
		case []string:
			return l
		}

	}

	return def
}
