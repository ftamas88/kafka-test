package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Default(t *testing.T) {
	conf, err := New()
	require.NoError(t, err)
	assert.Equal(t, conf.HTTPPort, 3000)
}

func Test_Configured(t *testing.T) {
	os.Setenv("HTTP_PORT", "8080")

	conf, err := New()
	require.NoError(t, err)
	assert.Equal(t, conf.HTTPPort, 8080)
}
