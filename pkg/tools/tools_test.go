package tools

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"go-clickhouse-migrator/internal/domain/model"
)

func TestGetDynamicQuery(t *testing.T) {
	query, err := GetDynamicQuery("select {{ .TableName }}", &model.MigrationQueryParams{TableName: "test"})
	assert.NoError(t, err)
	assert.Equal(t, query, "select test")
}

type testConfig struct {
	Prefixed    string `envconfig:"TEST_ENV_VALUE"`
	NonPrefixed string `envconfig:"TEST_ENV_VALUE2"`
}

func TestProcessEnv(t *testing.T) {
	expected := &testConfig{
		Prefixed:    "1",
		NonPrefixed: "2",
	}
	
	t.Setenv("TEST_PREFIX_TEST_ENV_VALUE", expected.Prefixed)
	t.Setenv("TEST_PREFIX_TEST_ENV_VALUE2", expected.NonPrefixed)
	
	cfg := &testConfig{}
	
	assert.NoError(t, ProcessEnv("TEST_PREFIX", cfg))
	assert.Equal(t, expected, cfg)
}
