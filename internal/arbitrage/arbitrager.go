package arbitrage

import (
	"github.com/shopspring/decimal"
)

// could be config
const exchangeRateFixed = 6

type Arbitrager struct {
	matrix     ExchangeMatrix
	currencies map[string]struct{}
}

func NewArbitrager(markets []Market) Arbitrager {
	currencies := map[string]struct{}{}

	for _, m := range markets {
		currencies[m.Base()] = struct{}{}
		currencies[m.Quote()] = struct{}{}
	}

	matrix := NewExchangeMatrix(markets)
	return Arbitrager{
		matrix:     matrix,
		currencies: currencies,
	}
}

func (a *Arbitrager) FindAllRoutes() []Route {
	return a.findRoutes()
}

func (a *Arbitrager) FindProfitableRoutes() []Route {
	routes := a.findRoutes()
	profitable := []Route{}

	// small number to filter out rounding errors, could be replaced by "desired profitability rate"
	epsilon := decimal.RequireFromString("1.000001")

	for _, p := range routes {
		if p.exchangeRate.Cmp(epsilon) > 0 {
			profitable = append(profitable, p)
		}
	}

	return profitable
}

func (a *Arbitrager) findRoutes() []Route {
	ps := []Route{}

	// we need to iterate through all of the currencies to find all possible routes
	for c := range a.currencies {
		routesForCurrency := a.findRouteForCurrency(c)

		for _, pfc := range routesForCurrency {
			if !pfc.IsFinalised() {
				continue
			}

			// avoid processing the same cycled routes, e.g. A->B->C->A and C->A->B->C
			isSameRoute := false
			for _, p := range ps {
				if pfc.IsSameRoute(p) {
					isSameRoute = true
					break
				}
			}

			if !isSameRoute {
				ps = append(ps, pfc)
			}
		}
	}

	return ps
}

func (a *Arbitrager) findRouteForCurrency(currency string) []Route {
	return a.extendRoutes(currency, NewRoute(currency))
}

func (a *Arbitrager) extendRoutes(currency string, route Route) []Route {
	routes := []Route{}

	for c, er := range a.matrix.ExchangesForCurrency(currency) {
		if route.StartsWith(c) {
			p := route
			p.Finalise(er)
			routes = append(routes, p)
			continue
		}

		// route will contain loop
		if route.Contains(c) {
			continue
		}

		routes = append(routes, a.extendRoutes(c, route.ExtendRoute(c, er))...)
	}

	return routes
}
