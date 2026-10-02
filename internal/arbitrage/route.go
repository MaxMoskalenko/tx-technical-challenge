package arbitrage

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Route struct {
	currencies   []string
	exchangeRate decimal.Decimal
	isFinalised  bool
}

func NewRoute(startingCurrency string) Route {
	r := Route{
		currencies:   []string{startingCurrency},
		exchangeRate: decimal.NewFromInt(1),
	}

	return r
}

func (r Route) IsFinalised() bool             { return r.isFinalised }
func (r Route) Currencies() []string          { return r.currencies }
func (r Route) ExchangeRate() decimal.Decimal { return r.exchangeRate }
func (r Route) FormattedExchangeRate() string { return r.exchangeRate.StringFixed(exchangeRateFixed) }

func (r Route) StartsWith(currency string) bool {
	return len(r.currencies) > 0 && r.currencies[0] == currency
}

func (r Route) Contains(currency string) bool {
	for _, cur := range r.currencies {
		if cur == currency {
			return true
		}
	}

	return false
}

func (r Route) ExtendRoute(nextCurrency string, exchangeRateToLast decimal.Decimal) Route {
	currencies := append([]string{}, r.currencies...)
	currencies = append(currencies, nextCurrency)

	return Route{
		exchangeRate: r.exchangeRate.Mul(exchangeRateToLast),
		currencies:   currencies,
	}
}

func (r1 Route) IsSameRoute(r2 Route) bool {
	if len(r1.currencies) != len(r2.currencies) {
		return false
	}

	for offset := range r2.currencies {
		if r1.currencies[0] != r2.currencies[offset] {
			continue
		}

		isSameRoute := true
		for i := range r1.currencies {
			if r1.currencies[i] != r2.currencies[(offset+i)%len(r2.currencies)] {
				isSameRoute = false
				break
			}
		}

		if isSameRoute {
			return true
		}
	}

	return false
}

// returns array of markets in order how currencies should be exchanged
func (r Route) MarketsPath(markets []Market) []Market {
	if len(r.currencies) < 2 {
		return nil
	}

	currencies := append([]string{}, r.currencies...)
	if r.isFinalised {
		currencies = append(currencies, r.currencies[0])
	}

	marketsPath := make([]Market, 0, len(currencies)-1)
	for i := 0; i < len(currencies)-1; i++ {
		market, ok := findMarket(markets, currencies[i], currencies[i+1])
		if !ok {
			return nil
		}

		marketsPath = append(marketsPath, market)
	}

	return marketsPath
}

func (r Route) String() string {
	s := ""

	for _, c := range r.currencies {
		s += fmt.Sprintf("%s:", c)
	}

	if r.isFinalised {
		s += r.currencies[0]
	}

	return fmt.Sprintf("%s=%s", s, r.exchangeRate.String())
}

func (r *Route) Finalise(er decimal.Decimal) {
	r.exchangeRate = r.exchangeRate.Mul(er)
	r.isFinalised = true
}
