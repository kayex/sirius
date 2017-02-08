package sirius

type Configuration struct {
	EID    EID
	Config map[string]interface{}
}

func NewConfiguration(eid EID) Configuration {
	return Configuration{
		EID:  eid,
	}
}
