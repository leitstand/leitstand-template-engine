package configen

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_featureIsEnable(t *testing.T) {

	tests := []struct {
		name       string
		constraint interface{}
		version    interface{}
		enabled    bool
		err        error
	}{
		{
			name:       "version matches feature toggle constraint",
			constraint: ">= 1.0.0",
			version:    "1.0.0",
			enabled:    true,
			err:        nil,
		},
		{
			name:       "version does not match feature toggle constraint",
			constraint: ">= 2.0.0",
			version:    "1.0.0",
			enabled:    false,
			err:        nil,
		},
		{
			name:       "nil constraint disables feature",
			constraint: nil,
			version:    "1.0.0",
			enabled:    false,
			err:        nil,
		},
		{
			name:       "empty constraint disables feature",
			constraint: "",
			version:    "1.0.0",
			enabled:    false,
			err:        nil,
		},
		{
			name:       "invalid constraint produces an error",
			constraint: "foo",
			version:    "1.0.0",
			enabled:    false,
			err:        fmt.Errorf("improper constraint: foo"),
		},
		{
			name:       "invalid version produces an error",
			constraint: ">=1.0.0",
			version:    "",
			enabled:    false,
			err:        fmt.Errorf("Invalid Semantic Version"),
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			enabled, err := featureIsEnabled(tt.constraint, tt.version)
			require.EqualValues(t, tt.enabled, enabled)
			if tt.err == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tt.err.Error())
			}
		})

	}

}
