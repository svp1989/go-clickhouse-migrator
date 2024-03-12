package migrator

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	
	"go-clickhouse-migrator/internal/domain/model"
	"go-clickhouse-migrator/internal/domain/model/message"
	"go-clickhouse-migrator/internal/domain/service/mocks"
)

const migrationDir = "./migrations"

var errTest = errors.New("test_error")

var queryParams = model.MigrationQueryParams{TableName: migrationDir}

func TestUseCase_diff(t *testing.T) {
	type ReturnMock struct {
		Command     func() (map[string]model.MigrationInfo, error)
		FileManager func() ([]model.MigrationFileInfo, error)
	}
	
	testCases := []struct {
		Name              string
		Mock              ReturnMock
		FileMigrationDiff []model.MigrationFileInfo
		MigrationFileDiff []model.MigrationFileInfo
		Err               error
	}{
		{
			Name: "not exist migration",
			Mock: ReturnMock{
				Command:     func() (map[string]model.MigrationInfo, error) { return nil, nil },
				FileManager: func() ([]model.MigrationFileInfo, error) { return nil, nil },
			},
			FileMigrationDiff: nil,
			MigrationFileDiff: nil,
			Err:               nil,
		},
		{
			Name: "err command query",
			Mock: ReturnMock{
				Command:     func() (map[string]model.MigrationInfo, error) { return nil, errTest },
				FileManager: func() ([]model.MigrationFileInfo, error) { return nil, nil },
			},
			FileMigrationDiff: nil,
			MigrationFileDiff: nil,
			Err:               errTest,
		},
		{
			Name: "err file manager",
			Mock: ReturnMock{
				Command:     func() (map[string]model.MigrationInfo, error) { return nil, nil },
				FileManager: func() ([]model.MigrationFileInfo, error) { return nil, errTest },
			},
			FileMigrationDiff: nil,
			MigrationFileDiff: nil,
			Err:               errTest,
		},
		{
			Name: "migration exist file not exist",
			Mock: ReturnMock{
				Command: func() (map[string]model.MigrationInfo, error) {
					return map[string]model.MigrationInfo{
						"1_one":   {Version: "1_one"},
						"2_two":   {Version: "2_two"},
						"3_three": {Version: "3_three"},
					}, nil
				},
				FileManager: func() ([]model.MigrationFileInfo, error) { return nil, nil },
			},
			FileMigrationDiff: nil,
			MigrationFileDiff: []model.MigrationFileInfo{
				{Version: "1_one"},
				{Version: "2_two"},
				{Version: "3_three"},
			},
			Err: nil,
		},
		{
			Name: "file exist migration not exist",
			Mock: ReturnMock{
				Command: func() (map[string]model.MigrationInfo, error) { return nil, nil },
				FileManager: func() ([]model.MigrationFileInfo, error) {
					return []model.MigrationFileInfo{
						{Version: "1_one"},
						{Version: "2_two"},
						{Version: "3_three"},
					}, nil
				},
			},
			FileMigrationDiff: []model.MigrationFileInfo{
				{Version: "1_one"},
				{Version: "2_two"},
				{Version: "3_three"},
			},
			MigrationFileDiff: nil,
			Err:               nil,
		}, {
			Name: "file exist migration exist",
			Mock: ReturnMock{
				Command: func() (map[string]model.MigrationInfo, error) {
					return map[string]model.MigrationInfo{
						"1_one":   {Version: "1_one"},
						"2_two":   {Version: "2_two"},
						"3_three": {Version: "3_three"},
					}, nil
				},
				FileManager: func() ([]model.MigrationFileInfo, error) {
					return []model.MigrationFileInfo{
						{Version: "1_one"},
						{Version: "2_two"},
						{Version: "3_three"},
					}, nil
				},
			},
			FileMigrationDiff: nil,
			MigrationFileDiff: nil,
			Err:               nil,
		},
	}
	
	migration := mocks.NewMigration(t)
	migration.On("GetQueryParams").Return(queryParams)
	
	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			command := mocks.NewCommand(t)
			fileManager := mocks.NewFileManager(t)
			
			command.On("ExecutedMigration", mock.Anything, queryParams).Return(test.Mock.Command())
			
			if _, err := test.Mock.Command(); err == nil {
				fileManager.On("SortedMigrationFilesData").Return(test.Mock.FileManager())
			}
			
			fileMigrationDiff, migrationFileDiff, err := New(command, fileManager, migration).diff(context.Background())
			
			if err != nil {
				assert.ErrorIs(t, err, test.Err)
			}
			
			assert.Equal(t, test.FileMigrationDiff, fileMigrationDiff)
			assert.Equal(t, test.MigrationFileDiff, migrationFileDiff)
		})
	}
}

