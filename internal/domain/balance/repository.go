package balance

import "context"

type Repository interface {
	MakeTransaction(transaction Transaction, ctx context.Context) error
	GetFinancialReport(username string, ctx context.Context) (*TransactionsReport, error)
}
