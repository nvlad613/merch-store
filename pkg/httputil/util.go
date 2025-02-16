package httputil

import (
	"errors"
	"fmt"
	"github.com/go-viper/mapstructure/v2"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserDetails struct {
	Username string `mapstructure:"username"`
	Id       int    `mapstructure:"user_id"`
}

func SendError(status int, details string, c echo.Context) error {
	return c.JSON(status, echo.Map{
		"error": details,
	})
}

func GetUserDetails(c echo.Context) (UserDetails, error) {
	var zero UserDetails

	// by default token is stored under `user` key
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return zero, errors.New("get jwt token claims")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return zero, errors.New("cast claims as jwt.MapClaims")
	}

	var userDetails UserDetails
	if err := mapstructure.Decode(claims, &userDetails); err != nil {
		return zero, fmt.Errorf("get user from jwt claims: %w", err)
	}

	return userDetails, nil
}

func JwtErrorHandler(ctx echo.Context, err error) error {
	var details = "authorization error"
	switch {
	case errors.Is(err, echojwt.ErrJWTMissing):
		details = "jwt is missing"
	case errors.Is(err, echojwt.ErrJWTInvalid):
		details = "jwt is invalid"
	}

	return SendError(http.StatusUnauthorized, details, ctx)
}
