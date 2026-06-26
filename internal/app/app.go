package app

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hopstee/svlt/internal/cmds"
	"github.com/hopstee/svlt/internal/keyring"
	"github.com/hopstee/svlt/internal/storage"
)

const (
	appName       = "svlt"
	appHiddenName = "." + appName
	appDbName     = appName + ".db"
)

type App struct {
	ctx    context.Context
	cancel context.CancelFunc
}

func InitApp(rootCtx context.Context) *App {
	ctx, cancel := context.WithCancel(rootCtx)
	return &App{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (a *App) Run() error {
	appDataPath, err := a.initAppDataPath()
	if err != nil {
		return fmt.Errorf("Failed init app data path: %v", err)
	}

	keystore := keyring.NewKeyring(appName)
	store, err := storage.NewStorage(appDataPath)
	if err != nil {
		return fmt.Errorf("Failed to create storage instance: %v", err)
	}
	if err := store.Init(a.ctx); err != nil {
		return fmt.Errorf("Failed init tables: %v", err)
	}
	cmds.InitCLI(store, keystore)

	return cmds.RootCmd.ExecuteContext(a.ctx)
}

func (a *App) Stop() error {
	a.cancel()
	return nil
}

func (a *App) initAppDataPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Failed get user directory: %w", err)
	}

	appDir := filepath.Join(homeDir, appHiddenName)
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", fmt.Errorf("Failed open or create app data directory")
	}

	return filepath.Join(appDir, appDbName), nil
}
