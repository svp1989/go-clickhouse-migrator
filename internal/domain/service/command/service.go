package command

import (
	"context"
	
	"go-clickhouse-migrator/internal/domain/model"
	"go-clickhouse-migrator/internal/domain/repository"
	"go-clickhouse-migrator/internal/domain/service"
)

var _ service.Command = (*Service)(nil)

type Service struct {
	repo repository.Command
}

func New(repo repository.Command) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Up(ctx context.Context, query string) error {
	return s.repo.Up(ctx, query)
}

func (s *Service) Init(ctx context.Context, params model.MigrationQueryParams) error {
	return s.repo.Init(ctx, params)
}

func (s *Service) CurrentVersion(ctx context.Context, params model.MigrationQueryParams) (model.MigrationInfo, error) {
	return s.repo.GetMigrationInfo(ctx, params)
}

func (s *Service) Save(ctx context.Context, migration model.MigrationInfo, params model.MigrationQueryParams) error {
	return s.repo.Save(ctx, migration, params)
}

func (s *Service) ExecutedMigration(ctx context.Context, params model.MigrationQueryParams) (map[string]model.MigrationInfo, error) {
	return s.repo.GetMigrationInfoMap(ctx, params)
}
