package command

const initMigrationVersionTable = `
create table if not exists {{ .TableName }} (
    version String,
    executed_at DateTime64(3, 'Europe/Moscow'),
    execution_time UInt64,
    error String
) engine = TinyLog;
`
const getMigrationInfo = `
select version, executed_at, execution_time
from {{ .TableName }}
where error = ''
order by version desc
limit 1
`

const getMigrationInfoList = `
select version, executed_at, execution_time
from {{ .TableName }}
where error = ''
order by version
`

const insertMigrationVersion = `insert into {{ .TableName }}`
