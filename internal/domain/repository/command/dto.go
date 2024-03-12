package command

import "time"

type MigrationInfo struct {
	Version       string    `ch:"version"`
	ExecutedAt    time.Time `ch:"executed_at"`
	ExecutionTime uint64    `ch:"execution_time"`
	Error         string    `ch:"error"`
}
