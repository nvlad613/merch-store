package store

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
	"merch-store/internal/domain/store"
	"merch-store/pkg/httputil"
	"net/http"
)

type GetBuyMerchRequest struct {
	ProductName string `param:"item"`
}

func (req GetBuyMerchRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(&req.ProductName, validation.Required, validation.Length(1, 32)),
	)
}

func (r *Router) GetBuyMerchHandler(c echo.Context) error {
	ctx := c.Request().Context()

	userDetails, err := httputil.GetUserDetails(c)
	if err != nil {
		r.logger.Error("failed to get user details", "error", err)

		return httputil.SendError(http.StatusInternalServerError, "internal error", c)
	}

	var requestBody GetBuyMerchRequest
	if c.Bind(&requestBody) != nil {
		return httputil.SendError(http.StatusBadRequest, "invalid path param", c)
	}

	if err := requestBody.Validate(); err != nil {
		return httputil.SendError(http.StatusBadRequest, err.Error(), c)
	}

	err = r.storeService.MakePurchase(
		userDetails.Username,
		requestBody.ProductName,
		1,
		ctx,
	)
	switch {
	case errors.Is(err, store.MerchItemNotFound):
		return httputil.SendError(http.StatusUnprocessableEntity, err.Error(), c)
	case err != nil:
		r.logger.Errorw(
			"failed to make purchase",
			"error", err, "user", userDetails.Username, "item", requestBody.ProductName,
		)

		return httputil.SendError(http.StatusInternalServerError, "internal error", c)
	}

	return c.NoContent(http.StatusOK)
}
