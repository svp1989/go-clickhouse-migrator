package clickhouse

import "errors"

var (
	ErrConnectionFailed = errors.New("failed to connect ClickHouse")
	ErrScan             = errors.New("failed to scan row")
	ErrRowsNext         = errors.New("failed to fetch next rows")
)
