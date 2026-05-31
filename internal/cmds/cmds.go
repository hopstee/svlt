package cmds

import (
	"fmt"
	"log/slog"

	"github.com/hopstee/svlt/internal/keyring"
	"github.com/hopstee/svlt/internal/storage"
	"github.com/spf13/cobra"
)

var (
	dataPath   string
	appKeyring *keyring.Keyring
)

func InitCLI(dp string, k *keyring.Keyring) {
	dataPath = dp
	appKeyring = k
}

var RootCmd = &cobra.Command{
	Use:   "svlt",
	Short: "SSH connections manager",
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := storage.NewStorage(dataPath)
		if err != nil {
			return fmt.Errorf("Failed to open storage: %v", err)
		}

		connections, err := store.GetConns()
		_ = store.Close()
		if err != nil {
			return fmt.Errorf("Failed to load connections: %v", err)
		}

		slog.Info("data loaded", slog.Any("connections", connections))
		// use connections in tui

		return nil
	},
}
