package balance

import (
	"context"
	"time"
)

type Repository interface {
	MakeCoinTransaction(fromUser string, toUser string, amount int, timestamp time.Time, ctx context.Context) error
	GetTransactionsReport(username string, ctx context.Context) (*TransactionsReport, error)
}
