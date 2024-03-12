package clickhouse

import (
	"context"
	"errors"
	"fmt"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

func Connect(cfg *Config) (driver.Conn, error) {
	options := getOptions(cfg)

	conn, err := ch.Open(&options)
	if err != nil {
		return nil, err
	}

	if err = conn.Ping(context.Background()); err != nil {
		var exception *ch.Exception

		if errors.As(err, &exception) {
			except := fmt.Errorf("exception [%d] %s \n%s", exception.Code, exception.Message, exception.StackTrace)

			return nil, errors.Join(ErrConnectionFailed, except, err)
		}

		return nil, errors.Join(ErrConnectionFailed, err)
	}

	return conn, nil
}

func getOptions(cfg *Config) ch.Options {
	return ch.Options{
		Addr: []string{fmt.Sprintf("%s:%s", cfg.Server, cfg.Port)},
		Auth: ch.Auth{
			Database: cfg.Database,
			Username: cfg.User,
			Password: cfg.Password,
		},
	}
}
