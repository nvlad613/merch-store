package store

import (
	"context"
	"maps"
	"merch-store/pkg/timeprovider"
	"slices"
)

type Service interface {
	MakePurchase(username string, productName string, quantity int, ctx context.Context) error
	GetUserInventory(userId int, ctx context.Context) ([]InventoryItem, error)
}

type ServiceImpl struct {
	repository   Repository
	timeProvider timeprovider.TimeProvider
}

func New(
	repository Repository,
	timeProvider timeprovider.TimeProvider,
) *ServiceImpl {
	return &ServiceImpl{
		repository:   repository,
		timeProvider: timeProvider,
	}
}

func (s ServiceImpl) MakePurchase(username string, productName string, quantity int, ctx context.Context) error {
	purchase := Purchase{
		ProductName: productName,
		Timestamp:   s.timeProvider.Now(),
		Quantity:    quantity,
	}

	return s.repository.MakePurchase(purchase, username, ctx)
}

func (s ServiceImpl) GetUserInventory(userId int, ctx context.Context) ([]InventoryItem, error) {
	purchases, err := s.repository.GetUserPurchases(userId, ctx)
	if err != nil {
		return nil, err
	}

	inventory := make(map[string]InventoryItem)
	for _, p := range purchases {
		if item, found := inventory[p.ProductName]; found {
			item.Quantity += p.Quantity
			inventory[p.ProductName] = item
		} else {
			inventory[p.ProductName] = InventoryItem{
				ProductName: p.ProductName,
				Quantity:    p.Quantity,
			}
		}
	}

	return slices.Collect(maps.Values(inventory)), nil
}
