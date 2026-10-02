package arbitrage

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Market struct {
	base  string
	quote string

	exchangeRate decimal.Decimal
}

func NewMarket(base, quote string, exchangeRate decimal.Decimal) Market {
	return Market{
		base:         base,
		quote:        quote,
		exchangeRate: exchangeRate,
	}
}

func (m Market) Base() string                  { return m.base }
func (m Market) Quote() string                 { return m.quote }
func (m Market) ExchangeRate() decimal.Decimal { return m.exchangeRate }

func (m Market) HasCurrencies(c1, c2 string) bool {
	return (m.base == c1 && m.quote == c2) || (m.base == c2 && m.quote == c1)
}

func (m Market) Symbol() string {
	return fmt.Sprintf("%s/%s", m.base, m.quote)
}

func findMarket(markets []Market, c1, c2 string) (Market, bool) {
	for _, market := range markets {
		if market.HasCurrencies(c1, c2) {
			return market, true
		}
	}

	return Market{}, false
}
