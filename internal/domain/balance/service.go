package balance

import "context"

type Service interface {
	MakeTransaction(fromUser string, toUser string, amount int, ctx context.Context) error
	MakeReport(username string, ctx context.Context) (*TransactionsReport, error)
}
