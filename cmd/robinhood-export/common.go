package main

import (
	"context"

	"github.com/vitoordaz/robinhood-export/internal/robinhood"
	"github.com/vitoordaz/robinhood-export/internal/utils"
)

func getInstrumentsMarketIds(instruments []*robinhood.Instrument) []string {
	return utils.GetIDs(instruments, func(instrument interface{}) string {
		return instrument.(*robinhood.Instrument).Market
	})
}

func loadInstruments(ctx context.Context, client robinhood.Client, ids []string) ([]*robinhood.Instrument, error) {
	instruments := make([]*robinhood.Instrument, 0, len(ids))
	err := utils.LoadDetails(ctx, ids, &instruments, func(ctx context.Context, id string) (interface{}, error) {
		return client.GetInstrument(ctx, id)
	})
	if err != nil {
		return nil, err
	}
	return instruments, nil
}

func loadMarkets(ctx context.Context, client robinhood.Client, ids []string) ([]*robinhood.Market, error) {
	markets := make([]*robinhood.Market, 0, len(ids))
	err := utils.LoadDetails(ctx, ids, &markets, func(ctx context.Context, id string) (interface{}, error) {
		return client.GetMarket(ctx, id)
	})
	if err != nil {
		return nil, err
	}
	return markets, nil
}

func getInstrumentByURL(instruments []*robinhood.Instrument) map[string]*robinhood.Instrument {
	instrumentByURL := make(map[string]*robinhood.Instrument, len(instruments))
	for _, instrument := range instruments {
		instrumentByURL[instrument.URL] = instrument
	}
	return instrumentByURL
}

func getMarketByURL(markets []*robinhood.Market) map[string]*robinhood.Market {
	marketByURL := make(map[string]*robinhood.Market, len(markets))
	for _, market := range markets {
		marketByURL[market.URL] = market
	}
	return marketByURL
}
