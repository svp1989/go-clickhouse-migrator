package clickhouse

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}

	cfg := Config{
		Server:   "localhost",
		Port:     "9000",
		Database: "",
		User:     "admin",
		Password: "123",
	}
	conn, err := Connect(&cfg)
	assert.NoError(t, err)

	var res uint8

	row := conn.QueryRow(context.Background(), "select 9")
	err = row.Scan(&res)

	assert.NoError(t, err)
	assert.Equal(t, res, uint8(9))
}
