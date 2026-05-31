package storage

import (
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

func (s *Storage) GetConns() ([]Connection, error) {
	var conns []Connection
	err := s.db.View(func(tx *bbolt.Tx) error {
		connsBucket := tx.Bucket([]byte(CONNECTIONS_BUCKET))
		if connsBucket == nil {
			return ErrConnectionBucketNotFound
		}

		return connsBucket.ForEach(func(k, v []byte) error {
			var c Connection
			if err := json.Unmarshal(v, &c); err != nil {
				return err
			}
			conns = append(conns, c)
			return nil
		})
	})
	return conns, err
}

func (s *Storage) GetOneByName(name string) (*Connection, error) {
	var conn *Connection
	err := s.db.View(func(tx *bbolt.Tx) error {
		connsBucket := tx.Bucket([]byte(CONNECTIONS_BUCKET))
		if connsBucket == nil {
			return ErrConnectionBucketNotFound
		}

		conn = s.getOneByName(connsBucket, name)
		if conn == nil {
			return ErrConnectionNotFound
		}
		return nil
	})
	return conn, err
}

func (s *Storage) getOneByName(bucket *bbolt.Bucket, name string) *Connection {
	key := []byte(s.getKey(name))
	v := bucket.Get(key)
	if v == nil {
		return nil
	}

	var conn Connection
	if err := json.Unmarshal(v, &conn); err != nil {
		return nil
	}
	return &conn
}

func (s *Storage) AddConn(newConn *UpsertConnectionDto) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		connsBucket := tx.Bucket([]byte(CONNECTIONS_BUCKET))
		if connsBucket == nil {
			return ErrConnectionBucketNotFound
		}

		existingConn := s.getOneByName(connsBucket, newConn.Label)
		if existingConn != nil {
			return ErrConnectionAlreadyExists
		}

		bConn, err := s.connectionToByte("", newConn)
		if err != nil {
			slog.Error(ErrConnectionToBytes.Error(), slog.Any("error", err))
			return ErrConnectionToBytes
		}

		return connsBucket.Put([]byte(s.getKey(newConn.Label)), bConn)
	})
}

func (s *Storage) Update(oldName string, updatedConn *UpsertConnectionDto) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		connsBucket := tx.Bucket([]byte(CONNECTIONS_BUCKET))
		if connsBucket == nil {
			return ErrConnectionBucketNotFound
		}

		existingConn := s.getOneByName(connsBucket, oldName)
		if existingConn == nil {
			return ErrConnectionNotFound
		}

		bConn, err := s.connectionToByte(existingConn.ID, updatedConn)
		if err != nil {
			slog.Error(ErrConnectionToBytes.Error(), slog.Any("error", err))
			return ErrConnectionToBytes
		}

		if s.getKey(oldName) != s.getKey(updatedConn.Label) {
			if s.getOneByName(connsBucket, updatedConn.Label) != nil {
				return ErrConnectionAlreadyExists
			}
			if err := connsBucket.Put([]byte(s.getKey(updatedConn.Label)), bConn); err != nil {
				slog.Error(ErrUpdateConnection.Error(), slog.Any("error", err))
				return ErrUpdateConnection
			}
			if err := connsBucket.Delete([]byte(s.getKey(oldName))); err != nil {
				return ErrFailedDelete
			}
			return nil
		}

		return connsBucket.Put([]byte(s.getKey(oldName)), bConn)
	})
}

func (s *Storage) connectionToByte(ID string, conn *UpsertConnectionDto) ([]byte, error) {
	if ID == "" {
		ID = uuid.New().String()
	}

	r := Connection{
		ID:         ID,
		Label:      conn.Label,
		Host:       conn.Host,
		Port:       conn.Port,
		User:       conn.User,
		AuthMethod: conn.AuthMethod,
		KeyPath:    conn.KeyPath,
	}

	v, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (s *Storage) DeleteConn(name string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		connsBucket := tx.Bucket([]byte(CONNECTIONS_BUCKET))
		if connsBucket == nil {
			return ErrConnectionBucketNotFound
		}

		existingConn := s.getOneByName(connsBucket, name)
		if existingConn == nil {
			return ErrConnectionNotFound
		}

		return connsBucket.Delete([]byte(s.getKey(name)))
	})
}

func (s *Storage) getKey(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}
