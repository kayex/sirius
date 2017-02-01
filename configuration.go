package sirius

type Configuration struct {
	ID     string
	Name   string
	User   *User
	EID    EID
	Config map[string]interface{}
}

func NewConfiguration(usr *User, eid EID) Configuration {
	return Configuration{
		User: usr,
		EID:  eid,
	}
}

func (c *Configuration) SetConfig(cfg map[string]interface{}) {
	c.Config = cfg
}
