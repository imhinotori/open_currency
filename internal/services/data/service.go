package data

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/imhinotori/open_currency/internal/configuration"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"os"
)

type Service struct {
	*sheets.Service
	cfg *configuration.Configuration
}

func New(cfg *configuration.Configuration) (*Service, error) {
	googleService, err := newSheetsService(cfg)
	if err != nil {
		return nil, err
	}

	service := &Service{
		Service: googleService,
		cfg:     cfg,
	}

	rates, err := service.getConvertionRates()
	if err != nil {
		return nil, err
	}

	log.Info("Rates loaded", "rates", rates)

	return service, nil
}

func newSheetsService(cfg *configuration.Configuration) (*sheets.Service, error) {
	ctx := context.Background()

	f, err := os.ReadFile(cfg.Google.CredentialsFile)
	if err != nil {
		return nil, err
	}

	credentials, err := google.CredentialsFromJSON(ctx, f, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}

	srv, err := sheets.NewService(ctx, option.WithCredentials(credentials))
	if err != nil {
		return nil, err
	}

	return srv, nil
}
