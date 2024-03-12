package config

import (
	"go-clickhouse-migrator/pkg/clickhouse"
	"go-clickhouse-migrator/pkg/migrator"
)

type Config struct {
	ClickHouse clickhouse.Config `envconfig:"clickhouse"`
	Migrator   migrator.Config   `envconfig:"migration"`
}
