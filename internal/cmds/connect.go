package cmds

import (
	"fmt"
	"log/slog"

	"github.com/hopstee/svlt/internal/storage"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect [name]",
	Short: "Fast SSH connection to server",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		store, err := storage.NewStorage(dataPath)
		if err != nil {
			return fmt.Errorf("Failed to open storage: %v", err)
		}

		conn, err := store.GetOneByName(name)
		if err != nil {
			return fmt.Errorf("Failed to get single connection: %v", err)
		}

		slog.Info("data loaded", slog.Any("connection", conn))
		// connection logic

		return nil
	},
}

func init() {
	RootCmd.AddCommand(connectCmd)
}
