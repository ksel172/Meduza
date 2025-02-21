package utils

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func WithTransactionTimeout(
	ctx context.Context,
	db *sql.DB,
	timeout int,
	txOpts sql.TxOptions,
	fn func(context.Context, *sql.Tx) error,
) error {
	tx, err := db.BeginTx(ctx, &txOpts)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx failed: %v, rollback failed: %v", err, rbErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}

func WithTransactionResultTimeout[T any](
	ctx context.Context,
	db *sql.DB,
	timeout int,
	txOpts sql.TxOptions,
	fn func(context.Context, *sql.Tx) (T, error),
) (T, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to start transaction: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	result, err := fn(ctx, tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return result, fmt.Errorf("tx failed: %v, rollback failed: %v", err, rbErr)
		}
		return result, err
	}

	if err := tx.Commit(); err != nil {
		return result, fmt.Errorf("failed to commit tx: %w", err)
	}

	return result, nil
}

// For simple db operations that make only one query
func WithTimeout(
	ctx context.Context,
	db *sql.DB,
	query string,
	timeout int,
	fn func(context.Context, *sql.Stmt) error,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	return fn(ctx, stmt)
}

// For simple db operations that make only one query
func WithResultTimeout[T any](
	ctx context.Context,
	db *sql.DB,
	query string,
	timeout int,
	fn func(context.Context, *sql.Stmt) (T, error),
) (T, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		var zero T
		return zero, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	return fn(ctx, stmt)
}
