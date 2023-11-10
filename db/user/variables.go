package user

import (
	"errors"
)

// errors that might be returned by the functions in this package
var (
	ErrEmailAndDisplayNameTaken = errors.New("that email and display name are both already in use")
	ErrEmailTaken               = errors.New("that email is already in use")
	ErrDisplayNameTaken         = errors.New("that display name is already in use")
	ErrDatabaseError            = errors.New("there was an error when executing a query on the database")
	ErrPwdTooLong               = errors.New("the password is exceeding the 72-byte limit")
	ErrNoAuthCookie             = errors.New("the request contained no authentication cookie")
	ErrExpiredSession           = errors.New("the request contained an expired authentication cookie")
	ErrUnrecognizedSession      = errors.New("the session cookie wasn't found in the database")
	ErrServer                   = errors.New("the server had some kind of unexpected error")
	ErrEmailNotFound            = errors.New("no user matching that email was found")
	ErrIncorrectPassword        = errors.New("incorrect password")
)

// this is the max amount of time that a session is allowed before it is
// considered to have expired. This constant represents 4 days in nanoseconds.
const maxSessionAge = 3.456e+14

// represents a session as it appears in the database.
type Session struct {
	Email       string
	LastRenewed string
	Uuid        string
}

// The email and password of a user; the `Password` should be unhashed.
type UserCredentials struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// A struct that represents all of the data necessary to register a new user.
type UserRegistration struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required"`
	FirstName   string `json:"firstName" validate:"required"`
	LastName    string `json:"lastName" validate:"required"`
	DisplayName string `json:"displayName" validate:"required"`
}
