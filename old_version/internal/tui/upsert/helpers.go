package upsert

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"charm.land/huh/v2"
	"github.com/hopstee/svlt/old_version/internal/storage"
)

func loadSSHConfigs() ([]huh.Option[string], error) {
	var sshConfigs []huh.Option[string]

	userDir, err := os.UserHomeDir()
	if err != nil {
		return sshConfigs, fmt.Errorf("Failed get user directory: %w", err)
	}

	sshDir := filepath.Join(userDir, ".ssh")
	entries, err := os.ReadDir(sshDir)
	if err != nil {
		return sshConfigs, fmt.Errorf("Failed to read configs from ssh directory")
	}

	for _, entry := range entries {
		entryName := entry.Name()
		if !entry.IsDir() && strings.HasSuffix(entryName, ".pub") {
			key := strings.TrimSuffix(entryName, ".pub")
			value := filepath.Join(sshDir, key)
			sshConfigs = append(sshConfigs, huh.NewOption(key, value))
		}
	}
	sshConfigs = append(sshConfigs, huh.NewOption("Generate new", GenerateNewOption))

	return sshConfigs, nil
}

func generateForm(values *storage.UpsertConnectionDto, sshConfigs []huh.Option[string]) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("label").
				Title("Connection Label In List").
				Accessor(huh.NewPointerAccessor(&values.Label)).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("Label cannot be empty")
					}
					return nil
				}),

			huh.NewInput().
				Key("host").
				Title("IP address/Host").
				Accessor(huh.NewPointerAccessor(&values.Host)).
				Validate(func(str string) error {
					if str == "" {
						return errors.New("IP address/host cannot be empty")
					}
					return nil
				}),

			huh.NewInput().
				Key("port").
				Title("Port (Default 22)").
				Accessor(huh.NewPointerAccessor(&values.Port)),

			huh.NewInput().
				Key("user").
				Title("User (Default 'root')").
				Accessor(huh.NewPointerAccessor(&values.User)),

			huh.NewSelect[storage.AuthMethod]().
				Key("auth_type").
				Options(huh.NewOptions(storage.PasswordMethod, storage.PassphraseMethod)...).
				Title("Auth method").
				Accessor(huh.NewPointerAccessor(&values.AuthMethod)),

			huh.NewSelect[string]().
				Key("auth_config").
				Title("SSH config").
				Accessor(huh.NewPointerAccessor(&values.KeyPath)).
				Options(sshConfigs...),
		),
	).WithShowErrors(true).WithShowHelp(true)
}
