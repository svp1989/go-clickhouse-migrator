package filemanager

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
	
	"go-clickhouse-migrator/internal/domain/model"
	"go-clickhouse-migrator/internal/domain/service"
	"go-clickhouse-migrator/pkg/migrator"
)

var _ service.FileManager = (*Service)(nil)

type Service struct {
	baseDir string
}

func New(cfg *migrator.Config) *Service {
	return &Service{
		baseDir: cfg.Dir,
	}
}

// SortedMigrationFilesData - возвращает файлы для миграции в отсортированном виде
func (s *Service) SortedMigrationFilesData() ([]model.MigrationFileInfo, error) {
	files, err := os.ReadDir(s.baseDir)
	if err != nil {
		return nil, err
	}
	
	filesData := make([]model.MigrationFileInfo, 0, len(files))
	for _, file := range files {
		filesData = append(filesData, model.MigrationFileInfo{
			Version: file.Name(),
			Dir:     s.baseDir,
		})
	}
	
	sort.Slice(filesData, func(i, j int) bool {
		return filesData[i].Version < filesData[j].Version
	})
	
	return filesData, nil
}

func (s *Service) Read(fileData model.MigrationFileInfo) (string, error) {
	readFile, err := os.ReadFile(filepath.Join(fileData.Dir, string(os.PathSeparator), fileData.Version))
	if err != nil {
		return "", err
	}
	
	return string(readFile), nil
}

func (s *Service) Create(name string) (string, error) {
	fileName := fmt.Sprintf("%v_%s.up.sql", time.Now().Unix(), name)
	
	path := filepath.Join(s.baseDir, string(os.PathSeparator), fileName)
	
	if _, err := os.Create(path); err != nil {
		return "", err
	}
	
	if err := os.WriteFile(path, []byte(templateMigrationFile), os.ModePerm); err != nil {
		return "", err
	}
	
	return fileName, nil
}

func (s *Service) ReadHelp(file string) (string, error) {
	readFile, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	
	return string(readFile), nil
}
