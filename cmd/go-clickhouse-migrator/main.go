package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	
	"go-clickhouse-migrator/internal/config"
	"go-clickhouse-migrator/internal/controller"
	"go-clickhouse-migrator/internal/controller/command"
	commandRepository "go-clickhouse-migrator/internal/domain/repository/command"
	commandService "go-clickhouse-migrator/internal/domain/service/command"
	filemanagerService "go-clickhouse-migrator/internal/domain/service/filemanager"
	migrationService "go-clickhouse-migrator/internal/domain/service/migration"
	"go-clickhouse-migrator/internal/domain/usecase/migrator"
	"go-clickhouse-migrator/pkg/clickhouse"
	"go-clickhouse-migrator/pkg/tools"
)

const envPrefix = "MIGRATOR"

func main() {
	var errorCode int
	defer func() { os.Exit(errorCode) }()
	
	logger := slog.Default()
	
	defer func() {
		if err := recover(); err != nil {
			logger.Error("❌ panic recovered in main", "err", err)
			
			errorCode = 1
		}
	}()
	
	flag.Parse()
	
	if err := run(flag.Args(), logger); err != nil {
		logger.Error("❌ execution failed", "err", err)
		
		errorCode = 1
	}
}

func run(args []string, logger *slog.Logger) (err error) {
	cfg := &config.Config{}
	if err := tools.ProcessEnv(envPrefix, cfg); err != nil {
		return fmt.Errorf("❌ failed to parse config: %w", err)
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	connect, err := clickhouse.Connect(&cfg.ClickHouse)
	if err != nil {
		return fmt.Errorf("❌ failed to connect: %w", err)
	}
	
	defer func() {
		if err := connect.Close(); err != nil {
			slog.Error("❌ error closing clickhouse", "err", err)
		}
	}()
	
	ctrl := controller.New(
		command.New(logger,
			migrator.New(
				commandService.New(commandRepository.New(connect)),
				filemanagerService.New(&cfg.Migrator),
				migrationService.New(&cfg.Migrator),
			),
		),
	)
	
	err = ctrl.Handle(ctx, args)
	if err != nil {
		return err
	}
	
	return nil
}
