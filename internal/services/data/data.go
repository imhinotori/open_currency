package data

import (
	"fmt"
	"github.com/charmbracelet/log"
	"strconv"
	"strings"
)

func (s *Service) getConvertionRates() (*ConversionRates, error) {
	sheetRange := "Converter!A1:C100"

	response, err := s.Spreadsheets.Values.Get(s.cfg.Google.SheetId, sheetRange).Do()
	if err != nil {
		return nil, err
	}

	if len(response.Values) == 0 {
		return nil, fmt.Errorf("no data found")
	}

	conversionRates.Lock()
	defer conversionRates.Unlock()

	conversionRates.Rates = make(map[string]float64)
	conversionRates.Names = make(map[string]string)

	for i, row := range response.Values {
		// Header
		if i == 0 {
			continue
		}

		if len(row) < 3 {
			continue
		}

		currencyName, ok := row[0].(string)
		if !ok {
			log.Error("there was an error converting currency name to string", "row", i)
			continue
		}
		currencyCode, ok := row[1].(string)
		if !ok {
			log.Error("there was an error converting currency code to string", "row", i, "name", currencyName)
			continue
		}

		currencyConversionAsString, ok := row[2].(string)
		if !ok {
			log.Error("there was an error converting currency conversion rate to string", "row", i, "name", currencyName)
			continue
		}

		currencyCode = strings.TrimSpace(strings.ToUpper(currencyCode))
		currencyName = strings.TrimSpace(currencyName)
		currencyConversionAsString = strings.TrimSpace(currencyConversionAsString)

		if strings.ToUpper(currencyConversionAsString) == "#N/A" || currencyConversionAsString == "" {
			log.Warn("currency conversion rate is not a number", "row", i, "name", currencyName)
			continue
		}

		currencyConversion, err := strconv.ParseFloat(strings.ReplaceAll(currencyConversionAsString, ",", "."), 64)
		if err != nil {
			log.Error("there was an error converting currency conversion rate to float", "row", i, "name", currencyName)
			continue
		}

		conversionRates.Rates[currencyCode] = currencyConversion
		conversionRates.Names[currencyCode] = currencyName
	}

	log.Print("conversion rates loaded", "total currencies", len(conversionRates.Rates))

	return &conversionRates, nil
}

func (s *Service) ConvertCurrency(from string, to string, amount float64) (float64, error) {
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	if from == "" || to == "" {
		return 0, fmt.Errorf("invalid currency codes")
	}

	// Fucking donkey
	if from == to {
		return amount, nil
	}

	conversionRates.RLock()
	defer conversionRates.RUnlock()

	fromRate, exists := conversionRates.Rates[from]
	if !exists {
		return 0, fmt.Errorf("currency code (from) not found")
	}

	toRate, exists := conversionRates.Rates[to]
	if !exists {
		return 0, fmt.Errorf("currency code (to) not found")
	}

	amountInUSD := amount / fromRate
	convertedAmount := amountInUSD * toRate

	return convertedAmount, nil
}
