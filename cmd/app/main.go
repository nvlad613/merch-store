package main

import (
	"github.com/samber/lo"
	"merch-store/config"
)

func main() {
	// Init configuration
	conf := lo.Must(config.Load())

	// Init logger
	logger := lo.Must(conf.Logger.Build())
	defer logger.Sync()

	logger.Info(conf.Database.User.Login)
}
