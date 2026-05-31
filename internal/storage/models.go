package storage

import "time"

type AuthMethod uint

const (
	PasswordMethod AuthMethod = iota
	PassphraseMethod
	AgentMethod
)

type Connection struct {
	ID       string    `json:"id"`
	Label    string    `json:"label"`
	Group    string    `json:"group"`
	Tags     []string  `json:"tags"`
	LastUsed time.Time `json:"last_used"`
	IsActive bool      `json:"is_active"`
	IsPinned bool      `json:"is_pinned"`

	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`

	AuthMethod AuthMethod `json:"auth_type"`
	KeyPath    string     `json:"key_path,omitempty"`
}

type UpsertConnectionDto struct {
	Label string   `json:"label"`
	Group string   `json:"group"`
	Tags  []string `json:"tags"`

	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`

	AuthMethod AuthMethod `json:"auth_type"`
	KeyPath    string     `json:"key_path,omitempty"`
}
