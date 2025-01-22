package services

import (
	"github.com/charmbracelet/log"
	"github.com/imhinotori/open_currency/internal/configuration"
	"github.com/imhinotori/open_currency/internal/services/data"
)

type Services struct {
	DataService *data.Service
}

func New(cfg *configuration.Configuration) (*Services, error) {
	dataService, err := data.New(cfg)
	if err != nil {
		log.Error("there was an error initializing data service", "error", err)
		return nil, err
	}

	services := &Services{
		DataService: dataService,
	}

	return services, nil
}
