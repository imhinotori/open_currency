package server

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/go-fuego/fuego"
	"github.com/imhinotori/open_currency/internal/configuration"
	"github.com/imhinotori/open_currency/internal/services"
	"go.uber.org/fx"
)

type Server struct {
	*fuego.Server
	services *services.Services
	cfg      *configuration.Configuration
}

func New(lc fx.Lifecycle, services *services.Services, cfg *configuration.Configuration) (*Server, error) {
	fuegoServer := fuego.NewServer()
	services.DataService.InitializeHandlers(fuegoServer)

	server := &Server{
		Server:   fuegoServer,
		services: services,
		cfg:      cfg,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.start(); err != nil {
					log.Fatal(err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})

	return server, nil
}

func (s *Server) start() error {
	if s.cfg.HTTP.SSL {
		return s.ListenAndServeTLS(s.cfg.HTTP.SSLCert, s.cfg.HTTP.SSLKey)
	}

	return s.ListenAndServe()
}
