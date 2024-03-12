package filemanager

import (
	"os"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"go-clickhouse-migrator/internal/domain/model"
	"go-clickhouse-migrator/pkg/migrator"
)

const migrationDir = "./migrations"

var cfg = migrator.Config{Dir: migrationDir}

func testBeforeRun(t testing.TB) error {
	t.Cleanup(func() {
		os.RemoveAll(migrationDir)
	})
	
	return os.Mkdir(migrationDir, os.ModePerm)
}

func TestService_Create(t *testing.T) {
	assert.NoError(t, testBeforeRun(t))
	
	fileManager := New(&cfg)
	for _, name := range []string{"test_create", "test_create_1", "test_create_2"} {
		createName, err := fileManager.Create(name)
		
		assert.NoError(t, err)
		assert.Contains(t, createName, name)
	}
}

func TestService_Read(t *testing.T) {
	assert.NoError(t, testBeforeRun(t))
	
	fileManager := New(&cfg)
	createName, err := fileManager.Create("test_create")
	assert.NoError(t, err)
	
	read, err := fileManager.Read(model.MigrationFileInfo{Version: createName, Dir: migrationDir})
	assert.NoError(t, err)
	assert.Equal(t, read, templateMigrationFile)
}

func TestService_SortedMigrationFilesData(t *testing.T) {
	assert.NoError(t, testBeforeRun(t))
	
	fileManager := New(&cfg)
	
	sortedCreatedFiles := make([]model.MigrationFileInfo, 0, 3)
	
	for _, name := range []string{"first", "second", "three"} {
		createName, err := fileManager.Create(name)
		
		assert.NoError(t, err)
		assert.Contains(t, createName, name)
		
		sortedCreatedFiles = append(sortedCreatedFiles, model.MigrationFileInfo{
			Version: createName,
			Dir:     migrationDir,
		})
	}
	
	data, err := fileManager.SortedMigrationFilesData()
	assert.NoError(t, err)
	
	for i, fileData := range sortedCreatedFiles {
		assert.Equal(t, fileData, data[i])
	}
}
