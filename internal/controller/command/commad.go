package command

import (
	"log/slog"
	
	"go-clickhouse-migrator/internal/domain/usecase"
)

type Command struct {
	logger   *slog.Logger
	migrator usecase.Migrator
}

func New(logger *slog.Logger, migrator usecase.Migrator) *Command {
	return &Command{
		logger:   logger,
		migrator: migrator,
	}
}
