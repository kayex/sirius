package model

type Configuration struct {
	Id         string
	Name       string
	UserId     string
	User       *User
	PluginGuid string
	Config     map[string]interface{}
}

func NewConfiguration(usr *User, pid string) Configuration {
	return Configuration{
		User:       usr,
		UserId:     usr.Id,
		PluginGuid: pid,
	}
}

func (c *Configuration) SetConfig(cfg map[string]interface{}) {
	c.Config = cfg
}
