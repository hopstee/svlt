package storage

import (
	"fmt"
	"time"

	"go.etcd.io/bbolt"
)

const CONNECTIONS_BUCKET = "ssh_connections"

type Storage struct {
	db *bbolt.DB
}

func openOrCreateDB(storePath string) (*bbolt.DB, error) {
	return bbolt.Open(storePath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
}

func createBucketsIfNotExsits(db *bbolt.DB) error {
	return db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(CONNECTIONS_BUCKET))
		if err != nil {
			return err
		}
		return nil
	})
}

func NewStorage(storePath string) (*Storage, error) {
	db, err := openOrCreateDB(storePath)
	if err != nil {
		return nil, fmt.Errorf("Failed to open or create database: %v", err)
	}

	if err := createBucketsIfNotExsits(db); err != nil {
		return nil, fmt.Errorf("Failed to create buckets in database: %v", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}
