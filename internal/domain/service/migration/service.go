package migration

import (
	"go-clickhouse-migrator/internal/domain/model"
	"go-clickhouse-migrator/pkg/migrator"
)

type Service struct {
	tableName string
}

func New(cfg *migrator.Config) *Service {
	return &Service{
		tableName: cfg.Table,
	}
}

func (m *Service) GetQueryParams() model.MigrationQueryParams {
	return model.MigrationQueryParams{
		TableName: m.tableName,
	}
}
