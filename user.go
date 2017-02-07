package sirius

type User struct {
	ID             SlackID
	Token          string
	Configurations []*Configuration
}

func NewUser(token string) User {
	return User{
		Token: token,
	}
}

func (u *User) AddConfiguration(cfg *Configuration) {
	u.Configurations = append(u.Configurations, cfg)
}
