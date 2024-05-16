package Validate

import (
	"regexp"
	"simple_account/app/Error"
)

func Username(username string) int {
	if username == "" {
		return Error.USERNAME_IS_EMPTY
	}
	usernameRegexp := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !usernameRegexp.MatchString(username) {
		return Error.INVALID_USERNAME
	}
	if len(username) < 4 {
		return Error.USERNAME_TOO_SHORT
	}
	if len(username) > 24 {
		return Error.USERNAME_TOO_LONG
	}

	return Error.NULL
}

func Email(email string) int {
	if email == "" {
		return Error.EMAIL_IS_EMPTY
	}
	emailRegexp := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	if !emailRegexp.MatchString(email) {
		return Error.INVALID_EMAIL_FORMAT
	}

	return Error.NULL
}

func Password(password string) int {
	if password == "" {
		return Error.PASSWORD_IS_EMPTY
	}
	passwordRegexp := regexp.MustCompile(`^[A-Za-z0-9!@#$%^&*_\-]+$`)
	if !passwordRegexp.MatchString(password) {
		return Error.INVALID_PASSWORD
	}
	if len(password) < 8 {
		return Error.PASSWORD_TOO_SHORT
	}

	return Error.NULL
}
