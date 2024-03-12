package command

import (
	"context"
	
	"go-clickhouse-migrator/pkg/tools"
)

func (c *Command) Version(ctx context.Context) error {
	msg, err := c.migrator.Version(ctx)
	if err != nil {
		return err
	}
	
	tools.PrintToConsole(c.logger, msg)
	
	return nil
}
