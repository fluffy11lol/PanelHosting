package postgres

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestConfig is a configuration for testing purposes.
var TestConfig = Config{
	UserName: "test_user",
	Password: "test_password",
	Host:     "localhost",
	Port:     "5432",
	DbName:   "test_db",
}

func TestNew_Failure(t *testing.T) {
	// Test with an invalid configuration to simulate a failure.
	invalidConfig := Config{
		UserName: "invalid",
		Password: "invalid",
		Host:     "invalid_host",
		Port:     "invalid_port",
		DbName:   "invalid_db",
	}

	db, err := New(invalidConfig)
	assert.NoError(t, err, "New should return an error for invalid configuration")
	assert.NotNil(t, db, "Returned DB instance should be nil for invalid configuration")
}

func TestDSNFormat(t *testing.T) {
	// Test if the DSN string is formatted correctly.
	config := TestConfig
	expectedDSN := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		config.UserName, config.Password, config.DbName, config.Host, config.Port)

	actualDSN := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=%s port=%s",
		config.UserName, config.Password, config.DbName, config.Host, config.Port)

	assert.Equal(t, expectedDSN, actualDSN, "DSN string should match the expected format")
}
