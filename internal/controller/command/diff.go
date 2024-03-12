package command

import (
	"context"
	
	"go-clickhouse-migrator/pkg/tools"
)

func (c *Command) Diff(ctx context.Context) error {
	messages, err := c.migrator.Diff(ctx)
	if err != nil {
		return err
	}
	
	for _, msg := range messages {
		tools.PrintToConsole(c.logger, msg)
	}
	
	return nil
}
