package command

import (
	"context"
	"errors"
	
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	
	"go-clickhouse-migrator/internal/domain/model"
	"go-clickhouse-migrator/internal/domain/repository"
	"go-clickhouse-migrator/pkg/clickhouse"
	"go-clickhouse-migrator/pkg/tools"
)

var _ repository.Command = (*Repository)(nil)

type Repository struct {
	conn driver.Conn
}

func New(conn driver.Conn) *Repository {
	return &Repository{
		conn: conn,
	}
}

func (c *Repository) Up(ctx context.Context, query string) error {
	err := c.conn.Exec(ctx, query)
	if err != nil {
		return err
	}
	
	return nil
}

func (c *Repository) Init(ctx context.Context, params model.MigrationQueryParams) error {
	query, err := tools.GetDynamicQuery(initMigrationVersionTable, &params)
	if err != nil {
		return err
	}
	
	err = c.conn.Exec(ctx, query)
	if err != nil {
		return err
	}
	
	return nil
}

func (c *Repository) Save(ctx context.Context, data model.MigrationInfo, params model.MigrationQueryParams) error {
	query, err := tools.GetDynamicQuery(insertMigrationVersion, &params)
	if err != nil {
		return err
	}
	
	batch, err := c.conn.PrepareBatch(ctx, query)
	if err != nil {
		return err
	}
	
	err = batch.Append(data.Version, data.ExecutedAt, data.ExecutionTime, data.Error)
	if err != nil {
		return err
	}
	
	err = batch.Send()
	if err != nil {
		return err
	}
	
	return nil
}

func (c *Repository) GetMigrationInfo(ctx context.Context, params model.MigrationQueryParams) (model.MigrationInfo, error) {
	query, err := tools.GetDynamicQuery(getMigrationInfo, &params)
	if err != nil {
		return model.MigrationInfo{}, err
	}
	
	row := c.conn.QueryRow(ctx, query)
	
	var data MigrationInfo
	
	if err := row.ScanStruct(&data); err != nil {
		return model.MigrationInfo{}, err
	}
	
	return model.MigrationInfo(data), nil
}

func (c *Repository) GetMigrationInfoMap(ctx context.Context, params model.MigrationQueryParams) (map[string]model.MigrationInfo, error) {
	query, err := tools.GetDynamicQuery(getMigrationInfoList, &params)
	if err != nil {
		return nil, err
	}
	
	rows, err := c.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	
	migrationDataMap := make(map[string]model.MigrationInfo)
	
	for rows.Next() {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		
		var migrationData MigrationInfo
		if err := rows.ScanStruct(&migrationData); err != nil {
			return nil, errors.Join(clickhouse.ErrScan, err)
		}
		
		migrationDataMap[migrationData.Version] = model.MigrationInfo(migrationData)
	}
	
	if err := rows.Err(); err != nil {
		return nil, errors.Join(clickhouse.ErrRowsNext, err)
	}
	
	return migrationDataMap, nil
}
