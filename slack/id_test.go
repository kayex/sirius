package slack

import "testing"

func TestUserID_Equals(t *testing.T) {
	cases := []struct {
		a   UserID
		b   UserID
		exp bool
	}{
		{
			a: UserID{
				UserID: "123",
				TeamID: "abc",
			},
			b: UserID{
				UserID: "123",
				TeamID: "abc",
			},
			exp: true,
		},
		{
			a: UserID{
				UserID: "456",
				TeamID: "abc",
			},
			b: UserID{
				UserID: "123",
				TeamID: "abc",
			},
			exp: false,
		},
		{
			a: UserID{
				UserID: "123",
				TeamID: "def",
			},
			b: UserID{
				UserID: "123",
				TeamID: "abc",
			},
			exp: false,
		},
		{
			a: UserID{
				UserID: "",
				TeamID: "",
			},
			b: UserID{
				UserID: "123",
				TeamID: "abc",
			},
			exp: false,
		},
		{
			a: UserID{
				UserID: "",
				TeamID: "",
			},
			b: UserID{
				UserID: "",
				TeamID: "",
			},
			exp: false,
		},
	}

	for _, c := range cases {
		if c.a.Equals(c.b) != c.exp {
			if c.exp {
				t.Fatalf("Expected (%v) to equal (%v)", c.a.String(), c.b.String())
			} else {
				t.Fatalf("Expected (%v) to not equal (%v)", c.a.String(), c.b.String())
			}
		}
	}
}

func TestUserID_Incomplete(t *testing.T) {
	cases := []struct {
		id  UserID
		exp bool
	}{
		{
			id: UserID{
				UserID: "",
				TeamID: "",
			},
			exp: true,
		},
		{
			id: UserID{
				UserID: "123",
				TeamID: "",
			},
			exp: true,
		},
		{
			id: UserID{
				UserID: "",
				TeamID: "abc",
			},
			exp: true,
		},
		{
			id: UserID{
				UserID: "123",
				TeamID: "abc",
			},
			exp: false,
		},
	}

	for _, c := range cases {
		if c.id.Incomplete() != c.exp {
			t.Fatalf("Expected UserID (%v) to count as Incomplete()", c.id.String())
		}
	}

}

func TestUserID_String(t *testing.T) {
	id := UserID{
		UserID: "123",
		TeamID: "abc",
	}
	exp := "abc.123"

	if id.String() != exp {
		t.Fatalf("Expected UserID to concatenate into (%v), got (%v)", exp, id.String())
	}

}

func TestUserID_Secure(t *testing.T) {
	cases := []struct {
		id UserID
	}{
		{
			id: UserID{
				UserID: "123",
				TeamID: "abc",
			},
		},
		{
			id: UserID{
				UserID: "",
				TeamID: "abc",
			},
		},
		{
			id: UserID{
				UserID: "123",
				TeamID: "",
			},
		},
		{
			id: UserID{
				UserID: "",
				TeamID: "",
			},
		},
	}

	for _, c := range cases {
		sid := c.id.Secure()

		if sid.Incomplete() {
			t.Fatalf("UserID (%v) secures into (%v), which is not a valid SecureID", c.id.String(), sid.String())
		}
	}
}

func TestSecureID_Incomplete(t *testing.T) {
	cases := []struct {
		id  SecureID
		exp bool
	}{
		{
			id:  SecureID{"test-id"},
			exp: false,
		},
		{
			id:  SecureID{},
			exp: true,
		},
	}

	for _, c := range cases {
		if c.id.Incomplete() != c.exp {
			t.Fatalf("Expected SecureID (%v) to count as Incomplete()", c.id.String())
		}
	}

}

func TestSecureID_Equals(t *testing.T) {
	cases := []struct {
		a   SecureID
		b   SecureID
		exp bool
	}{
		{
			a:   SecureID{"Han Solo"},
			b:   SecureID{"Han Solo"},
			exp: true,
		},
		{
			a:   SecureID{"Han Solo"},
			b:   SecureID{"Darth Vader"},
			exp: false,
		},
		{
			a:   SecureID{},
			b:   SecureID{"Darth Vader"},
			exp: false,
		},
		{
			a:   SecureID{},
			b:   SecureID{},
			exp: false,
		},
	}

	for _, c := range cases {
		if c.a.Equals(c.b) != c.exp {
			if c.exp {
				t.Fatalf("Expected SecureID (%v) to equal (%v)", c.a.String(), c.b.String())
			} else {
				t.Fatalf("Expected SecureID (%v) to not equal (%v)", c.a.String(), c.b.String())
			}
		}
	}
}

func TestSecureID_String(t *testing.T) {
	id := SecureID{"test-1"}
	exp := "test-1"

	if id.String() != exp {
		t.Fatalf("Expected SecureID with HashSum (%v) to String() into (%v), got (%v)", id.HashSum, exp, id.String())
	}
}
