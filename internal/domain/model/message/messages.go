package message

const (
	MigrationsNotFoundSuccess      = "✅ new migration not found"
	MigrationTableCreated          = "✅ migration table created"
	MigrationsFileGeneratedSuccess = "✅ migration file generated"
	MigrationsExecuted             = "✅ migration executed"

	MigrationsNotExecutedInfo = "🔹 migration not executed"
	MigrationsLastVersionInfo = "🔹 latest migration version"

	MigrationsFilesNotFoundWarning = "❗️migration file not found"
	MigrationsNotFoundWarning      = "❗️executed migration not found"
)

type ConsoleMessage struct {
	Message string
	Type    Type
	Data    Data
}

type Data struct {
	Key   string
	Value any
}
