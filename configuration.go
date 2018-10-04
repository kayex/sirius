package sirius

import "net/url"

type ExtensionConfig map[string]interface{}

type Configuration struct {
	URL string
	EID EID
	Cfg ExtensionConfig
}

func FromConfigurationMap(cfg map[string]interface{}) []Configuration {
	var cfgs []Configuration

	for eid, settings := range cfg {
		var c Configuration

		// Check for HTTP extensions
		_, err := url.ParseRequestURI(eid)
		if err == nil {
			c = NewHTTPConfiguration(eid)
		} else {
			c = NewConfiguration(EID(eid))
		}

		switch ec := settings.(type) {
		case ExtensionConfig:
			c.Cfg = ec
		case map[string]interface{}:
			c.Cfg = ExtensionConfig(ec)
		}

		cfgs = append(cfgs, c)
	}

	return cfgs
}

func NewConfiguration(eid EID) Configuration {
	return Configuration{
		EID: eid,
	}
}

func NewHTTPConfiguration(url string) Configuration {
	return Configuration{
		URL: url,
	}
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
			switch b {
			case 0:
				return false
			case 1:
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
