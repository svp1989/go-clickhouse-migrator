package command

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"
	
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/stretchr/testify/assert"
	
	"go-clickhouse-migrator/internal/domain/model"
	"go-clickhouse-migrator/pkg/clickhouse"
)

const testTableName = "test_migrations_schema"

var queryParams = model.MigrationQueryParams{TableName: testTableName}

var cfg = clickhouse.Config{
	Server:   "localhost",
	Port:     "9000",
	Database: "",
	User:     "admin",
	Password: "123",
}

func getConnectBeforeTest(t *testing.T) driver.Conn {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
	
	conn, err := clickhouse.Connect(&cfg)
	assert.NoError(t, err)
	
	err = conn.Exec(context.Background(), fmt.Sprintf("drop table if exists default.%s sync", testTableName))
	assert.NoError(t, err)
	
	repo := New(conn)
	err = repo.Init(context.Background(), queryParams)
	
	assert.NoError(t, err)
	
	return conn
}

func getTestTime() time.Time {
	loc, _ := time.LoadLocation("Europe/Moscow")
	
	return time.Date(2023, time.January, 1, 0, 0, 0, 0, loc)
}

func TestRepository_Up(t *testing.T) {
	assert.NoError(t,
		New(getConnectBeforeTest(t)).Up(context.Background(), `select 1`),
	)
}

func TestRepository_GetMigrationInfo(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		error    error
		expected model.MigrationInfo
	}{
		{
			name:     "empty table",
			query:    `select 1`,
			error:    sql.ErrNoRows,
			expected: model.MigrationInfo{},
		},
		{
			name: "table don't have success migration",
			query: fmt.Sprintf(`INSERT INTO %s (version, executed_at, execution_time, error)
						VALUES ('version', toDateTime('2023-01-01'), 10, 'some error');`, testTableName),
			error:    sql.ErrNoRows,
			expected: model.MigrationInfo{},
		},
		{
			name: "success migration",
			query: fmt.Sprintf(`INSERT INTO %s (version, executed_at, execution_time, error)
						VALUES ('version', toDateTime('2023-01-01'), 10, '');`, testTableName),
			error: nil,
			expected: model.MigrationInfo{
				Version:       "version",
				ExecutedAt:    getTestTime(),
				ExecutionTime: 10,
				Error:         "",
			},
		},
	}
	
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctx := context.Background()
			conn := getConnectBeforeTest(t)
			assert.NoError(t, conn.Exec(ctx, testCase.query))
			
			repo := New(conn)
			actual, err := repo.GetMigrationInfo(ctx, queryParams)
			
			if testCase.error != nil {
				assert.ErrorIs(t, err, testCase.error)
			}
			
			assert.Equal(t, actual, testCase.expected)
		})
	}
}

func TestRepository_Save(t *testing.T) {
	ctx := context.Background()
	conn := getConnectBeforeTest(t)
	
	repo := New(conn)
	
	info := model.MigrationInfo{
		Version:       "1234_test.sql",
		ExecutionTime: 154,
		ExecutedAt:    getTestTime(),
		Error:         "some error",
	}
	err := repo.Save(ctx, info, queryParams)
	
	assert.NoError(t, err)
	
	row := conn.QueryRow(ctx, fmt.Sprintf("select version, executed_at, execution_time, error from %s", testTableName))
	dto := MigrationInfo{}
	
	assert.NoError(t, row.ScanStruct(&dto))
	assert.Equal(t, dto, MigrationInfo{
		Version:       info.Version,
		ExecutionTime: info.ExecutionTime,
		ExecutedAt:    getTestTime(),
		Error:         info.Error,
	})
}

func TestRepository_GetMigrationInfoMap(t *testing.T) {
	info := []model.MigrationInfo{
		{Version: "1234_test.sql", ExecutionTime: 312, ExecutedAt: getTestTime(), Error: "some error"},
		{Version: "4235235_test.sql", ExecutionTime: 2, ExecutedAt: getTestTime(), Error: ""},
		{Version: "123123_test.sql", ExecutionTime: 6, ExecutedAt: getTestTime(), Error: ""},
		{Version: "5555_test.sql", ExecutionTime: 1, ExecutedAt: getTestTime(), Error: ""},
	}
	
	ctx := context.Background()
	conn := getConnectBeforeTest(t)
	
	repo := New(conn)
	for _, migrationInfo := range info {
		assert.NoError(t, repo.Save(ctx, migrationInfo, queryParams))
	}
	
	actual, err := repo.GetMigrationInfoMap(ctx, queryParams)
	assert.NoError(t, err)
	assert.Equal(t, actual, map[string]model.MigrationInfo{
		"4235235_test.sql": {Version: "4235235_test.sql", ExecutionTime: 2, ExecutedAt: getTestTime(), Error: ""},
		"123123_test.sql":  {Version: "123123_test.sql", ExecutionTime: 6, ExecutedAt: getTestTime(), Error: ""},
		"5555_test.sql":    {Version: "5555_test.sql", ExecutionTime: 1, ExecutedAt: getTestTime(), Error: ""},
	})
}
