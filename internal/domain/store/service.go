package store

import "context"

type Service interface {
	MakePurchase(username string, productName string, quantity int, ctx context.Context) error
	GetUserInventory(username string, ctx context.Context) ([]InventoryItem, error)
}
