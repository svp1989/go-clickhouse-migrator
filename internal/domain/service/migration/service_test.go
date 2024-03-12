package migration

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"go-clickhouse-migrator/internal/domain/model"
	"go-clickhouse-migrator/pkg/migrator"
)

func TestService_GetQueryParams(t *testing.T) {
	assert.Equal(t, New(&migrator.Config{Table: "migration_version"}).GetQueryParams(), model.MigrationQueryParams{
		TableName: "migration_version",
	})
}
