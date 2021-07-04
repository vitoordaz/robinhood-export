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

	var deduplicated []string
	var duplicates []*item
	for i := 0; i < 100; i++ {
		id := strconv.FormatInt(int64(i), 10)
		deduplicated = append(deduplicated, id)
		for j := 0; j < 2; j++ {
			duplicates = append(duplicates, &item{id})
		}
	}

	result := GetIDs(duplicates, func(i interface{}) string { return i.(*item).ID })
	require.Equal(t, deduplicated, result)
}
