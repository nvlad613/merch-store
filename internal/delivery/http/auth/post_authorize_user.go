package auth

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
	"merch-store/internal/domain/auth"
	"merch-store/pkg/httputil"
	"net/http"
)

type PostAuthorizeUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PostAuthorizeUserResponse struct {
	Token string `json:"token"`
}

func (req PostAuthorizeUserRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.Password, validation.Required, validation.Length(8, 64)),
		validation.Field(&req.Username, validation.Required, validation.Length(1, 32)),
	)
}

func (r *Router) PostAuthorizeUserHandler(c echo.Context) error {
	ctx := c.Request().Context()

	var requestBody PostAuthorizeUserRequest
	if c.Bind(&requestBody) != nil {
		return httputil.SendError(http.StatusBadRequest, "invalid request body", c)
	}

	if err := requestBody.Validate(); err != nil {
		return httputil.SendError(http.StatusBadRequest, err.Error(), c)
	}

	userModel := auth.User{
		Username: requestBody.Username,
		Password: requestBody.Password,
	}
	token, err := r.authService.MakeAuth(userModel, ctx)
	switch {
	case errors.Is(err, auth.WrongCredentialsError):
		return httputil.SendError(http.StatusForbidden, err.Error(), c)
	case err != nil:
		return httputil.SendError(http.StatusInternalServerError, "internal error", c)
	}

	return c.JSON(http.StatusOK, &PostAuthorizeUserResponse{
		Token: token,
	})
}
