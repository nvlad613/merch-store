package balance

import (
	"context"
	"merch-store/pkg/timeprovider"
)

type Service interface {
	MakeTransaction(fromUser string, toUser string, amount int, ctx context.Context) error
	MakeReport(username string, ctx context.Context) (*TransactionsReport, error)
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

func (s *ServiceImpl) MakeTransaction(fromUser string, toUser string, amount int, ctx context.Context) error {
	return s.repository.MakeCoinTransaction(fromUser, toUser, amount, s.timeProvider.Now(), ctx)
}

func (s *ServiceImpl) MakeReport(username string, ctx context.Context) (*TransactionsReport, error) {
	return s.repository.GetTransactionsReport(username, ctx)
}
