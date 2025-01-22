package main

import (
	"github.com/charmbracelet/log"
	"github.com/imhinotori/open_currency/internal/configuration"
	"github.com/imhinotori/open_currency/internal/server"
	"github.com/imhinotori/open_currency/internal/services"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			configuration.Load,
			services.New,
			server.New,
		),
		fx.Invoke(func(_ *server.Server) {
			log.Info("Server started")
		}),
	)

	app.Run()
}
