package main

import (
	"context"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitoordaz/robinhood-export/internal/robinhood"
)

func TestLoadMarkets(t *testing.T) {
	ctx := context.Background()
	mockClient := &robinhood.MockClient{}
	markets, err := loadMarkets(ctx, mockClient, []string{})
	require.NoError(t, err)
	require.Empty(t, markets)

	marketByID := map[string]*robinhood.Market{"m1": {URL: "m1", Acronym: "m1"}, "m2": {URL: "m2", Acronym: "m2"}}
	mockClient.GetMarketFunc = func(ctx context.Context, id string) (*robinhood.Market, error) {
		if market, ok := marketByID[id]; ok {
			return market, nil
		}
		return nil, &robinhood.ResponseError{Detail: "not found"}
	}

	markets, err = loadMarkets(ctx, mockClient, []string{"m1", "m2"})
	require.NoError(t, err)
	require.Len(t, markets, 2)
	sort.Slice(markets, func(i, j int) bool { return markets[i].URL < markets[j].URL })
	require.Equal(t, "m1", markets[0].URL)
	require.Equal(t, "m2", markets[1].URL)

	markets, err = loadMarkets(ctx, mockClient, []string{"m1", "m2", "m3"})
	require.EqualError(t, err, "not found")
	require.Nil(t, markets)
}
