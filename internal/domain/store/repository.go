package store

import "context"

type Repository interface {
	CreatePurchase(purchase Purchase, ctx context.Context) error
}
