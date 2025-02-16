package store

import "context"

type Repository interface {
	MakePurchase(purchase Purchase, ctx context.Context) error
	GetUserPurchases(username string, ctx context.Context) ([]Purchase, error)
}
