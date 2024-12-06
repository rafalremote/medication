package tests

import (
	"os"
	"testing"

	"medication/config"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("PORT", "3000")
	os.Setenv("DB_DRIVER", "postgres")
	os.Setenv("DB_HOST", "test-db")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "medication_test")
	os.Setenv("JWT_SECRET", "your_jwt_secret")

	// Load configuration
	cfg, err := config.LoadConfig()
	assert.NoError(t, err)

	// Verify the loaded configuration
	assert.Equal(t, "3000", cfg.Port)
	assert.Equal(t, "postgres", cfg.DBDriver)
	assert.Equal(t, "test-db", cfg.DBHost)
	assert.Equal(t, 5433, cfg.DBPort)
	assert.Equal(t, "postgres", cfg.DBUser)
	assert.Equal(t, "postgres", cfg.DBPassword)
	assert.Equal(t, "medication_test", cfg.DBName)
	assert.Equal(t, "your_jwt_secret", cfg.JWTSecret)
}

func TestGetEnvDefault(t *testing.T) {
	// Ensure the environment variable is not set
	os.Unsetenv("SOME_ENV_KEY")

	// Test getEnv with a default value
	value := config.GetEnv("SOME_ENV_KEY", "defaultvalue")
	assert.Equal(t, "defaultvalue", value)
}

func TestGetEnvSetValue(t *testing.T) {
	// Set an environment variable
	os.Setenv("SOME_ENV_KEY", "setvalue")

	// Test getEnv with an existing value
	value := config.GetEnv("SOME_ENV_KEY", "defaultvalue")
	assert.Equal(t, "setvalue", value)

	// Cleanup
	os.Unsetenv("SOME_ENV_KEY")
}
