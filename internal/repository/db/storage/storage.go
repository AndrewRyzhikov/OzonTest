package storage

import (
	"context"
	"database/sql"
	"fmt"

	"OzonTest/internal/config"
)

type Storage[T any] struct {
	*sql.DB
}

func NewStorage[T any](cfg config.DataBaseConfig) (*Storage[T], error) {
	db, err := sql.Open(cfg.Driver, dataToPSQLConnection(cfg.Port, cfg.Host, cfg.User, cfg.Password, cfg.DbName))
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}
	return &Storage[T]{db}, nil
}

func dataToPSQLConnection(port uint64, host, user, password, dbname string) string {
	return fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func (d *Storage[T]) executeQueryRow(ctx context.Context, query string, args ...interface{}) (T, error) {
	var record T

	tx, err := d.Begin()
	if err != nil {
		return record, fmt.Errorf("transaction cannot start: %w", err)
	}

	err = tx.QueryRowContext(ctx, query+" RETURNING *", args...).Scan(&record)
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return record, fmt.Errorf("transaction cannot execute query and rollback: %w, %w", err, errRollback)
		}
		return record, fmt.Errorf("transaction cannot execute query: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return record, fmt.Errorf("transaction cannot commit: %w", err)
	}

	return record, nil
}

func (d *Storage[T]) ExecuteInsert(ctx context.Context, query string, args ...interface{}) (T, error) {
	return d.executeQueryRow(ctx, query, args...)
}

func (d *Storage[T]) ExecuteUpdate(ctx context.Context, query string, args ...interface{}) (T, error) {
	return d.executeQueryRow(ctx, query, args...)
}

func (d *Storage[T]) ExecuteDelete(ctx context.Context, query string, args ...interface{}) error {
	tx, err := d.Begin()
	if err != nil {
		return fmt.Errorf("transaction cannot start: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return fmt.Errorf("transaction cannot execute DELETE and rollback: %w, %w", err, errRollback)
		}
		return fmt.Errorf("transaction cannot execute DELETE: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("transaction cannot commit: %w", err)
	}

	return nil
}
