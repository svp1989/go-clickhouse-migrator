package usecase

import (
	"context"
	
	"go-clickhouse-migrator/internal/domain/model/message"
)

type Migrator interface {
	Up(ctx context.Context, force bool) ([]message.ConsoleMessage, error)
	Version(ctx context.Context) (message.ConsoleMessage, error)
	Init(ctx context.Context) (message.ConsoleMessage, error)
	Generate(name string) (message.ConsoleMessage, error)
	Diff(ctx context.Context) ([]message.ConsoleMessage, error)
	Help(file string) (string, error)
}
