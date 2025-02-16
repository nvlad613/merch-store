package http

import (
	"errors"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"merch-store/config"
	"merch-store/pkg/httputil"
	"net/http"
)

type Server struct {
}

func NewServer(
	logger *zap.Logger,
	conf config.Server,
) *Server {
	e := echo.New()
	e.Use(ZapLogger(logger))

	jwtMiddlewareFunc := echojwt.WithConfig(echojwt.Config{
		SigningKey:    []byte(conf.Auth.Key),
		SigningMethod: conf.Auth.Method,
		ErrorHandler: func(ctx echo.Context, err error) error {
			var details = "authorization error"
			switch {
			case errors.Is(err, echojwt.ErrJWTMissing):
				details = "jwt is missing"
			case errors.Is(err, echojwt.ErrJWTInvalid):
				details = "jwt is invalid"
			}

			return httputil.SendError(http.StatusUnauthorized, details, ctx)
		},
	})

	return nil
}

func (s *Server) Start() error {
	return nil
}
