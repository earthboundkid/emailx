# emailx [![GoDoc](https://godoc.org/github.com/carlmjohnson/emailx?status.svg)](https://godoc.org/github.com/carlmjohnson/emailx) [![Calver v0.YY.Minor](https://img.shields.io/badge/calver-v0.YY.Minor-22bfda.svg)](https://calver.org)

Go package for email validation and normalization.

Forked from [goware/emailx](https://github.com/goware/emailx) with some breaking changes to make the API more convenient.

## Email validation

Simple email format check (not a complicated regexp, [this is why](http://davidcel.is/posts/stop-validating-email-addresses-with-regex/)).

```go
import "github.com/carlmjohnson/emailx"

func main() {
    if email := "email.@example.com"; !emailx.Valid(email) {
        fmt.Printf("%q is not valid\n", email)
        // "email.@example.com" is not valid
    }
}
```

## Email resolving

Check whether the domain has a valid DNS record:

```go
    if err := emailx.Resolve("My+Email@wrong.example.com"); err != nil {
        fmt.Println("Email is not valid.")

        if errors.Is(err, emailx.ErrUnresolvableHost) {
            fmt.Println("Unresolvable host.")
        }
    }
    // Output:
    // Email is not valid.
    // Unresolvable host.
```

## Email normalization

```go
import "github.com/carlmjohnson/emailx"

func main() {
    fmt.Print(emailx.Normalize(" My+Email@example.com. "))
    // Prints my+email@example.com
}
```

## License
Emailx is licensed under the [MIT License](./LICENSE).
