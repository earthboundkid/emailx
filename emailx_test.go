package emailx_test

import (
	"errors"
	"fmt"
	"runtime"
	"testing"

	"github.com/carlmjohnson/emailx"
)

var validitycases = []struct {
	email      string
	valid      bool
	resolvable bool
}{
	// Invalid format.
	{"", false, false},
	{"email@", false, false},
	{"email@x", false, false},
	{"email@@example.com", false, false},
	{".email@example.com", false, false},
	{"email.@example.com", false, false},
	{"email..test@example.com", false, false},
	{".email..test.@example.com", false, false},
	{"email@at@example.com", false, false},
	{"some whitespace@example.com", false, false},
	{"email@whitespace example.com", false, false},
	{`email'\@example.com`, false, false},
	{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@example.com", false, false},
	{"email@aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.com", false, false},
	{"emailaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@email.com", false, false},
	// Unresolvable domain.
	{"email+extra@wrong.example.com", true, false},
	{"{email+extra}@wrong.example.com", true, false},
	{
		"abcdefghijklmnopqrstuvwxyz" +
			"0123456789" +
			"!#$%&'*+/=?^_`{|}~.-]@t.d", true, false},
	{"emailaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@12345678012345678012345678012345678012345678012345678012345678012345678012345678012345678012345678012345678012345678012345678012345678012345678012345678012345678001234567800123456780123.com", true, false},
	{"0@0\x8000", false, false},
	{"0@\xe7̾\xb9\xf2\xd5", false, false},
	{"0@00000000000000000000000000000000", false, false},
	{"0@\xff\xb1\xb1\xb1\xb1\xb1\xb1\xff", false, false},
	{"0@000 00000", false, false},

	// Valid + resolvable
	{"{email}@gmail.com", true, true},
	{"email@gmail.com", true, true},
	{"email.email@gmail.com", true, true},
	{"email+extra@example.com", true, true},
	{"EMAIL@aol.co.uk", true, true},
	{"emailaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa@email.com", true, true},
	{"EMAIL+EXTRA@aol.co.uk", true, true},
}

func TestResolvable(t *testing.T) {
	for _, tt := range validitycases {
		if ok := emailx.Resolvable(tt.email); ok != tt.resolvable {
			t.Errorf("%q: got resolvable %t", tt.email, ok)
		}
	}
}

func TestResolve(t *testing.T) {
	for _, tt := range validitycases {
		err := emailx.Resolve(tt.email)
		if err == nil != tt.resolvable {
			t.Errorf("%q: unexpected error: %v", tt.email, err)
		}
	}
}

func TestValidate(t *testing.T) {
	for _, tt := range validitycases {
		err := emailx.Validate(tt.email)
		if err == nil != tt.valid {
			t.Errorf("%q: unexpected error: %v", tt.email, err)
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

var benchSink bool

func BenchmarkValidate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tt := validitycases[i%len(validitycases)]
		benchSink = emailx.Valid(tt.email)
		if benchSink != tt.valid {
			b.FailNow()
		}
	}
	runtime.KeepAlive(benchSink)
}

func FuzzValidate(f *testing.F) {
	for _, tt := range validitycases {
		f.Add(tt.email)
	}

	f.Fuzz(func(t *testing.T, email string) {
		emailx.Valid(email)
	})
}
