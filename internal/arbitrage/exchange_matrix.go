package arbitrage

import "github.com/shopspring/decimal"

type ExchangeMatrix struct {
	rates map[string]map[string]decimal.Decimal // matrix (graph) of exchange rates and routes e.g. USD->EUR->1.2
}

func NewExchangeMatrix(markets []Market) ExchangeMatrix {
	rates := map[string]map[string]decimal.Decimal{}

	for _, m := range markets {
		if _, ok := rates[m.Base()]; !ok {
			rates[m.Base()] = map[string]decimal.Decimal{}
		}

		if _, ok := rates[m.Quote()]; !ok {
			rates[m.Quote()] = map[string]decimal.Decimal{}
		}

		// for simplicity, we are saving A->B and B->A routes
		rates[m.Base()][m.Quote()] = m.ExchangeRate()

		// this type operates under assumption, that exchange rate was validated before and is not 0
		rates[m.Quote()][m.Base()] = decimal.NewFromInt(1).Div(m.ExchangeRate())
	}

	return ExchangeMatrix{
		rates: rates,
	}
}

func (em *ExchangeMatrix) ExchangesForCurrency(currency string) map[string]decimal.Decimal {
	return em.rates[currency]
}
