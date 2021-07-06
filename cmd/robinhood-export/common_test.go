package main

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vitoordaz/robinhood-export/internal/robinhood"
)

func TestLoadMarkets(t *testing.T) {
	t.Parallel()

	mockClient := &robinhood.MockClient{}
	marketByID := map[string]*robinhood.Market{"m1": {URL: "m1", Acronym: "m1"}, "m2": {URL: "m2", Acronym: "m2"}}
	mockClient.GetMarketFunc = func(ctx context.Context, id string) (*robinhood.Market, error) {
		if market, ok := marketByID[id]; ok {
			return market, nil
		}
		return nil, &robinhood.ResponseError{Detail: "not found"}
	}

	testCases := []struct {
		ids         []string
		errExpected bool
	}{
		{[]string{}, false},
		{[]string{"m1", "m2"}, false},
		{[]string{"m1", "m2", "m3"}, true},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("loadMarkets(%s)", strings.Join(tc.ids, ", ")), func(t *testing.T) {
			t.Parallel()
			markets, err := loadMarkets(context.Background(), mockClient, tc.ids)
			if tc.errExpected {
				require.Error(t, err)
				require.EqualError(t, err, "not found")
			} else {
				require.NoError(t, err)
				require.Len(t, markets, len(tc.ids))
				sort.Slice(markets, func(i, j int) bool { return markets[i].URL < markets[j].URL })
				marketIds := make([]string, 0, len(markets))
				for _, market := range markets {
					marketIds = append(marketIds, market.URL)
				}
				require.Equal(t, tc.ids, marketIds)
			}
		})
	}
}
