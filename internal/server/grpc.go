package server

import (
	"context"
	"fmt"

	"github.com/MaxMoskalenko/tx-technical-challenge/internal/arbitrage"
	pb "github.com/MaxMoskalenko/tx-technical-challenge/internal/proto"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	pb.UnimplementedAPIServer
}

func NewGRPCServer() *GRPCServer {
	return &GRPCServer{}
}

func (s *GRPCServer) GetArbitrageOpportunities(ctx context.Context, req *pb.GetArbitrageOpportunitiesRequest) (*pb.GetArbitrageOpportunitiesResponse, error) {
	markets, err := marketsFromProto(req.GetMarkets())
	if err != nil {
		return nil, err
	}

	arbitrager := arbitrage.NewArbitrager(markets)
	routes := arbitrager.FindProfitableRoutes()

	opportunities := routesToProtoOpportunities(routes, markets)

	return &pb.GetArbitrageOpportunitiesResponse{
		Opportunities: opportunities,
	}, nil
}

func marketsFromProto(protoMarkets []*pb.Market) ([]arbitrage.Market, error) {
	markets := make([]arbitrage.Market, 0, len(protoMarkets))

	for i, market := range protoMarkets {
		if market == nil {
			continue
		}
		if market.GetBaseCurrency() == "" {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("markets[%d].base_currency is required", i))
		}
		if market.GetQuoteCurrency() == "" {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("markets[%d].quote_currency is required", i))
		}
		if market.GetBaseCurrency() == market.GetQuoteCurrency() {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("markets[%d].base_currency = markets[%d].quote_currency", i, i))
		}

		exchangeRate, err := decimal.NewFromString(market.GetExchangeRate())
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("markets[%d].exchange_rate must be a decimal", i))
		}
		if !exchangeRate.IsPositive() {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("markets[%d].exchange_rate must be positive", i))
		}

		markets = append(markets, arbitrage.NewMarket(
			market.GetBaseCurrency(),
			market.GetQuoteCurrency(),
			exchangeRate,
		))
	}

	return markets, nil
}

func routesToProtoOpportunities(routes []arbitrage.Route, markets []arbitrage.Market) []*pb.ArbitrageOpportunity {
	opportunities := make([]*pb.ArbitrageOpportunity, 0, len(routes))
	for _, route := range routes {
		marketsPath := []string{}
		for _, m := range route.MarketsPath(markets) {
			marketsPath = append(marketsPath, m.Symbol())
		}

		opportunities = append(opportunities, &pb.ArbitrageOpportunity{
			Path:          marketsPath,
			Profitability: route.FormattedExchangeRate(),
		})
	}

	return opportunities
}
