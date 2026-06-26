package submodels

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/hopstee/svlt/internal/storage"
	"github.com/hopstee/svlt/internal/tui/cmds"
)

type UpsertModelMode int

const (
	UpsertCreateMode UpsertModelMode = iota
	UpsertUpdateMode
)

const GenerateNewOption = "generate_new"

type UpsertModel struct {
	ctx    context.Context
	cancel context.CancelFunc

	form        *huh.Form
	currentData *storage.UpsertConnectionDto
	store       *storage.Storage

	mode  UpsertModelMode
	oldID *string

	dbErr error
}

func NewUpsertModel(
	rootCtx context.Context,
	store *storage.Storage,
	connData storage.UpsertConnectionDto,
	connID *string,
	mode UpsertModelMode,
) *UpsertModel {
	sshConfigs, _ := loadSSHConfigs()
	ctx, cancel := context.WithCancel(rootCtx)
	return &UpsertModel{
		ctx:         ctx,
		cancel:      cancel,
		form:        generateForm(&connData, sshConfigs),
		currentData: &connData,
		store:       store,
		mode:        mode,
		oldID:       connID,
	}
}

func (um *UpsertModel) Init() tea.Cmd {
	return um.form.Init()
}

func (um *UpsertModel) Update(msg tea.Msg) (*UpsertModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m, c := um.handleKeyMsg(msg)
		if m != nil {
			return m, c
		}
	}

	form, cmd := um.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		um.form = f
	}

	if um.form.State == huh.StateCompleted {
		if um.currentData.Port == "" {
			um.currentData.Port = "22"
		}
		if um.currentData.User == "" {
			um.currentData.User = "root"
		}

		var err error
		switch um.mode {
		case UpsertCreateMode:
			err = um.store.AddConn(um.ctx, um.currentData)
		case UpsertUpdateMode:
			err = um.store.Update(um.ctx, *um.oldID, um.currentData)
		}

		connections, err := um.store.GetConns(um.ctx)

		if err != nil {
			um.dbErr = err
			um.form.State = huh.StateNormal
			return um, um.form.Init()
		}

		return um, cmds.RefreshList(connections)
	}

	return um, cmd
}

func (um *UpsertModel) handleKeyMsg(msg tea.KeyMsg) (*UpsertModel, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		um.cancel()
		return um, cmds.BackToList
	}

	return nil, nil
}

func (um *UpsertModel) View() tea.View {
	v := tea.NewView(um.form.View())
	v.AltScreen = true
	return v
}

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
