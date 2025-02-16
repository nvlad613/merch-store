package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/pgdriver"
	"merch-store/internal/domain"
	"merch-store/internal/domain/balance"
	"time"
)

func (r *RepositoryImpl) MakeCoinTransaction(
	fromUser string,
	toUser string,
	amount int,
	timestamp time.Time,
	ctx context.Context,
) error {
	return r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		var fromUserId, toUserId int

		_, err := tx.NewUpdate().
			Table("users").
			Set("coins = coins - ?", amount).
			Where("name = ?", fromUser).
			Returning("id").
			Exec(ctx, &fromUserId)

		if err != nil {
			if isCheckError(err) {
				return balance.NotEnoughCoinsError
			}

			return err
		}

		_, err = tx.NewUpdate().
			Table("users").
			Set("coins = coins + ?", amount).
			Where("name = ?", toUser).
			Returning("id").
			Exec(ctx, &toUserId)
		if err != nil {
			return err
		}

		transaction := Transaction{
			SenderId:    fromUserId,
			RecipientId: toUserId,
			Amount:      amount,
			Occurred:    timestamp,
		}
		_, err = tx.NewInsert().
			Model(&transaction).
			Exec(ctx)

		return err
	})
}

func (r *RepositoryImpl) GetTransactionsReport(username string, ctx context.Context) (*balance.TransactionsReport, error) {
	var (
		user         User
		transactions []TransactionPreview
	)

	err := r.db.NewSelect().
		Model(&user).
		Where("name = ?", username).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.UserNotFoundError
		}
	}

	err = r.db.NewRaw(`
		with user_transactions as (
			select t.sender_id, t.recipient_id, t.amount, t.id from users u
			join transactions t on u.id in (t.recipient_id, t.sender_id)
			where u.id = ?
		)
		select ut.amount, u.name as sender_name, uu.name as recipient_name from user_transactions ut
		join users u on (ut.sender_id = u.id)
		join users uu on (ut.recipient_id = uu.id);`,
		user.Id,
	).Scan(ctx, &transactions)
	if err != nil {
		return nil, err
	}

	transactionModels := lo.Map(transactions, func(item TransactionPreview, _ int) balance.Transaction {
		return item.ToModel(username)
	})

	return &balance.TransactionsReport{
		User:         username,
		Coins:        user.Coins,
		Transactions: transactionModels,
	}, nil
}

func isCheckError(err error) bool {
	var pgErr *pgdriver.Error
	if errors.As(err, &pgErr) {
		return pgErr.Field('n') != ""
	}
	return false
}
