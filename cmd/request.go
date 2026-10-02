package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/MaxMoskalenko/tx-technical-challenge/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func request(args []string) error {
	addr := "localhost:50051"
	if len(args) > 0 {
		addr = args[0]
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	resp, err := proto.NewAPIClient(conn).GetArbitrageOpportunities(ctx, sampleRequest())
	if err != nil {
		return err
	}

	if len(resp.GetOpportunities()) == 0 {
		fmt.Println("No arbitrage opportunities found")
		return nil
	}

	for _, opportunity := range resp.GetOpportunities() {
		fmt.Printf("%s -> %s\n", strings.Join(opportunity.GetPath(), " -> "), opportunity.GetProfitability())
	}

	return nil
}

func sampleRequest() *proto.GetArbitrageOpportunitiesRequest {
	return &proto.GetArbitrageOpportunitiesRequest{
		Markets: []*proto.Market{
			{
				BaseCurrency:  "EUR",
				QuoteCurrency: "USD",
				ExchangeRate:  "1.1586",
			},
			{
				BaseCurrency:  "EUR",
				QuoteCurrency: "GBP",
				ExchangeRate:  "0.6849",
			},
			{
				BaseCurrency:  "USD",
				QuoteCurrency: "GBP",
				ExchangeRate:  "0.5904",
			},
		},
	}
}
