// Package emailx contains helpers for email address validation and normalization.
package emailx

import (
	"errors"
	"net"
	"strings"
)

var (
	// ErrInvalidFormat is returned when email's format is invalid
	ErrInvalidFormat = errors.New("invalid format")
	// ErrUnresolvableHost is returned when Resolve couldn't resolve email's host
	ErrUnresolvableHost = errors.New("unresolvable host")
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

var isValidUser = func() func(s string) bool {
	const validChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789" +
		"!#$%&'*+/=?^_`{|}~.-]+$"
	var m [256]bool
	for _, c := range validChars {
		m[c] = true
	}
	return func(s string) bool {
		for _, b := range []byte(s) {
			if !m[b] {
				return false
			}
		}
		return true
	}
}()

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
		// As per RFC 5332 section 3.2.3:
		// https://tools.ietf.org/html/rfc5322#section-3.2.3
		!isValidUser(user),
		// Dots are not allowed in the beginning, end
		// or in groups of more than 1 in the user address
		strings.HasPrefix(user, "."),
		strings.HasSuffix(user, "."),
		strings.Contains(user, ".."),
		// No whitespace in host
		strings.ContainsAny(host, "\t\n\f\r "),
		// Host must contain .
		!strings.Contains(host, "."):
		return ErrInvalidFormat
	}

	return nil
}

// Split an address into user and host portions.
// Split does not perform any validation or normalization.
func Split(email string) (user, host string) {
	at := strings.LastIndex(email, "@")
	if at == -1 {
		return
	}

	user = email[:at]
	host = email[at+1:]
	return
}

// Normalize an email address.
func Normalize(email string) string {
	// Trim whitespaces.
	email = strings.TrimSpace(email)

	// Trim extra dot in hostname.
	email = strings.TrimRight(email, ".")

	// Lowercase.
	email = strings.ToLower(email)

	return email
}
