package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"merch-store/internal/domain"
	"merch-store/internal/domain/balance"
	"merch-store/internal/domain/store"
)

func (r *RepositoryImpl) MakePurchase(purchase store.Purchase, username string, ctx context.Context) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		var merch Merch
		err := tx.NewSelect().
			Model(&merch).
			Where("name = ?", purchase.ProductName).
			Scan(ctx)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return store.MerchItemNotFound
			}

			return err
		}

		var userId int
		_, err = tx.NewUpdate().
			Table("users").
			Set("coins = coins - ?", merch.Price).
			Where("name = ?", username).
			Returning("id").
			Exec(ctx, &userId)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return domain.UserNotFoundError
			}

			return balance.NotEnoughCoinsError
		}

		purchaseEntity := Purchase{
			UserId:   userId,
			MerchId:  merch.Id,
			Quantity: purchase.Quantity,
			Occurred: purchase.Timestamp,
		}

		_, err = tx.NewInsert().
			Model(&purchaseEntity).
			Exec(ctx)

		return err
	})
}

func (r *RepositoryImpl) GetUserPurchases(userId int, ctx context.Context) ([]store.Purchase, error) {
	var purchases []PurchasePreview

	err := r.db.NewSelect().
		TableExpr("purchases as p").
		ColumnExpr("p.quantity, p.occurred, m.name as merch_name").
		Where("p.user_id = ?", userId).
		Join("join merch as m").
		JoinOn("p.merch_id = m.id").
		Scan(ctx, &purchases)
	if err != nil {
		return nil, err
	}

	models := lo.Map(purchases, func(item PurchasePreview, _ int) store.Purchase {
		return store.Purchase{
			ProductName: item.MerchName,
			Quantity:    item.Quantity,
			Timestamp:   item.Occurred,
		}
	})

	return models, nil
}
