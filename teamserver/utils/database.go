package utils

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func WithTransactionTimeout(db *sql.DB, ctx context.Context, timeout int, fn func(context.Context, *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
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

func WithTransactionResultTimeout[T any](db *sql.DB, ctx context.Context, timeout int, fn func(context.Context, *sql.Tx) (T, error)) (T, error) {
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

func WithTimeout(ctx context.Context, timeout int, fn func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()
	return fn(ctx)
}

func WithResultTimeout[T any](ctx context.Context, timeout int, fn func(context.Context) (T, error)) (T, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()
	return fn(ctx)
}
