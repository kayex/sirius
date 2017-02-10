package slack

import (
	"crypto/sha256"
	"encoding/hex"
)

type UserID struct {
	UserID string
	TeamID string
}

type SecureID struct {
	HashSum string
}

/*
Notice that user IDs are not guaranteed to be globally unique across all Slack users.
The combination of user ID and team ID, on the other hand, is guaranteed to be globally unique.

- Slack API documentation
*/
func (id UserID) Equals(o UserID) bool {
	return id.UserID == o.UserID && id.TeamID == o.TeamID
}

func (id UserID) Empty() bool {
	return id.UserID == "" && id.TeamID == ""
}

func (id UserID) Secure() SecureID {
	concat := id.TeamID + "." + id.UserID

	h := sha256.New()
	h.Write([]byte(concat))

	s := hex.EncodeToString(h.Sum(nil))

	return SecureID{
		HashSum: s,
	}
}

func (id SecureID) Equals(o SecureID) bool {
	return id.HashSum == o.HashSum
}

func (id *SecureID) Empty() bool {
	return id.HashSum == ""
}
