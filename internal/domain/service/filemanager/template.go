package filemanager

const templateMigrationFile = `-- if you need create table please use create table if not exist
create table if not exists {{ tableName }}
`
