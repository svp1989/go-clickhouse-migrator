package migrator

// Config - конфиг мигратора
// Dir дирректория в которой лежат файлы миграции
// Table название таблицы для хранения данных о миграции
type Config struct {
	Dir   string `envconfig:"directory" default:"./migrations"`
	Table string `envconfig:"table" default:"migration_versions"`
}
