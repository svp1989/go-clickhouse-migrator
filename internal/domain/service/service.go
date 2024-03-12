package service

import (
	"context"
	
	"go-clickhouse-migrator/internal/domain/model"
)

//go:generate mockery --name FileManager

type FileManager interface {
	SortedMigrationFilesData() ([]model.MigrationFileInfo, error)
	Read(data model.MigrationFileInfo) (string, error)
	Create(name string) (string, error)
	ReadHelp(file string) (string, error)
}

//go:generate mockery --name Command

type Command interface {
	Up(ctx context.Context, query string) error
	Init(ctx context.Context, params model.MigrationQueryParams) error
	CurrentVersion(ctx context.Context, params model.MigrationQueryParams) (model.MigrationInfo, error)
	Save(ctx context.Context, migration model.MigrationInfo, params model.MigrationQueryParams) error
	ExecutedMigration(ctx context.Context, params model.MigrationQueryParams) (map[string]model.MigrationInfo, error)
}

//go:generate mockery --name Migration

type Migration interface {
	GetQueryParams() model.MigrationQueryParams
}
