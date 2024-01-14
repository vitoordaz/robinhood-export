package main

import (
	"context"

	"github.com/vitoordaz/robinhood-export/internal/robinhood"
	"github.com/vitoordaz/robinhood-export/internal/utils"
)

func getInstrumentsMarketIds(instruments []*robinhood.Instrument) []string {
	return utils.GetIDs(instruments, func(instrument *robinhood.Instrument) string {
		return instrument.Market
	})
}

func loadAccounts(
	ctx context.Context,
	client robinhood.Client,
	auth *robinhood.ResponseToken,
	ids []string,
) ([]*robinhood.Account, error) {
	return utils.LoadDetails(ctx, ids, func(ctx context.Context, id string) (*robinhood.Account, error) {
		return client.GetAccount(ctx, auth, id)
	})
}

func loadInstruments(ctx context.Context, client robinhood.Client, ids []string) ([]*robinhood.Instrument, error) {
	return utils.LoadDetails(ctx, ids, client.GetInstrument)
}

func loadMarkets(ctx context.Context, client robinhood.Client, ids []string) ([]*robinhood.Market, error) {
	return utils.LoadDetails(ctx, ids, client.GetMarket)
}

func getAccountByURL(accounts []*robinhood.Account) map[string]*robinhood.Account {
	accountByURL := make(map[string]*robinhood.Account, len(accounts))
	for _, account := range accounts {
		accountByURL[account.URL] = account
	}
	return accountByURL
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
