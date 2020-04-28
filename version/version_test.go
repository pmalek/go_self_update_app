// +build !windows

package version

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetNumber(t *testing.T) {
	testcases := []struct {
		filename        string
		expectedVersion int
		expectedError   bool
	}{
		{
			filename:        `server_v12`,
			expectedVersion: 12,
		},
		{
			filename:        `server_v1`,
			expectedVersion: 1,
		},
		{
			filename:      `server_x`,
			expectedError: true,
		},
	}

	for i, tc := range testcases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			v, err := GetNumber(tc.filename)

			if tc.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedVersion, v)
			}
		})
	}
}
