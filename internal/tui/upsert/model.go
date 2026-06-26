package upsert

import (
	"context"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/hopstee/svlt/internal/storage"
	"github.com/hopstee/svlt/internal/tui/styles"
)

type UpsertModelMode int

const (
	CreateMode UpsertModelMode = iota
	UpdateMode
)

const GenerateNewOption = "generate_new"

type Model struct {
	ctx    context.Context
	cancel context.CancelFunc

	form        *huh.Form
	currentData *storage.UpsertConnectionDto
	store       *storage.Storage

	mode  UpsertModelMode
	oldID *string

	dbErr error

	darkBG bool
	width  int
	height int
	styles styles.Styles
}

func NewModel(
	rootCtx context.Context,
	store *storage.Storage,
	connData storage.UpsertConnectionDto,
	connID *string,
	mode UpsertModelMode,
) *Model {
	sshConfigs, _ := loadSSHConfigs()
	ctx, cancel := context.WithCancel(rootCtx)
	return &Model{
		ctx:         ctx,
		cancel:      cancel,
		form:        generateForm(&connData, sshConfigs),
		currentData: &connData,
		store:       store,
		mode:        mode,
		oldID:       connID,
	}
}

func (m *Model) Init() tea.Cmd {
	return m.form.Init()
}
