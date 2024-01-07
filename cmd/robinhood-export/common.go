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
	return utils.LoadDetails[robinhood.Instrument](ctx, ids, client.GetInstrument)
}

func loadMarkets(ctx context.Context, client robinhood.Client, ids []string) ([]*robinhood.Market, error) {
	return utils.LoadDetails[robinhood.Market](ctx, ids, client.GetMarket)
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
