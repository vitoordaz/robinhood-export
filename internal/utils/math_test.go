package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsZero(t *testing.T) {
	cases := []struct {
		value       string
		isZero      bool
		errExpected bool
	}{
		{"0", true, false},
		{"0.0", true, false},
		{"1", false, false},
		{"1.0", false, false},
		{"", false, true},
	}

	for _, tc := range cases {
		isZero, err := IsZero(tc.value)
		if tc.errExpected {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
		if tc.isZero {
			require.True(t, isZero)
		} else {
			require.False(t, isZero)
		}
	}
}
