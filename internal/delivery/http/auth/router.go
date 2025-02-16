package auth

import (
	"go.uber.org/zap"
	"merch-store/internal/domain/auth"
)

type Router struct {
	logger      *zap.SugaredLogger
	authService auth.Service
}

func NewRouter(
	logger *zap.Logger,
	authService auth.Service,
) *Router {
	return &Router{
		logger:      logger.Sugar(),
		authService: authService,
	}
}
