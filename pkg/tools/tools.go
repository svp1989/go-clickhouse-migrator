package tools

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"text/template"
	
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	
	"go-clickhouse-migrator/internal/domain/model"
	"go-clickhouse-migrator/internal/domain/model/message"
)

func GetDynamicQuery(query string, params *model.MigrationQueryParams) (string, error) {
	t, err := template.New("").Parse(query)
	if err != nil {
		return "", fmt.Errorf("failed parse query %w", err)
	}
	
	result := &bytes.Buffer{}
	if err := t.Execute(result, params); err != nil {
		return "", fmt.Errorf("failed execute %w", err)
	}
	
	return result.String(), nil
}

func PrintToConsole(logger *slog.Logger, msg message.ConsoleMessage) {
	switch msg.Type {
	case message.Info:
		logger.Info(msg.Message, msg.Data.Key, msg.Data.Value)
	case message.Warning:
		logger.Warn(msg.Message, msg.Data.Key, msg.Data.Value)
	case message.Success:
		logger.Info(msg.Message, msg.Data.Key, msg.Data.Value)
	case message.Error:
		logger.Error(msg.Message, msg.Data.Key, msg.Data.Value)
	default:
		logger.Error(msg.Message, msg.Data.Key, msg.Data.Value)
	}
}

const defaultEnvFilePath = ".env"

func ProcessEnv(prefix string, cfg any) error {
	// Load env vars from a file if it exists
	// Return an error if file is present but couldn't be parsed
	err := godotenv.Load(defaultEnvFilePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	
	// Load config
	if err := envconfig.Process(prefix, cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}
	
	return nil
}
