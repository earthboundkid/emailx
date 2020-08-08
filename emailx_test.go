package emailx_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/carlmjohnson/emailx"
)

func TestResolvable(t *testing.T) {
	tests := []struct {
		in  string
		out string
		err bool
	}{
		// Invalid format.
		{in: "", err: true},
		{in: "email@", err: true},
		{in: "email@x", err: true},
		{in: "email@@example.com", err: true},
		{in: ".email@example.com", err: true},
		{in: "email.@example.com", err: true},
		{in: "email..test@example.com", err: true},
		{in: ".email..test.@example.com", err: true},
		{in: "email@at@example.com", err: true},
		{in: "some whitespace@example.com", err: true},
		{in: "email@whitespace example.com", err: true},
		{in: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@example.com", err: true},
		{in: "email@aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.com", err: true},

		// Unresolvable domain.
		{in: "email+extra@wrong.example.com", err: true},

		// Valid.
		{in: "email@gmail.com"},
		{in: "email.email@gmail.com"},
		{in: "email+extra@example.com"},
		{in: "EMAIL@aol.co.uk"},
		{in: "EMAIL+EXTRA@aol.co.uk"},
	}

	for _, tt := range tests {
		if ok := emailx.Resolvable(tt.in); ok == tt.err {
			t.Errorf("%q: got resolvable %t", tt.in, ok)
		}
	}
}

func TestResolve(t *testing.T) {
	tests := []struct {
		in  string
		out string
		err bool
	}{
		// Invalid format.
		{in: "", err: true},
		{in: "email@", err: true},
		{in: "email@x", err: true},
		{in: "email@@example.com", err: true},
		{in: ".email@example.com", err: true},
		{in: "email.@example.com", err: true},
		{in: "email..test@example.com", err: true},
		{in: ".email..test.@example.com", err: true},
		{in: "email@at@example.com", err: true},
		{in: "some whitespace@example.com", err: true},
		{in: "email@whitespace example.com", err: true},
		{in: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@example.com", err: true},
		{in: "email@aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.com", err: true},

		// Unresolvable domain.
		{in: "email+extra@wrong.example.com", err: true},

		// Valid.
		{in: "email@gmail.com"},
		{in: "email.email@gmail.com"},
		{in: "email+extra@example.com"},
		{in: "EMAIL@aol.co.uk"},
		{in: "EMAIL+EXTRA@aol.co.uk"},
	}

	for _, tt := range tests {
		err := emailx.Resolve(tt.in)
		if err != nil {
			if !tt.err {
				t.Errorf(`"%s": unexpected error \"%v\"`, tt.in, err)
			}
			continue
		}
		if tt.err && err == nil {
			t.Errorf(`"%s": expected error`, tt.in)
			continue
		}
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		in  string
		out string
		err bool
	}{
		// Invalid format.
		{in: "", err: true},
		{in: "email@", err: true},
		{in: "email@x", err: true},
		{in: "email@@example.com", err: true},
		{in: ".email@example.com", err: true},
		{in: "email.@example.com", err: true},
		{in: "email..test@example.com", err: true},
		{in: ".email..test.@example.com", err: true},
		{in: "email@at@example.com", err: true},
		{in: "some whitespace@example.com", err: true},
		{in: "email@whitespace example.com", err: true},
		{in: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@example.com", err: true},
		{in: "email@aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.com", err: true},

		// Valid.
		{in: "email@gmail.com"},
		{in: "email.email@gmail.com"},
		{in: "email+extra@example.com"},
		{in: "EMAIL@aol.co.uk"},
		{in: "EMAIL+EXTRA@aol.co.uk"},
	}

	for _, tt := range tests {
		err := emailx.Validate(tt.in)
		if err != nil {
			if !tt.err {
				t.Errorf(`"%s": unexpected error \"%v\"`, tt.in, err)
			}
			continue
		}
		if tt.err && err == nil {
			t.Errorf(`"%s": expected error`, tt.in)
			continue
		}
	}
}

func ExampleValid() {
	if email := "email.@example.com"; !emailx.Valid(email) {
		fmt.Printf("%q is not valid\n", email)
	}
	// Output:
	// "email.@example.com" is not valid
}

func ExampleResolve() {
	if err := emailx.Resolve("My+Email@wrong.example.com"); err != nil {
		fmt.Println("Email is not valid.")

		if errors.Is(err, emailx.ErrInvalidFormat) {
			fmt.Println("Wrong format.")
		}

		if errors.Is(err, emailx.ErrUnresolvableHost) {
			fmt.Println("Unresolvable host.")
		}
	}
	// Output:
	// Email is not valid.
	// Unresolvable host.
}

func TestNormalize(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{in: "email@EXAMPLE.COM. ", out: "email@example.com"},
		{in: " Email+Me@example.com. ", out: "email+me@example.com"},
	}

	for _, tt := range tests {
		normalized := emailx.Normalize(tt.in)
		if normalized != tt.out {
			t.Errorf(`%v: got "%v", want "%v"`, tt.in, normalized, tt.out)
		}
	}
}

func ExampleNormalize() {
	fmt.Println(emailx.Normalize(" Email+Me@example.com. "))
	// Output: email+me@example.com
}

func TestSplit(t *testing.T) {
	tests := []struct {
		in, user, host string
	}{
		{"", "", ""},
		{"user@", "user", ""},
		{"@host", "", "host"},
		{"user@host", "user", "host"},
		{"user@subuser@host", "user@subuser", "host"},
	}

	for _, tt := range tests {
		user, host := emailx.Split(tt.in)
		if tt.user != user {
			t.Errorf("%q user %q != %q", tt.in, tt.user, user)
		}
		if tt.host != host {
			t.Errorf("%q host %q != %q", tt.in, tt.host, host)
		}
	}
}
