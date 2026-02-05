package utils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetIDs(t *testing.T) {
	type item struct {
		ID string
	}

	deduplicated := make([]string, 0, 100)
	duplicates := make([]*item, 0, 200)
	for i := range 100 {
		id := strconv.FormatInt(int64(i), 10)
		deduplicated = append(deduplicated, id)
		for range 2 {
			duplicates = append(duplicates, &item{id})
		}
	}

	result := GetIDs(duplicates, func(i *item) string { return i.ID })
	require.Equal(t, deduplicated, result)
}
