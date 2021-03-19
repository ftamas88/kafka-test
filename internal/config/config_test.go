package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Default(t *testing.T) {
	conf, err := New()

	var es []string

	require.NoError(t, err)
	assert.Equal(t, conf.HTTPPort, 3000)
	assert.Equal(t, conf.Topic, "test")
	assert.Equal(t, conf.SchemaPath, "pkg/avro/transaction_with_amount.avsc")
	assert.Equal(t, conf.CloudSchemaRegistryServers, es)
	assert.Equal(t, conf.SchemaRegistryServers, []string{"http://localhost:8081"})
	assert.Equal(t, conf.CloudServer, "")
	assert.Equal(t, conf.Servers, []string{"http://localhost:9092"})
}

func Test_Configured(t *testing.T) {
	os.Setenv("HTTP_PORT", "8080")
	os.Setenv("KAFKA_TOPIC", "test-2")

	conf, err := New()
	require.NoError(t, err)
	assert.Equal(t, conf.HTTPPort, 8080)
	assert.Equal(t, conf.Topic, "test-2")
}
