package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	getSizes    = "SELECT * FROM sizes ORDER BY size ASC"
	deleteSizes = "DELETE FROM sizes"
	insertSizes = "INSERT INTO sizes (size) VALUES ($1)"
)

type DB struct {
	conn *pgxpool.Pool
}

func NewDB(ctx context.Context, url string) (*DB, error) {
	conn, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("unable to initialise database: %w", err)
	}
	if err = conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &DB{conn: conn}, nil
}

func (db *DB) Close() {
	db.conn.Close()
}

func (db *DB) GetPackSizes(ctx context.Context) ([]int, error) {
	rows, err := db.conn.Query(ctx, getSizes)
	if err != nil {
		return nil, fmt.Errorf("failed to get pack sizes: %w", err)
	}
	defer rows.Close()

	var sizes []int
	for rows.Next() {
		var size int
		if err := rows.Scan(&size); err != nil {
			return nil, fmt.Errorf("failed to scan pack size: %w", err)
		}
		sizes = append(sizes, size)
	}

	return sizes, nil
}

func (db *DB) StorePackSizes(ctx context.Context, sizes []int) error {
	tx, err := db.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = db.conn.Exec(ctx, deleteSizes)
	if err != nil {
		return fmt.Errorf("failed to delete pack sizes: %w", err)
	}

	for _, size := range sizes {
		_, err := db.conn.Exec(ctx, insertSizes, size)
		if err != nil {
			return fmt.Errorf("failed to insert pack size: %w", err)
		}
	}

	return tx.Commit(ctx)
}
