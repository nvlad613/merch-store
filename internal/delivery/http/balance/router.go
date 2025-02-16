package balance

import (
	"go.uber.org/zap"
	"merch-store/internal/domain/balance"
	"merch-store/internal/domain/store"
)

type Router struct {
	logger         *zap.SugaredLogger
	balanceService balance.Service
	storeService   store.Service
}

func NewRouter(
	logger *zap.Logger,
	balanceService balance.Service,
	storeService store.Service,
) *Router {
	return &Router{
		logger:         logger.Sugar(),
		balanceService: balanceService,
		storeService:   storeService,
	}
}
