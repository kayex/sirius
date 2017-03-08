package slack

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
)

type ID interface {
	String() string
	Equals(ID) bool
	Valid() bool
}

type UserID struct {
	UserID string
	TeamID string
}

// SecureID is an opaque, deterministic representation of a Slack user identity
// that can be used in place of UserID to reduce the risk of compromising
// a user's real identity.
//
// A SecureID can be constructed from a UserID
// by calling UserID.Secure()
type SecureID struct {
	HashSum string
}

// Equals indicates if id and o represent the same user identity.
func (id UserID) Equals(o ID) bool {
	switch o := o.(type) {
	case UserID:
		if !(id.Valid() && o.Valid()) {
			return false
		}

		// > Notice that user IDs are not guaranteed to be globally unique across all Slack users.
		// > The combination of user ID and team ID, on the other hand, is guaranteed to be globally unique.
		//
		// - Slack API documentation
		return id.UserID == o.UserID && id.TeamID == o.TeamID
	case SecureID:
		return id.Secure().Equals(o)
	}

	return false
}

func (id UserID) Valid() bool {
	return id.UserID != "" && id.TeamID != ""
}

func (id UserID) String() string {
	return id.TeamID + "." + id.UserID
}

// Secure converts id into a SecureID.
func (id UserID) Secure() SecureID {
	if !id.Valid() {
		return SecureID{}
	}

	var buf bytes.Buffer
	buf.WriteString(id.TeamID)
	buf.WriteRune('.')
	buf.WriteString(id.UserID)

	h := sha256.New()
	h.Write(buf.Bytes())
	s := hex.EncodeToString(h.Sum(nil))

	return SecureID{
		HashSum: s,
	}
}

// Equals indicates if id and o represent the same user identity.
func (id SecureID) Equals(o ID) bool {
	switch o := o.(type) {
	case SecureID:
		if !id.Valid() || !o.Valid() {
			return false
		}

		return id.HashSum == o.HashSum
	case UserID:
		return o.Secure().Equals(id)
	}
	return false
}

func (id SecureID) Valid() bool {
	return id.HashSum != ""
}

func (id SecureID) String() string {
	return id.HashSum
}
