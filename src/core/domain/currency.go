package domain

import (
	"errors"
	"strings"
)

type Currency string

type CurrencyPair struct {
	From Currency
	To   Currency
}

var CurrencyLookupMap = map[Currency]string{
	"USD": "USD",
	"EUR": "EUR",
}

var ConversionRateMap = map[CurrencyPair]float64{
	{"USD", "EUR"}: 0.85,
	{"EUR", "USD"}: 1.15,
}

func NewCurrencyPair(from Currency , to Currency) *CurrencyPair {
	return &CurrencyPair{
		From: from,
		To: to,
	}
}

func (c CurrencyPair) Calculate(amount float64) float64 {
	if c.From == c.To {
		return amount
	}
	return amount * ConversionRateMap[c]
}

func (c CurrencyPair) String() string {
	return CurrencyLookupMap[c.From] + "-" + CurrencyLookupMap[c.To]
}

func CurrencyPairParse(pair string) (CurrencyPair, error) {
	pairSplit := strings.Split(pair, "-")

	if len(pairSplit) != 2 {
		return CurrencyPair{}, errors.New("Bad formatting use: ###-@@@")
	}

	return CurrencyPair{
		From: Currency(pairSplit[0]),
		To:   Currency(pairSplit[1]),
	}, nil
	
}
