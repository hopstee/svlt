package keyring

import (
	"errors"
	"fmt"

	"github.com/zalando/go-keyring"
)

var (
	ErrSecreteNotFound  = errors.New("Secrete not found in keyring")
	ErrKeyringOperation = errors.New("Keyring operation failed")
)

type Keyring struct {
	appName string
}

func NewKeyring(appName string) *Keyring {
	return &Keyring{appName: appName}
}

func (k *Keyring) Set(label, secret string) error {
	if label == "" || secret == "" {
		return fmt.Errorf("%w: label and secret cannot be empty", ErrKeyringOperation)
	}

	err := keyring.Set(k.appName, label, secret)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrKeyringOperation, err)
	}

	return nil
}

func (k *Keyring) Get(label string) (string, error) {
	if label == "" {
		return "", fmt.Errorf("%w: label cannot be empty", ErrKeyringOperation)
	}

	secret, err := keyring.Get(k.appName, label)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return "", ErrSecreteNotFound
		}
		return "", fmt.Errorf("%w: %v", ErrKeyringOperation, err)
	}

	return secret, nil
}

func (k *Keyring) Delete(label string) error {
	if label == "" {
		return fmt.Errorf("%w: label cannot be empty", ErrKeyringOperation)
	}

	err := keyring.Delete(k.appName, label)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("%w: %v", ErrKeyringOperation, err)
	}

	return nil
}
