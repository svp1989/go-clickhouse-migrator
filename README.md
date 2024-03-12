# ClickHouse migrator
Clickhouse Migration Application Tool

## How to use it
The configuration of the migrator is set in the environment variables

```dotenv
#Коннект к кликхаусу
MIGRATOR_CLICKHOUSE_SERVER:   localhost #clickhouse server
MIGRATOR_CLICKHOUSE_PORT:     9000 #port (default is 9000)
MIGRATOR_CLICKHOUSE_DATABASE: default #database ("default" by default)
MIGRATOR_CLICKHOUSE_USER:     admin #user
MIGRATOR_CLICKHOUSE_PASSWORD: 123 #password

#Migratory Settings
MIGRATOR_MIGRATIONS_DIRECTORY: ./migrations #the path to the directory with the migration files
MIGRATOR_MIGRATIONS_TABLE:     migration_versions #the name of the table used to track the versioning of migrations
```

### Applying migrations to gitlab-ci in third-party packages:
```bash
#!/usr/bin/env bash
set -e

echo Starting clickhouse migrations;

echo $SECRET | docker login --username $LOGIN $REGISTRY --password-stdin

# Test positive case
docker run --init --rm -v ${CI_PROJECT_DIR}/migrations:/migrations \
    -e MIGRATOR_CLICKHOUSE_SERVER=${MIGRATOR_CLICKHOUSE_SERVER} \
    -e MIGRATOR_CLICKHOUSE_PORT=${MIGRATOR_CLICKHOUSE_PORT} \
    -e MIGRATOR_CLICKHOUSE_USER=${MIGRATOR_CLICKHOUSE_USER} \
    -e MIGRATOR_CLICKHOUSE_PASSWORD=${MIGRATOR_CLICKHOUSE_PASSWORD} \
    -e MIGRATOR_MIGRATION_DIRECTORY=${MIGRATOR_MIGRATION_DIRECTORY} \
    -e MIGRATOR_MIGRATION_TABLE=${MIGRATOR_MIGRATION_TABLE} \
    $IMAGE_MIGRATION up || error_code=$?

if [ ! -z $error_code ]; then
        exit 1;
fi

```
- `$IMAGE_MIGRATION` - the name of the migration image in harbor
- `-v ${CI_PROJECT_DIR}/migrations:` - directory with third-party application migrations

### Migrator Commands

[go-clickhouse-migrator](bin%2Fgo-clickhouse-migrator)
- `help` - description of the migration commands
- `up` - is the command for applying migrations
- `diff` - shows which migrations have not been applied and which have been applied but do not have migration files
- `init` - creates a table in which migrations will be stored. The table name is taken from the variable ${MIGRATOR_MIGRATION_TABLE}
- `version` - the last migration applied
- `gen`     - creates a new migration the second argument passes the name of the migration. You should only generate new migrations using this command
