package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	configDir := "./configs"
	err := os.MkdirAll(configDir, 0755)
	assert.NoError(t, err)

	tempConfigContent := `
GRPC_SERVER_PORT=9091
REST_SERVER_PORT=8081
POSTGRES_HOST=test-host
POSTGRES_PORT=5432
POSTGRES_USER=test-user
POSTGRES_PASSWORD=test-password
POSTGRES_DB=test-db
REDIS_HOST=localhost
REDIS_PORT=6379
`
	err = os.WriteFile(configDir+"/local.env", []byte(tempConfigContent), 0644)
	assert.NoError(t, err)

	t.Cleanup(func() {
		os.RemoveAll(configDir)
	})

	t.Run("successful config load", func(t *testing.T) {
		cfg := New()
		assert.NotNil(t, cfg)
		assert.Equal(t, 9091, cfg.GRPCServerPort)
		assert.Equal(t, 8081, cfg.RestServerPort)
		assert.Equal(t, "test-host", cfg.Config.Host)
		assert.Equal(t, "5432", cfg.Config.Port)
		assert.Equal(t, "test-user", cfg.Config.UserName)
		assert.Equal(t, "test-password", cfg.Config.Password)
		assert.Equal(t, "test-db", cfg.Config.DbName)
		assert.Equal(t, "localhost", cfg.RedisConfig.Host)
		assert.Equal(t, "6379", cfg.RedisConfig.Port)
	})
}

func TestNew_NoConfigFile(t *testing.T) {
	cfg := New()
	assert.Nil(t, cfg)
}
