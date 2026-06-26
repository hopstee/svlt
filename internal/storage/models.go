package storage

import "time"

type AuthMethod string

const (
	PasswordMethod   AuthMethod = "password"
	PassphraseMethod AuthMethod = "passphrase"
)

type Connection struct {
	ID       string     `json:"id"`
	Label    string     `json:"label"`
	FolderID *string    `json:"folder_id"`
	LastUsed *time.Time `json:"last_used"`
	IsActive bool       `json:"is_active"`
	IsPinned bool       `json:"is_pinned"`

	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`

	AuthMethod AuthMethod `json:"auth_type"`
	KeyPath    string     `json:"key_path,omitempty"`
}

type UpsertConnectionDto struct {
	Label    string  `json:"label"`
	FolderID *string `json:"folder_id"`

	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`

	AuthMethod AuthMethod `json:"auth_type"`
	KeyPath    string     `json:"key_path,omitempty"`
}
