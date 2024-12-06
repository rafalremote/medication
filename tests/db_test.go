package tests

import (
	"os"
	"testing"

	"medication/internal/db"

	"github.com/stretchr/testify/assert"
)

func TestNewDatabaseConnection(t *testing.T) {
	os.Setenv("DB_DRIVER", "postgres")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_NAME", "medication_test")

	// Test with valid configuration
	conn, err := db.New()
	if err != nil {
		t.Skipf("Skipping test: database connection unavailable: %v", err)
		return
	}

	assert.NoError(t, err)
	assert.NotNil(t, conn)
	if conn != nil {
		defer conn.Close()
	}

	// Test with missing driver
	os.Setenv("DB_DRIVER", "")
	conn, err = db.New()
	assert.Error(t, err)
	assert.Nil(t, conn)

	// Test with invalid port
	os.Setenv("DB_DRIVER", "postgres")
	os.Setenv("DB_PORT", "invalid_port")
	conn, err = db.New()
	assert.Error(t, err)
	assert.Nil(t, conn)

	// Reset environment variables for other tests
	os.Unsetenv("DB_DRIVER")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("DB_NAME")
}
