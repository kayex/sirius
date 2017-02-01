package sirius

import (
	"time"
)

type SlackID struct {
	UserID string
	TeamID string
}

func NewSlackID(userID, teamID string) SlackID {
	return SlackID{
		UserID: userID,
		TeamID: teamID,
	}
}

func (s *SlackID) equals(o *SlackID) bool {
	return s.UserID == o.UserID && s.TeamID == o.TeamID
}

type User struct {
	ID             string
	Token          string
	CreatedAt      time.Time
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
