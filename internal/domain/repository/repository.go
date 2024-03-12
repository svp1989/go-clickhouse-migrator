package repository

import (
	"context"
	
	"go-clickhouse-migrator/internal/domain/model"
)

type Command interface {
	Up(ctx context.Context, query string) error
	Init(ctx context.Context, params model.MigrationQueryParams) error
	GetMigrationInfo(ctx context.Context, params model.MigrationQueryParams) (model.MigrationInfo, error)
	GetMigrationInfoMap(ctx context.Context, params model.MigrationQueryParams) (map[string]model.MigrationInfo, error)
	Save(ctx context.Context, migration model.MigrationInfo, params model.MigrationQueryParams) error
}
