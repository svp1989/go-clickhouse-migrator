package migrator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"time"
	
	"go-clickhouse-migrator/internal/domain/model"
	"go-clickhouse-migrator/internal/domain/model/message"
	"go-clickhouse-migrator/internal/domain/service"
	"go-clickhouse-migrator/internal/domain/usecase"
)

var _ usecase.Migrator = (*UseCase)(nil)

type UseCase struct {
	command     service.Command
	fileManager service.FileManager
	migration   service.Migration
}

func New(command service.Command, fileManager service.FileManager, migration service.Migration) *UseCase {
	return &UseCase{
		command:     command,
		fileManager: fileManager,
		migration:   migration,
	}
}

// Up - накатывает миграции
// -force  необязательный параметр накатывает миграции без проверок
func (uc *UseCase) Up(ctx context.Context, force bool) ([]message.ConsoleMessage, error) {
	migrationsNotExecuted, fileNotFound, err := uc.diff(ctx)
	if err != nil {
		return nil, err
	}
	
	if fileNotFound != nil && !force {
		messages := make([]message.ConsoleMessage, 0, len(fileNotFound))
		for _, info := range fileNotFound {
			messages = append(messages, message.ConsoleMessage{
				Message: message.MigrationsFilesNotFoundWarning,
				Type:    message.Error,
				Data:    message.Data{Key: "file_name", Value: info.Version},
			})
		}
		
		return messages, nil
	}
	
	version, err := uc.command.CurrentVersion(ctx, uc.migration.GetQueryParams())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	
	messages := make([]message.ConsoleMessage, 0, len(migrationsNotExecuted))
	
	for _, currentMigration := range migrationsNotExecuted {
		if currentMigration.Version < version.Version && !force {
			return nil, fmt.Errorf(message.ErrVersionMigration+"%s < %s", currentMigration.Version, version.Version)
		}
		
		start := time.Now()
		
		query, err := uc.fileManager.Read(currentMigration)
		if err != nil {
			return messages, err
		}
		
		if err := uc.command.Up(ctx, query); err != nil {
			if errSave := uc.saveMigrationInfo(ctx, currentMigration, start, err.Error()); errSave != nil {
				return messages, errors.Join(err, errSave)
			}
			
			return messages, err
		}
		
		if err := uc.saveMigrationInfo(ctx, currentMigration, start, ""); err != nil {
			return messages, err
		}
		
		messages = append(messages, message.ConsoleMessage{
			Message: message.MigrationsExecuted,
			Type:    message.Success,
			Data:    message.Data{Key: "migration", Value: currentMigration},
		})
	}
	
	return messages, nil
}

// Version - возвращает последнюю версию миграции
func (uc *UseCase) Version(ctx context.Context) (message.ConsoleMessage, error) {
	currentVersionInfo, err := uc.command.CurrentVersion(ctx, uc.migration.GetQueryParams())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return message.ConsoleMessage{Message: message.MigrationsNotFoundWarning, Type: message.Warning}, nil
		}
		
		return message.ConsoleMessage{}, err
	}
	
	return message.ConsoleMessage{
		Message: message.MigrationsLastVersionInfo,
		Type:    message.Info,
		Data:    message.Data{Key: "migration", Value: currentVersionInfo.Version},
	}, nil
}

// Init - создает таблицу для хранения миграции
func (uc *UseCase) Init(ctx context.Context) (message.ConsoleMessage, error) {
	err := uc.command.Init(ctx, uc.migration.GetQueryParams())
	if err != nil {
		return message.ConsoleMessage{}, err
	}
	
	return message.ConsoleMessage{
		Message: message.MigrationTableCreated,
		Type:    message.Success,
		Data:    message.Data{Key: "table", Value: uc.migration.GetQueryParams().TableName},
	}, nil
}

// Generate - создает новую миграцию
// name - аргумент для создания миграции
func (uc *UseCase) Generate(name string) (message.ConsoleMessage, error) {
	fileName, err := uc.fileManager.Create(name)
	if err != nil {
		return message.ConsoleMessage{}, err
	}
	
	return message.ConsoleMessage{
		Message: message.MigrationsFileGeneratedSuccess,
		Type:    message.Success,
		Data:    message.Data{Key: "file_name", Value: fileName},
	}, nil
}

func (uc *UseCase) Help(file string) (string, error) {
	return uc.fileManager.ReadHelp(file)
}

func (uc *UseCase) Diff(ctx context.Context) ([]message.ConsoleMessage, error) {
	migrationsNotExecuted, migrationsFilesNotFound, err := uc.diff(ctx)
	if err != nil {
		return nil, err
	}
	
	if migrationsNotExecuted == nil && migrationsFilesNotFound == nil {
		return []message.ConsoleMessage{
			{
				Message: message.MigrationsNotFoundSuccess,
				Type:    message.Success,
			},
		}, nil
	}
	
	diff := make([]message.ConsoleMessage, 0, len(migrationsNotExecuted)+len(migrationsFilesNotFound))
	for _, info := range migrationsNotExecuted {
		diff = append(diff, message.ConsoleMessage{
			Message: message.MigrationsNotExecutedInfo,
			Type:    message.Info,
			Data:    message.Data{Key: "migration", Value: info.Version},
		})
	}
	
	for _, info := range migrationsFilesNotFound {
		diff = append(diff, message.ConsoleMessage{
			Message: message.MigrationsFilesNotFoundWarning,
			Type:    message.Warning,
			Data:    message.Data{Key: "migration", Value: info.Version},
		})
	}
	
	return diff, nil
}

// diff - возвращает разницу между выполенными и исполняемыми миграциями
func (uc *UseCase) diff(ctx context.Context) (migrationsNotExecuted, migrationsFilesNotFound []model.MigrationFileInfo, err error) {
	executedMigration, err := uc.command.ExecutedMigration(ctx, uc.migration.GetQueryParams())
	if err != nil {
		return nil, nil, err
	}
	
	migrationsFilesData, err := uc.fileManager.SortedMigrationFilesData()
	if err != nil {
		return nil, nil, err
	}
	
	for _, currentMigration := range migrationsFilesData {
		if _, ok := executedMigration[currentMigration.Version]; !ok {
			migrationsNotExecuted = append(migrationsNotExecuted, currentMigration)
		}
		
		delete(executedMigration, currentMigration.Version)
	}
	
	for version := range executedMigration {
		migrationsFilesNotFound = append(migrationsFilesNotFound, model.MigrationFileInfo{
			Version: version,
		})
	}
	
	sort.Slice(migrationsFilesNotFound, func(i, j int) bool {
		return migrationsFilesNotFound[i].Version < migrationsFilesNotFound[j].Version
	})
	
	return migrationsNotExecuted, migrationsFilesNotFound, nil
}

// saveMigrationInfo - сохраняем информацию о миграции
func (uc *UseCase) saveMigrationInfo(ctx context.Context, fileInfo model.MigrationFileInfo, start time.Time, errMessage string) error {
	migrationData := new(model.MigrationInfo).FromMigrationFileInfo(fileInfo, start, errMessage)
	
	return uc.command.Save(ctx, migrationData, uc.migration.GetQueryParams())
}
