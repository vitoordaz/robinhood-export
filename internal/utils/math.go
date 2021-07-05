package utils

import (
	"strconv"
	"strings"
)

// IsZero returns true if a given numerical string represents zero.
func IsZero(n string) (bool, error) {
	v, err := strconv.ParseFloat(strings.TrimSpace(n), 64)
	if err != nil {
		return false, err
	}
	return v == 0.0, nil
}
