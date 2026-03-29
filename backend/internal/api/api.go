package api

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/2yuri/pack-calculator/pkg/config"
	"github.com/2yuri/pack-calculator/pkg/rest"
)

type Deps struct {
	Port   string
	Prefix string
	Routes []rest.Routes
}

type API struct {
	echo *echo.Echo
	port string
}

func New(deps Deps) *API {
	e := echo.New()
	e.HideBanner = true

	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: strings.Split(config.Instance().App.AllowOrigins, ","),
	}))
	e.Use(echomiddleware.RequestLoggerWithConfig(echomiddleware.RequestLoggerConfig{
		LogURI:     true,
		LogStatus:  true,
		LogMethod:  true,
		LogLatency: true,
		LogValuesFunc: func(c echo.Context, v echomiddleware.RequestLoggerValues) error {
			zap.L().Info("request",
				zap.String("method", v.Method),
				zap.String("uri", v.URI),
				zap.Int("status", v.Status),
				zap.Duration("latency", v.Latency),
			)
			return nil
		},
	}))

	group := e.Group(deps.Prefix)
	for _, routes := range deps.Routes {
		routes.Register(group)
	}

	return &API{echo: e, port: deps.Port}
}

func (a *API) Start() error {
	zap.L().Info("starting server", zap.String("port", a.port))
	return a.echo.Start(fmt.Sprintf(":%s", a.port))
}
