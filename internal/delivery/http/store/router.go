package store

import (
	"go.uber.org/zap"
	"merch-store/internal/domain/store"
)

type Router struct {
	logger       *zap.SugaredLogger
	storeService store.Service
}

func NewRouter(
	logger *zap.Logger,
	storeService store.Service,
) *Router {
	return &Router{
		logger:       logger.Sugar(),
		storeService: storeService,
	}
}
