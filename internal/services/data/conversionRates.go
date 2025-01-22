package data

import "sync"

type ConversionRates struct {
	Rates map[string]float64
	Names map[string]string
	sync.RWMutex
}

var conversionRates = ConversionRates{
	Rates: make(map[string]float64),
	Names: make(map[string]string),
}
