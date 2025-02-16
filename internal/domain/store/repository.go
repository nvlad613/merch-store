package store

import "context"

type Repository interface {
	MakePurchase(purchase Purchase, username string, ctx context.Context) error
	GetUserPurchases(userId int, ctx context.Context) ([]Purchase, error)
}
