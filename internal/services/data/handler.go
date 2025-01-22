package data

import (
	"github.com/charmbracelet/log"
	"github.com/go-fuego/fuego"
)

func (s *Service) InitializeHandlers(srv *fuego.Server) {
	h := &Handler{
		service: s,
	}

	fuego.Post(srv, "/convert", h.convertCurrency)
}

type Handler struct {
	service *Service
}

type CurrencyConvertionRequest struct {
	FromCurrency string  `json:"from_currency" validate:"required"`
	ToCurrency   string  `json:"to_currency" validate:"required"`
	Amount       float64 `json:"amount" validate:"required"`
}

type CurrencyConvertionResponse struct {
	Result float64 `json:"result"`
}

func (h *Handler) convertCurrency(c fuego.ContextWithBody[CurrencyConvertionRequest]) (*CurrencyConvertionResponse, error) {
	body, err := c.Body()
	if err != nil {
		log.Error("there was an error getting request body", "error", err)
		return nil, err
	}

	result, err := h.service.ConvertCurrency(body.FromCurrency, body.ToCurrency, body.Amount)
	if err != nil {
		log.Error("there was an error converting currency", "error", err)
		return nil, err
	}

	return &CurrencyConvertionResponse{Result: result}, err
}
