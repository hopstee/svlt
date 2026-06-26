package cmds

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/hopstee/svlt/internal/keyring"
	"github.com/hopstee/svlt/internal/storage"
	"github.com/hopstee/svlt/internal/tui/root"
	"github.com/spf13/cobra"
)

var (
	store      *storage.Storage
	appKeyring *keyring.Keyring
)

func InitCLI(s *storage.Storage, k *keyring.Keyring) {
	store = s
	appKeyring = k
}

var RootCmd = &cobra.Command{
	Use:   "svlt",
	Short: "SSH connections manager",
	RunE: func(cmd *cobra.Command, args []string) error {
		connections, err := store.GetConns(cmd.Context())
		if err != nil {
			return fmt.Errorf("Failed to load connections: %v", err)
		}

		tui := root.NewModel(cmd.Context(), connections, appKeyring, store)
		tuiProgram := tea.NewProgram(tui)
		if _, err := tuiProgram.Run(); err != nil {
			return fmt.Errorf("TUI terminated with error: %v", err)
		}

		return nil
	},
}
