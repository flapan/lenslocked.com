package models

import "strings"

const (
	// ErrNotFound is returned when a resource cannot be found in the database
	ErrNotFound modelError = "models: resource not found"
	// ErrIDInvalid is returned if an invalid id is supplied
	ErrIDInvalid modelError = "models: ID must be > 0"
	// ErrPasswordIncorrect is returned when a provided password is invalid
	ErrPasswordIncorrect modelError = "models: incorrect password provided"
	// ErrPasswordTooShort is returned when an update or create is attempted with a password that is less than 8 characters
	ErrPasswordTooShort modelError = "models: password must be at least 8 characters long"
	// ErrPasswordRequired is retunred when an update or create is attempted with the empty password
	ErrPasswordRequired modelError = "models: password is required"

	// ErrEmailRequired is returned when an email is not provided when creating a new user
	ErrEmailRequired modelError = "models: email address is required"
	// ErrEmailInvalid is returned when a provided email does not match emailRegExp
	ErrEmailInvalid modelError = "models: email address is not valid"
	// ErrEmailNotAvail is returned when upon user create the provided email is already used by another user in the system
	ErrEmailNotAvail modelError = "models: email is already used by a user"

	//ErrRememberHashRequired
	ErrRememberHashRequired modelError = "models: remember token is required"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}
