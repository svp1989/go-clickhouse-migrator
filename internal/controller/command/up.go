package command

import (
	"context"
	"fmt"
	
	"go-clickhouse-migrator/internal/domain/model/message"
	"go-clickhouse-migrator/pkg/tools"
)

const optionForce = "--force"

func (c *Command) Up(ctx context.Context, args []string) error {
	force := false
	if len(args[1:]) != 0 && args[1] == optionForce {
		force = true
	}
	
	messages, err := c.migrator.Up(ctx, force)
	
	isErrorMessage := false
	
	for _, msg := range messages {
		if msg.Type == message.Error && !force {
			isErrorMessage = true
		}
		
		tools.PrintToConsole(c.logger, msg)
	}
	
	if isErrorMessage {
		return fmt.Errorf("please check migration files")
	}
	
	return err
}
