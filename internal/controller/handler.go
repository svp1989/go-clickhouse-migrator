package controller

import (
	"context"
	"fmt"
	
	"go-clickhouse-migrator/internal/controller/command"
	"go-clickhouse-migrator/internal/domain/model/message"
)

const (
	commandUp       = "up"
	commandVersion  = "version"
	commandGenerate = "gen"
	commandInit     = "init"
	commandDiff     = "diff"
	commandHelp     = "help"
)

var allowedCommands = []string{
	commandUp,
	commandVersion,
	commandGenerate,
	commandInit,
	commandDiff,
	commandHelp,
}

type Handler struct {
	command command.Command
}

func New(cmd *command.Command) *Handler {
	return &Handler{
		command: *cmd,
	}
}

func (h *Handler) Handle(ctx context.Context, args []string) error {
	if err := validateArgs(args); err != nil {
		return err
	}
	
	switch args[0] {
	case commandUp:
		return h.command.Up(ctx, args)
	case commandVersion:
		return h.command.Version(ctx)
	case commandGenerate:
		return h.command.Generate(args)
	case commandInit:
		return h.command.Init(ctx)
	case commandDiff:
		return h.command.Diff(ctx)
	case commandHelp:
		return h.command.Help()
	default:
		return fmt.Errorf(message.ErrInvalidCommand)
	}
}

func validateArgs(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf(message.ErrInvalidCommand)
	}
	
	for _, allowedCommand := range allowedCommands {
		if args[0] == allowedCommand {
			return nil
		}
	}
	
	return fmt.Errorf(message.ErrInvalidCommand)
}
