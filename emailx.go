package emailx

import (
	"errors"
	"net"
	"regexp"
	"strings"
)

var (
	//ErrInvalidFormat returns when email's format is invalid
	ErrInvalidFormat = errors.New("invalid format")
	//ErrUnresolvableHost returns when validator couldn't resolve email's host
	ErrUnresolvableHost = errors.New("unresolvable host")

	userRegexp = regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$")
	hostRegexp = regexp.MustCompile("^[^\\s]+\\.[^\\s]+$")
	// As per RFC 5332 secion 3.2.3: https://tools.ietf.org/html/rfc5322#section-3.2.3
	// Dots are not allowed in the beginning, end or in occurances of more than 1 in the email address
	userDotRegexp = regexp.MustCompile("(^[.]{1})|([.]{1}$)|([.]{2,})")
)

// Resolvable reports if Resolve succeeds
func Resolvable(email string) bool {
	return Resolve(email) == nil
}

// Resolve checks the validity of a given address and resolves its host name.
func Resolve(email string) error {
	if err := Validate(email); err != nil {
		return err
	}
	_, host := Split(email)

	if _, err := net.LookupMX(host); err != nil {
		if _, err := net.LookupIP(host); err != nil {
			// Only fail if both MX and A records are missing - any of the
			// two is enough for an email to be deliverable
			return ErrUnresolvableHost
		}
	}

	return nil
}

// Valid reports if Validate succeeds.
func Valid(email string) bool {
	return Validate(email) == nil
}

// Validate checks format of a given email.
func Validate(email string) error {
	if len(email) < 6 || len(email) > 254 {
		return ErrInvalidFormat
	}

	user, host := Split(email)
	switch {
	case len(user) < 1,
		len(user) > 64,
		len(host) < 3,
		userDotRegexp.MatchString(user),
		!userRegexp.MatchString(user),
		!hostRegexp.MatchString(host):
		return ErrInvalidFormat
	}

	return nil
}

// Split attempts to split an address into user and host portions.
// Split does not perform any validation.
func Split(email string) (user, host string) {
	at := strings.LastIndex(email, "@")
	if at == -1 {
		return
	}

	user = email[:at]
	host = email[at+1:]
	return
}

// Normalize normalizes email address.
func Normalize(email string) string {
	// Trim whitespaces.
	email = strings.TrimSpace(email)

	// Trim extra dot in hostname.
	email = strings.TrimRight(email, ".")

	// Lowercase.
	email = strings.ToLower(email)

	return email
}