func TestUseCase_Up(t *testing.T) {
	migration := mocks.NewMigration(t)
	migration.On("GetQueryParams").Return(queryParams)
	
	fileManager := mocks.NewFileManager(t)
	fileManager.On("Read", model.MigrationFileInfo{Version: "1_one"}).Return("select 1", nil)
	fileManager.On("SortedMigrationFilesData").Return([]model.MigrationFileInfo{{Version: "1_one"}}, nil)
	
	command := mocks.NewCommand(t)
	command.On("Up", mock.Anything, "select 1").Return(nil)
	command.On("Save", mock.Anything, mock.Anything, queryParams).Return(nil)
	command.On("CurrentVersion", mock.Anything, queryParams).Return(model.MigrationInfo{}, nil)
	command.On("ExecutedMigration", mock.Anything, queryParams).Return(nil, nil)
	
	_, err := New(command, fileManager, migration).Up(context.Background(), false)
	assert.NoError(t, err)
}

func TestUseCase_Generate(t *testing.T) {
	migration := mocks.NewMigration(t)
	
	name := "test_1"
	fileManager := mocks.NewFileManager(t)
	fileManager.On("Create", name).Return(name, nil)
	
	command := mocks.NewCommand(t)
	
	_, err := New(command, fileManager, migration).Generate(name)
	
	assert.NoError(t, err)
	
	fileManager.On("Create", "test2").Return("", errTest)
	_, err = New(command, fileManager, migration).Generate("test2")
	
	assert.Equal(t, errTest, err)
}

func TestUseCase_Init(t *testing.T) {
	migration := mocks.NewMigration(t)
	migration.On("GetQueryParams").Return(queryParams)
	
	fileManager := mocks.NewFileManager(t)
	
	command := mocks.NewCommand(t)
	command.On("Init", mock.Anything, queryParams).Return(nil)
	
	_, err := New(command, fileManager, migration).Init(context.Background())
	assert.NoError(t, err)
	
	command = mocks.NewCommand(t)
	command.On("Init", mock.Anything, queryParams).Return(errTest)
	
	_, err = New(command, fileManager, migration).Init(context.Background())
	assert.Error(t, errTest, err)
}

func TestUseCase_Diff(t *testing.T) {
	migration := mocks.NewMigration(t)
	migration.On("GetQueryParams").Return(queryParams)
	
	fileManager := mocks.NewFileManager(t)
	fileManager.On("SortedMigrationFilesData").Return(nil, nil)
	
	command := mocks.NewCommand(t)
	command.On("ExecutedMigration", mock.Anything, queryParams).Return(nil, nil)
	
	_, err := New(command, fileManager, migration).Diff(context.Background())
	assert.NoError(t, err)
	
	command = mocks.NewCommand(t)
	command.On("ExecutedMigration", mock.Anything, queryParams).Return(nil, errTest)
	
	_, err = New(command, fileManager, migration).Diff(context.Background())
	assert.Error(t, errTest, err)
}

func TestUseCase_Version(t *testing.T) {
	testCases := []struct {
		Name     string
		Input    model.MigrationInfo
		Expected message.ConsoleMessage
		Err      error
	}{
		{
			Name:     "empty migration",
			Input:    model.MigrationInfo{},
			Expected: message.ConsoleMessage{Message: message.MigrationsLastVersionInfo, Type: "info", Data: message.Data{Key: "migration", Value: ""}},
			Err:      nil,
		},
		{
			Name:     "err get version query",
			Input:    model.MigrationInfo{},
			Expected: message.ConsoleMessage{},
			Err:      errTest,
		},
		{
			Name:     "no rows",
			Input:    model.MigrationInfo{},
			Expected: message.ConsoleMessage{Message: message.MigrationsNotFoundWarning, Type: message.Warning, Data: message.Data{}},
			Err:      sql.ErrNoRows,
		},
		{
			Name:     "return last version",
			Input:    model.MigrationInfo{Version: "one"},
			Expected: message.ConsoleMessage{Message: message.MigrationsLastVersionInfo, Type: message.Info, Data: message.Data{Key: "migration", Value: "one"}},
			Err:      nil,
		},
	}
	
	migration := mocks.NewMigration(t)
	migration.On("GetQueryParams").Return(queryParams)
	
	fileManager := mocks.NewFileManager(t)
	
	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			command := mocks.NewCommand(t)
			command.On("CurrentVersion", mock.Anything, queryParams).Return(testCase.Input, testCase.Err)
			
			version, err := New(command, fileManager, migration).Version(context.Background())
			
			assert.Equal(t, testCase.Expected, version)
			
			if !errors.Is(testCase.Err, sql.ErrNoRows) {
				assert.Equal(t, testCase.Err, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
