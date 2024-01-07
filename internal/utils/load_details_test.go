package utils

import (
	"context"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadDetails(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type item struct {
		Field int
	}

	ids := make([]string, 0, 100)
	items := make([]*item, 0, 100)
	itemByID := make(map[string]*item, 100)
	for i := 0; i < 100; i++ {
		id := strconv.FormatInt(int64(i), 10)
		ids = append(ids, id)
		ii := &item{Field: i}
		itemByID[id] = ii
		items = append(items, ii)
	}

	result, err := LoadDetails(ctx, ids, func(ctx context.Context, id string) (*item, error) {
		return itemByID[id], nil
	})
	require.NoError(t, err)
	sort.Slice(result, func(i, j int) bool { return result[i].Field < result[j].Field })
	require.Equal(t, items, result)
}
