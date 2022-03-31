package acceptance

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func RequireJSONEq(t *testing.T, v1, v2 interface{}) {
	json1, err := json.Marshal(v1)
	require.NoError(t, err)

	json2, err := json.Marshal(v2)
	require.NoError(t, err)

	require.Equal(t, string(json1), string(json2))
}
