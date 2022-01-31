package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestConfig_VarsFromFile(t *testing.T) {
	v := viper.GetInt("UNIT_TEST_VALUE")
	require.Equal(t, 10, v)
}
