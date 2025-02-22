package balance

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
	"merch-store/internal/domain"
	"merch-store/internal/domain/balance"
	"merch-store/pkg/httputil"
	"net/http"
)

type PostSendCoinsRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

func (req PostSendCoinsRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.ToUser, validation.Required, validation.Length(1, 32)),
		validation.Field(&req.Amount, validation.Required, validation.Min(1)),
	)
}

func (r *Router) PostSendCoinsHandler(c echo.Context) error {
	ctx := c.Request().Context()

	userDetails, err := httputil.GetUserDetails(c)
	if err != nil {
		r.logger.Error("failed to get user details", "error", err)

		return httputil.SendError(http.StatusInternalServerError, "internal error", c)
	}

	var requestBody PostSendCoinsRequest
	if c.Bind(&requestBody) != nil {
		return httputil.SendError(http.StatusBadRequest, "invalid request body", c)
	}

	if err := requestBody.Validate(); err != nil {
		return httputil.SendError(http.StatusBadRequest, err.Error(), c)
	}

	err = r.balanceService.MakeTransaction(
		userDetails.Username,
		requestBody.ToUser,
		requestBody.Amount,
		ctx,
	)
	switch {
	case errors.Is(err, balance.NotEnoughCoinsError), errors.Is(err, domain.UserNotFoundError):
		return httputil.SendError(http.StatusUnprocessableEntity, err.Error(), c)
	case err != nil:
		r.logger.Errorw(
			"failed to make transaction",
			"error", err, "fromUser", userDetails.Username, "toUser", requestBody.ToUser,
		)

		return httputil.SendError(http.StatusInternalServerError, "internal error", c)
	}

	return c.NoContent(http.StatusOK)
}
