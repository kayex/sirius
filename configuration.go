package sirius

type Configuration struct {
	Id         string
	Name       string
	User       *User
	ExtensionGUID string
	Config     map[string]interface{}
}

func NewConfiguration(usr *User, pid string) Configuration {
	return Configuration{
		User:       usr,
		ExtensionGUID: pid,
	}
}

func (c *Configuration) SetConfig(cfg map[string]interface{}) {
	c.Config = cfg
}
