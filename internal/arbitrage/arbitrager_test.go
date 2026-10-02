package arbitrage

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

const testExchangeRateFixed = 4

func TestArbitrager_FindRoutes(t *testing.T) {
	tests := []struct {
		name             string
		markets          []Market
		startingCurrency string
		expectedRoutes   []Route
	}{
		{
			name: "test_1",
			markets: []Market{
				NewMarket("EUR", "USD", mustNewFromString("1.1586")),
				NewMarket("EUR", "GBP", mustNewFromString("0.6849")),
				NewMarket("USD", "GBP", mustNewFromString("0.5904")),
			},
			expectedRoutes: []Route{
				newMockRoute([]string{"USD", "GBP"}, mustNewFromString("1")),
				newMockRoute([]string{"USD", "EUR"}, mustNewFromString("1")),
				newMockRoute([]string{"EUR", "GBP"}, mustNewFromString("1")),
				newMockRoute([]string{"USD", "EUR", "GBP"}, mustNewFromString("1.0013")),
				newMockRoute([]string{"USD", "GBP", "EUR"}, mustNewFromString("0.9987")),
			},
		},
		{
			name: "test_2",
			markets: []Market{
				NewMarket("EUR", "USD", mustNewFromString("1.1")),
				NewMarket("EUR", "GBP", mustNewFromString("0.8")),
				NewMarket("USD", "GBP", mustNewFromString("0.7")),
				NewMarket("USD", "PLN", mustNewFromString("4")),
				NewMarket("PLN", "UAH", mustNewFromString("12")),
				NewMarket("UAH", "GBP", mustNewFromString("0.02")),
			},
			expectedRoutes: []Route{
				newMockRoute([]string{"GBP", "EUR"}, mustNewFromString("1")),
				newMockRoute([]string{"GBP", "UAH"}, mustNewFromString("1")),
				newMockRoute([]string{"PLN", "USD"}, mustNewFromString("1")),
				newMockRoute([]string{"PLN", "UAH"}, mustNewFromString("1")),
				newMockRoute([]string{"EUR", "USD"}, mustNewFromString("1")),
				newMockRoute([]string{"GBP", "USD"}, mustNewFromString("1")),
				newMockRoute([]string{"GBP", "EUR", "USD"}, mustNewFromString("0.9625")),
				newMockRoute([]string{"GBP", "USD", "EUR"}, mustNewFromString("1.039")),
				newMockRoute([]string{"GBP", "USD", "PLN", "UAH"}, mustNewFromString("1.3714")),
				newMockRoute([]string{"GBP", "UAH", "PLN", "USD"}, mustNewFromString("0.7292")),
				newMockRoute([]string{"GBP", "EUR", "USD", "PLN", "UAH"}, mustNewFromString("1.32")),
				newMockRoute([]string{"GBP", "UAH", "PLN", "USD", "EUR"}, mustNewFromString("0.7576")),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			arb := NewArbitrager(test.markets)

			routes := arb.FindAllRoutes()

			require.Len(t, routes, len(test.expectedRoutes), routes)

			for _, er := range test.expectedRoutes {
				assertRoute(t, er, routes)
			}
		})
	}
}

func TestArbitrager_FindProfitableRoutes(t *testing.T) {
	tests := []struct {
		name             string
		markets          []Market
		startingCurrency string
		expectedRoutes   []Route
	}{
		{
			name: "test_1",
			markets: []Market{
				NewMarket("EUR", "USD", mustNewFromString("1.1586")),
				NewMarket("EUR", "GBP", mustNewFromString("0.6849")),
				NewMarket("USD", "GBP", mustNewFromString("0.5904")),
			},
			expectedRoutes: []Route{
				newMockRoute([]string{"USD", "EUR", "GBP"}, mustNewFromString("1.0013")),
			},
		},
		{
			name: "test_2",
			markets: []Market{
				NewMarket("EUR", "USD", mustNewFromString("1.1")),
				NewMarket("EUR", "GBP", mustNewFromString("0.8")),
				NewMarket("USD", "GBP", mustNewFromString("0.7")),
				NewMarket("USD", "PLN", mustNewFromString("4")),
				NewMarket("PLN", "UAH", mustNewFromString("12")),
				NewMarket("UAH", "GBP", mustNewFromString("0.02")),
			},
			expectedRoutes: []Route{
				newMockRoute([]string{"GBP", "USD", "EUR"}, mustNewFromString("1.039")),
				newMockRoute([]string{"GBP", "USD", "PLN", "UAH"}, mustNewFromString("1.3714")),
				newMockRoute([]string{"GBP", "EUR", "USD", "PLN", "UAH"}, mustNewFromString("1.32")),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			arb := NewArbitrager(test.markets)

			routes := arb.FindProfitableRoutes()

			require.Len(t, routes, len(test.expectedRoutes), routes)

			for _, er := range test.expectedRoutes {
				assertRoute(t, er, routes)
			}
		})
	}
}

func TestRoute_MarketsPath(t *testing.T) {
	markets := []Market{
		NewMarket("EUR", "USD", mustNewFromString("1.1586")),
		NewMarket("EUR", "GBP", mustNewFromString("0.6849")),
		NewMarket("USD", "GBP", mustNewFromString("0.5904")),
	}

	route := newMockRoute([]string{"USD", "EUR", "GBP"}, mustNewFromString("1.0013"))

	actual := route.MarketsPath(markets)

	require.Equal(t, []Market{
		markets[0],
		markets[1],
		markets[2],
	}, actual)
}

func mustNewFromString(s string) decimal.Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		panic(err)
	}
	return d
}

func newMockRoute(currencies []string, exchangeRate decimal.Decimal) Route {
	return Route{
		currencies:   currencies,
		exchangeRate: exchangeRate,
		isFinalised:  true,
	}
}

func assertRoute(t *testing.T, expected Route, actualRoutes []Route) {
	t.Helper()

	for _, actual := range actualRoutes {
		if !expected.IsSameRoute(actual) {
			continue
		}

		require.Equal(t, len(expected.currencies), len(actual.currencies), expected.String(), actual.String())
		require.Equal(t, expected.exchangeRate.StringFixed(testExchangeRateFixed), actual.exchangeRate.StringFixed(testExchangeRateFixed), expected.String(), actual.String())
		require.True(t, actual.isFinalised, expected.String())
		return
	}

	require.Failf(t, "route not found", "expected route %s in %v", expected.String(), actualRoutes)
}
