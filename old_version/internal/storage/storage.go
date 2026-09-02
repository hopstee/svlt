package storage

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func openOrCreateDB(dbFile string) (*sql.DB, error) {
	return sql.Open("sqlite", dbFile)
}

func createTablesIfNotExsits(ctx context.Context, db *sql.DB) error {
	var err error

	tx, err := db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("Failed begin transaction: %w", err)
	}

	for _, q := range stmts {
		_, err = tx.ExecContext(ctx, q.SQL)
		if err != nil {
			return fmt.Errorf("Failed to create '%s' table: %w", q.Name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("Failed commit transaction: %w", err)
	}

	return nil
}

func NewStorage(dbFile string) (*Storage, error) {
	db, err := openOrCreateDB(dbFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to open or create database: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Init(ctx context.Context) error {
	if err := createTablesIfNotExsits(ctx, s.db); err != nil {
		return fmt.Errorf("Failed to create tables in database: %w", err)
	}
	return nil
}

func (s *Storage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
