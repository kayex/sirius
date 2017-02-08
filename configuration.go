package sirius

type Configuration struct {
	EID    EID
	Config ExtensionConfig
}

func NewConfiguration(eid EID) Configuration {
	return Configuration{
		EID: eid,
	}
}
