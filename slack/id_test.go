package slack

import "testing"

func TestUserID_Equals(t *testing.T) {
	cases := []struct {
		a   UserID
		b   ID
		exp bool
	}{
		{
			a:   UserID{"123", "abc"},
			b:   UserID{"123", "abc"},
			exp: true,
		},
		{
			a:   UserID{"456", "abc"},
			b:   UserID{"123", "abc"},
			exp: false,
		},
		{
			a:   UserID{"123", "def"},
			b:   UserID{"123", "abc"},
			exp: false,
		},
		{
			a:   UserID{"", ""},
			b:   UserID{"123", "abc"},
			exp: false,
		},
		{
			a:   UserID{"", ""},
			b:   UserID{"", ""},
			exp: false,
		},
		{
			a:   UserID{"123", "abc"},
			b:   UserID{"123", "abc"}.Secure(),
			exp: true,
		},
		{
			a:   UserID{"456", "abc"},
			b:   UserID{"123", "abc"}.Secure(),
			exp: false,
		},
		{
			a:   UserID{"123", "def"},
			b:   UserID{"123", "abc"}.Secure(),
			exp: false,
		},
		{
			a:   UserID{"", ""},
			b:   UserID{"123", "abc"}.Secure(),
			exp: false,
		},
		{
			a:   UserID{"", ""},
			b:   UserID{"", ""}.Secure(),
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

func TestUserID_Valid(t *testing.T) {
	cases := []struct {
		id  UserID
		exp bool
	}{
		{
			id:  UserID{"", ""},
			exp: false,
		},
		{
			id:  UserID{"123", ""},
			exp: false,
		},
		{
			id:  UserID{"", "abc"},
			exp: false,
		},
		{
			id:  UserID{"123", "abc"},
			exp: true,
		},
	}

	for _, c := range cases {
		v := c.id.Valid()

		if v != c.exp {
			t.Fatalf("Expected UserID(%v).Valid() to be %v, got %v", c.id.String(), c.exp, v)
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
		id  UserID
		exp bool
	}{
		{
			id:  UserID{"123", "abc"},
			exp: true,
		},
		{
			id:  UserID{"", "abc"},
			exp: false,
		},
		{
			id:  UserID{"123", ""},
			exp: false,
		},
		{
			id:  UserID{"", ""},
			exp: false,
		},
	}

	for _, c := range cases {
		sid := c.id.Secure()

		if sid.Valid() != c.exp {
			t.Fatalf("UserID (%v) secures into (%v), which is not a valid SecureID", c.id.String(), sid.String())
		}
	}
}

func TestSecureID_Valid(t *testing.T) {
	cases := []struct {
		id  SecureID
		exp bool
	}{
		{
			id:  SecureID{"test-id"},
			exp: true,
		},
		{
			id:  SecureID{},
			exp: false,
		},
	}

	for _, c := range cases {
		v := c.id.Valid()
		if v != c.exp {
			t.Fatalf("Expected SecureID(%v).Valid() to be %v, got %v", c.id.String(), c.exp, v)
		}
	}

}

func TestSecureID_Equals(t *testing.T) {
	cases := []struct {
		a   SecureID
		b   ID
		exp bool
	}{
		{
			a:   UserID{"123", "abc"}.Secure(),
			b:   UserID{"123", "abc"}.Secure(),
			exp: true,
		},
		{
			a:   UserID{"123", "abc"}.Secure(),
			b:   UserID{"456", "abc"}.Secure(),
			exp: false,
		},
		{
			a:   UserID{"123", "abc"}.Secure(),
			b:   UserID{"123", "def"}.Secure(),
			exp: false,
		},
		{
			a:   UserID{}.Secure(),
			b:   UserID{}.Secure(),
			exp: false,
		},
		{
			a:   UserID{"123", "abc"}.Secure(),
			b:   UserID{"123", "abc"},
			exp: true,
		},
		{
			a:   UserID{"123", "abc"}.Secure(),
			b:   UserID{"456", "abc"},
			exp: false,
		},
		{
			a:   UserID{"123", "abc"}.Secure(),
			b:   UserID{"123", "def"},
			exp: false,
		},
		{
			a:   UserID{"123", "abc"}.Secure(),
			b:   UserID{"", ""},
			exp: false,
		},
		{
			a:   UserID{"", ""}.Secure(),
			b:   UserID{"", ""},
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
