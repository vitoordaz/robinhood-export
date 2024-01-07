package utils

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadList(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type item struct {
		Field int
	}

	var items []*item
	for i := 0; i < 100; i++ {
		items = append(items, &item{Field: i})
	}

	result, err := LoadList(ctx, func(c context.Context, cursor string) ([]*item, string, error) {
		var err error
		var idx int64
		if cursor != "" {
			idx, err = strconv.ParseInt(cursor, 10, 64)
			if err != nil {
				return nil, "", err
			}
		}
		if idx+10 < int64(len(items)) {
			cursor = strconv.FormatInt(idx+10, 10)
		} else {
			cursor = ""
		}
		return items[idx : idx+10], cursor, nil
	})
	require.NoError(t, err)
	require.Equal(t, items, result)
}
