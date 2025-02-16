package jwtutil

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

func SigningMethodFromString(s string) (jwt.SigningMethod, error) {
	switch strings.ToLower(s) {
	case "hs256":
		return jwt.SigningMethodHS256, nil
	case "es256":
		return jwt.SigningMethodES256, nil
	}

	return nil, errors.New("unsupported signing method")
}
