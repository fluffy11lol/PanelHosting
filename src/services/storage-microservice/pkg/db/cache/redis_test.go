package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		cfg  RedisConfig
	}{
		{
			name: "valid config",
			cfg: RedisConfig{
				Host: "localhost",
				Port: "6379",
			},
		},
		{
			name: "custom port",
			cfg: RedisConfig{
				Host: "localhost",
				Port: "6380",
			},
		},
		{
			name: "custom host",
			cfg: RedisConfig{
				Host: "redis",
				Port: "6379",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := New(tt.cfg)
			assert.NotNil(t, client)
		})
	}
}
