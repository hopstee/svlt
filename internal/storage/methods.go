package storage

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
)

func (s *Storage) GetConns(ctx context.Context) ([]Connection, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			id,
			label,
			folder_id,
			last_used,
			is_active,
			is_pinned,
			host,
			port,
			user,
			auth_method,
			key_path
		FROM connections
		ORDER BY label;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conns []Connection

	for rows.Next() {
		var c Connection

		err := rows.Scan(
			&c.ID,
			&c.Label,
			&c.FolderID,
			&c.LastUsed,
			&c.IsActive,
			&c.IsPinned,
			&c.Host,
			&c.Port,
			&c.User,
			&c.AuthMethod,
			&c.KeyPath,
		)
		if err != nil {
			return nil, err
		}

		conns = append(conns, c)
	}

	return conns, rows.Err()
}

func (s *Storage) GetOneByName(ctx context.Context, name string) (*Connection, error) {
	var c Connection

	err := s.db.QueryRowContext(ctx, `
		SELECT
			id,
			label,
			folder_id,
			last_used,
			is_active,
			is_pinned,
			host,
			port,
			user,
			auth_method,
			key_path
		FROM connections
		WHERE LOWER(label) = LOWER(?)
		LIMIT 1;
	`, strings.TrimSpace(name)).
		Scan(
			&c.ID,
			&c.Label,
			&c.FolderID,
			&c.LastUsed,
			&c.IsActive,
			&c.IsPinned,
			&c.Host,
			&c.Port,
			&c.User,
			&c.AuthMethod,
			&c.KeyPath,
		)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrConnectionNotFound
	}

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (s *Storage) AddConn(ctx context.Context, dto *UpsertConnectionDto) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO connections (
			id,
			label,
			host,
			port,
			user,
			auth_method,
			key_path
		)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		uuid.NewString(),
		dto.Label,
		dto.Host,
		dto.Port,
		dto.User,
		dto.AuthMethod,
		dto.KeyPath,
	)

	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return ErrConnectionAlreadyExists
		}
		return err
	}

	return nil
}

func (s *Storage) Update(ctx context.Context, id string, dto *UpsertConnectionDto) error {
	res, err := s.db.ExecContext(ctx, `
		UPDATE connections
		SET
			label = ?,
			host = ?,
			port = ?,
			user = ?,
			auth_method = ?,
			key_path = ?
		WHERE id = ?
	`,
		dto.Label,
		dto.Host,
		dto.Port,
		dto.User,
		dto.AuthMethod,
		dto.KeyPath,
		id,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrConnectionNotFound
	}

	return nil
}

func (s *Storage) DeleteConn(ctx context.Context, ID string) error {
	res, err := s.db.ExecContext(ctx, `
		DELETE FROM connections
		WHERE id = ?
	`, ID)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return ErrConnectionNotFound
	}

	return nil
}
