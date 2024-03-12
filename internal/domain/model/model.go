package model

import (
	"time"
)

// MigrationInfo - модель хранит информацию о миграции
// Version версия миграция
// ExecutedAt дата выполнения миграции
// ExecutionTime время выполнения миграции
// Error ошибка в случае ее возникновения
type MigrationInfo struct {
	Version       string
	ExecutedAt    time.Time
	ExecutionTime uint64
	Error         string
}

func (m *MigrationInfo) FromMigrationFileInfo(fileInfo MigrationFileInfo, endTime time.Time, message string) MigrationInfo {
	executionTime := time.Since(endTime).Seconds()

	return MigrationInfo{
		Version:       fileInfo.Version,
		Error:         message,
		ExecutedAt:    time.Now(),
		ExecutionTime: uint64(executionTime),
	}
}

// MigrationQueryParams - динамичесике параметры миграции
// Table название таблицы по умолчанию migration_version
// меняется с помощью переменной окружения MIGRATOR_MIGRATOR_TABLE
type MigrationQueryParams struct {
	TableName string
}

// MigrationFileInfo - данные файла для миграции
// Version имя файла/версия миграции
// Dir директория где лежит файл
type MigrationFileInfo struct {
	Version string
	Dir     string
}
