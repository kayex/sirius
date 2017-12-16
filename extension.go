package sirius

type EID string

type Extension interface {
	Run(Message, ExtensionConfig) (MessageAction, error)
}

type ExtensionLoader interface {
	Load(EID) (Extension, error)
}

// ConfigExtension represents an extension with an associated configuration.
type ConfigExtension struct {
	Extension
	Cfg ExtensionConfig
}

func NewConfigExtension(ex Extension, cfg ExtensionConfig) *ConfigExtension {
	return &ConfigExtension{
		Extension: ex,
		Cfg: cfg,
	}
}

func FromConfiguration(l ExtensionLoader, cfg *Configuration) (*ConfigExtension, error) {
	// Check for HTTP extensions
	if cfg.URL != "" {
		ex := NewHttpExtension(cfg.URL, nil)
		return NewConfigExtension(ex, cfg.Cfg), nil
	} else {
		ex, err := l.Load(cfg.EID)
		if err != nil {
			return nil, err
		}

		return NewConfigExtension(ex, cfg.Cfg), nil
	}
}

func (ex *ConfigExtension) Run(msg Message) (MessageAction, error) {
	return ex.Extension.Run(msg, ex.Cfg)
}

func LoadFromSettings(l ExtensionLoader, s Settings) ([]ConfigExtension, error) {
	var exe []ConfigExtension

	for _, cfg := range s {
		x, err := FromConfiguration(l, &cfg)
		if err != nil {
			return nil, err
		}

		exe = append(exe, *x)
	}

	return exe, nil
}
