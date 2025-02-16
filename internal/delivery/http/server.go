package http

import (
	"context"
	"fmt"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"merch-store/config"
	authdelivery "merch-store/internal/delivery/http/auth"
	balancedelivery "merch-store/internal/delivery/http/balance"
	storedelivery "merch-store/internal/delivery/http/store"
	"merch-store/internal/domain/auth"
	"merch-store/internal/domain/balance"
	"merch-store/internal/domain/store"
	"merch-store/pkg/httputil"
	"strings"
)

type Server struct {
	inner *echo.Echo
	conf  config.Server
}

func NewServer(
	authService auth.Service,
	balanceService balance.Service,
	storeService store.Service,
	logger *zap.Logger,
	conf config.Server,
) *Server {
	e := echo.New()
	e.Use(ZapLogger(logger))

	jwtMiddlewareFunc := echojwt.WithConfig(echojwt.Config{
		SigningKey:    []byte(conf.Auth.Key),
		SigningMethod: strings.ToUpper(conf.Auth.Method),
		ErrorHandler:  httputil.JwtErrorHandler,
	})

	authRouter := authdelivery.NewRouter(logger, authService)
	e.POST("/api/auth", authRouter.PostAuthorizeUserHandler)

	balanceRouter := balancedelivery.NewRouter(logger, balanceService, storeService)
	e.GET("/api/info", balanceRouter.GetUserInfoHandler, jwtMiddlewareFunc)
	e.POST("/api/sendCoin", balanceRouter.PostSendCoinsHandler, jwtMiddlewareFunc)

	storeRouter := storedelivery.NewRouter(logger, storeService)
	e.GET("/api/buy/:item", storeRouter.GetBuyMerchHandler, jwtMiddlewareFunc)

	return &Server{
		inner: e,
		conf:  conf,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.conf.Hostname, s.conf.Port)
	return s.inner.Start(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.inner.Shutdown(ctx)
}
