package model

type Configuration struct {
	Id         string
	Name       string
	User       *User
	PluginGuid string
	Config     map[string]interface{}
}

func NewConfiguration(usr *User, pid string) Configuration {
	return Configuration{
		User:       usr,
		PluginGuid: pid,
	}
}

func (c *Configuration) SetConfig(cfg map[string]interface{}) {
	c.Config = cfg
}
