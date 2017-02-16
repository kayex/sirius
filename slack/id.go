package slack

import (
	"crypto/sha256"
	"encoding/hex"
)

type ID interface {
	String() string
	Equals(ID) bool
}

type UserID struct {
	UserID string
	TeamID string
}

// SecureID is an opaque representation of a Slack user identity
// which can be used in place of UserID to lower the risk of
// accidentally printing/logging real (and confidential) Slack IDs.
//
// A SecureID can be constructed from a UserID
// by calling UserID.Secure()
type SecureID struct {
	HashSum string
}

// Equals indicates if id can be considered to be the same user identity as o.
func (id UserID) Equals(o ID) bool {
	switch o := o.(type) {
	case UserID:
		if id.Incomplete() || o.Incomplete() {
			return false
		}

		// Notice that user IDs are not guaranteed to be globally unique across all Slack users.
		// The combination of user ID and team ID, on the other hand, is guaranteed to be globally unique.
		//
		// - Slack API documentation
		return id.UserID == o.UserID && id.TeamID == o.TeamID
	case SecureID:
		return id.Secure().Equals(o)
	}

	return false
}

// Incomplete indicates if id has been populated properly.
func (id UserID) Incomplete() bool {
	return id.UserID == "" || id.TeamID == ""
}

func (id UserID) String() string {
	return id.TeamID + "." + id.UserID
}

// Secure converts id into a SecureID.
func (id UserID) Secure() SecureID {
	if id.Incomplete() {
		return SecureID{}
	}

	concat := id.TeamID + "." + id.UserID

	h := sha256.New()
	h.Write([]byte(concat))

	s := hex.EncodeToString(h.Sum(nil))

	return SecureID{
		HashSum: s,
	}
}

// Equals indicates if id can be considered to be the same user identity as o.
func (id SecureID) Equals(o ID) bool {
	switch o := o.(type) {
	case SecureID:
		if id.Incomplete() || o.Incomplete() {
			return false
		}

		return id.HashSum == o.HashSum
	case UserID:
		return o.Secure().Equals(id)
	}

	return false
}

// Incomplete indicates if id has been populated properly.
func (id *SecureID) Incomplete() bool {
	return id.HashSum == ""
}

func (id SecureID) String() string {
	return id.HashSum
}
