package domain

import (
	"errors"

	"github.com/google/uuid"
)

type AuthType string

const (
	PasswordAuth AuthType = "password"
	PasskeyAuth  AuthType = "passkey"
)

var (
	ErrEmptyHost = errors.New("connection: host cannot be empty")
	ErrWrongPort = errors.New("connection: port should be valid number")
	ErrEmptyUser = errors.New("connection: user cannot be empty")
	ErrEmptyName = errors.New("connection: name cannot be empty")
)

type Connection struct {
	ID   uuid.UUID
	Name string
	Tag  string

	Host string
	Port uint8
	User string

	AuthType AuthType
}

func NewConnection(name, tag, host, user string, port uint8, authType AuthType) (Connection, error) {
	c := Connection{
		ID:       uuid.New(),
		Name:     name,
		Tag:      tag,
		Host:     host,
		Port:     port,
		User:     user,
		AuthType: authType,
	}

	if err := c.validate(); err != nil {
		return Connection{}, err
	}

	return c, nil
}

func (c Connection) validate() error {
	if c.Host == "" {
		return ErrEmptyHost
	}
	if c.Port == 0 {
		return ErrWrongPort
	}
	if c.User == "" {
		return ErrEmptyUser
	}
	if c.Name == "" {
		return ErrEmptyName
	}
	return nil
}
