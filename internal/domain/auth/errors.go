package auth

import "errors"

var (
	WrongCredentialsError = errors.New("authorization failed: wrong credentials")
)
