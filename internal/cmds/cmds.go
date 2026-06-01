package cmds

import (
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/hopstee/svlt/internal/keyring"
	"github.com/hopstee/svlt/internal/storage"
	"github.com/hopstee/svlt/internal/tui"
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
		connections, err := storage.Execute(dataPath, func(store *storage.Storage) ([]storage.Connection, error) {
			return store.GetConns()
		})
		if err != nil {
			return fmt.Errorf("Failed to load connections: %v", err)
		}

		connections = mockConnections()

		tui := tui.NewRootModel(connections, appKeyring, dataPath)
		tuiProgram := tea.NewProgram(tui)
		if _, err := tuiProgram.Run(); err != nil {
			return fmt.Errorf("TUI terminated with error: %v", err)
		}

		return nil
	},
}

func mockConnections() []storage.Connection {
	return []storage.Connection{
		{
			ID:       "fd7sa89-fdsaf9-fd7s89af",
			Label:    "Prod DB",
			Group:    "",
			Tags:     []string{},
			LastUsed: time.Now().AddDate(0, 0, 1),
			IsActive: true,
			IsPinned: false,

			Host: "0.0.0.0",
			Port: 22,
			User: "root",

			AuthMethod: storage.PassphraseMethod,
			KeyPath:    "~/.ssh/id_rsa",
		},
		{
			ID:       "fd7sa89-fds2f9-fd7gdfs9af",
			Label:    "CB Bot V1 Sasha",
			Group:    "",
			Tags:     []string{},
			LastUsed: time.Now().AddDate(0, 0, 2),
			IsActive: true,
			IsPinned: false,

			Host: "0.0.0.1",
			Port: 22,
			User: "root",

			AuthMethod: storage.PassphraseMethod,
			KeyPath:    "~/.ssh/id_rsa",
		},
	}
}
