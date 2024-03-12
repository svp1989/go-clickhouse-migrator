package command

import (
	"fmt"
	
	"go-clickhouse-migrator/internal/domain/model/message"
	"go-clickhouse-migrator/pkg/tools"
)

func (c *Command) Generate(args []string) error {
	if len(args[1:]) == 0 {
		return fmt.Errorf(message.ErrEmptyMigrationName)
	}
	
	migrationName := args[1]
	
	msg, err := c.migrator.Generate(migrationName)
	if err != nil {
		return err
	}
	
	tools.PrintToConsole(c.logger, msg)
	
	return nil
}
