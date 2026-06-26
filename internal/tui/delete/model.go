package delete

import (
	"context"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/hopstee/svlt/internal/storage"
)

type Model struct {
	ctx    context.Context
	cancel context.CancelFunc

	store *storage.Storage

	connID  string
	confirm *bool
	form    *huh.Form
}

func NewModel(
	rootCtx context.Context,
	connID string,
	store *storage.Storage,
) *Model {
	confirm := false
	ctx, cancel := context.WithCancel(rootCtx)
	return &Model{
		ctx:    ctx,
		cancel: cancel,

		store: store,

		connID:  connID,
		confirm: &confirm,
		form:    generateForm(&confirm),
	}
}

func (m *Model) Init() tea.Cmd {
	return m.form.Init()
}
