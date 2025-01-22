package data

import (
	"context"
	"github.com/charmbracelet/log"
	"github.com/imhinotori/open_currency/internal/configuration"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"os"
	"time"
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

	_, err = service.getConvertionRates()
	if err != nil {
		return nil, err
	}

	go service.startRateUpdater(context.Background(), 15*time.Minute)

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

func (s *Service) startRateUpdater(ctx context.Context, interval time.Duration) {
	log.Info("Starting rate updater")
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.updateConversionRates()
		case <-ctx.Done():
			log.Print("Rate updater stopped")
			return
		}
	}
}

func (s *Service) updateConversionRates() {
	_, err := s.getConvertionRates()
	if err != nil {
		log.Error("Error updating conversion rates", "error", err)
		return
	}

	log.Print("Rates updated")
}
