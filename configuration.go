package sirius

type Configuration struct {
	EID EID
	Cfg ExtensionConfig
}

func NewConfiguration(eid EID) Configuration {
	return Configuration{
		EID: eid,
	}
}
