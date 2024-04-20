package types

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

func (c CurrencyPair) Calculate(amount float64) float64 {
	if c.From == c.To {
		return amount
	}
	return amount * ConversionRateMap[c]
}