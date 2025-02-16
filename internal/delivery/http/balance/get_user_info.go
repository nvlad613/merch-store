package balance

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
	"merch-store/internal/delivery"
	"merch-store/internal/domain/store"
	"merch-store/pkg/httputil"
	"net/http"
)

type GetUserInfoResponse struct {
	Coins       int                         `json:"coins"`
	Inventory   []delivery.InventoryItem    `json:"inventory"`
	CoinHistory delivery.TransactionsReport `json:"coinHistory"`
}

func (r *Router) GetUserInfoHandler(c echo.Context) error {
	ctx := c.Request().Context()

	userDetails, err := httputil.GetUserDetails(c)
	if err != nil {
		r.logger.Error("failed to get user details", "error", err)

		return httputil.SendError(http.StatusInternalServerError, "internal error", c)
	}

	report, err := r.balanceService.MakeReport(userDetails.Username, ctx)
	if err != nil {
		r.logger.Errorw(
			"failed to make report about users purchases and transactions",
			"error", err, "user", userDetails.Username, "userId", userDetails.Id,
		)

		return httputil.SendError(http.StatusInternalServerError, "internal error", c)
	}

	inventory, err := r.storeService.GetUserInventory(userDetails.Id, ctx)
	if err != nil {
		r.logger.Errorw(
			"failed to make get user inventory",
			"error", err, "user", userDetails.Username, "userId", userDetails.Id,
		)

		return httputil.SendError(http.StatusInternalServerError, "internal error", c)
	}

	responseBody := GetUserInfoResponse{
		Coins:       report.Coins,
		CoinHistory: new(delivery.TransactionsReport).FromModel(*report),
		Inventory: lo.Map(inventory, func(item store.InventoryItem, _ int) delivery.InventoryItem {
			return delivery.InventoryItem{
				Type:     item.ProductName,
				Quantity: item.Quantity,
			}
		}),
	}

	return c.JSON(http.StatusOK, responseBody)
}
