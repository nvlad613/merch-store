package main

import (
	"context"
	"github.com/samber/lo"
	"merch-store/config"
	"merch-store/internal/delivery/http"
	"merch-store/internal/domain/auth"
	"merch-store/internal/domain/balance"
	"merch-store/internal/domain/store"
	"merch-store/internal/infra"
	"merch-store/internal/infra/repository"
	"merch-store/pkg/timeprovider"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Init configuration
	conf := lo.Must(config.Load())

	// Init logger
	logger := lo.Must(conf.Logger.Build())
	defer logger.Sync()

	db := lo.Must(infra.InitDb(conf.Database))
	rep := repository.NewRepository(db)
	timeProvider := timeprovider.NewProvider(timeprovider.Moscow)

	authService := lo.Must(auth.NewService(rep, timeProvider, conf.Server.Auth))
	balanceService := balance.New(rep, timeProvider)
	storeService := store.New(rep, timeProvider)

	server := http.NewServer(
		authService,
		balanceService,
		storeService,
		logger,
		conf.Server,
	)

	go func() {
		if err := server.Start(); err != nil {
			logger.Fatal(err.Error())
		}
	}()

	WaitInterrupt()
	logger.Info("Gracefully shutting down...")

	timeout := time.Duration(conf.Server.ShutdownTimeoutSec) * time.Second
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	if err := server.Shutdown(ctx); err != nil {
		logger.Sugar().Errorw("failed to shutdown gracefully", "error", err)
	}

	logger.Info("Server shutdown")
}

func WaitInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
