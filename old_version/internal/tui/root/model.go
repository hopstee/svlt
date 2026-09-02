package root

import (
	"context"

	tea "charm.land/bubbletea/v2"
	"github.com/hopstee/svlt/old_version/internal/keyring"
	"github.com/hopstee/svlt/old_version/internal/storage"
	"github.com/hopstee/svlt/old_version/internal/tui/delete"
	"github.com/hopstee/svlt/old_version/internal/tui/list"
	"github.com/hopstee/svlt/old_version/internal/tui/upsert"
)

type RootState int

const (
	StateList RootState = iota
	StateCreate
	StateEdit
	StateDelete
)

type Model struct {
	ctx    context.Context
	cancel context.CancelFunc

	state RootState
	kr    *keyring.Keyring
	store *storage.Storage

	listModel   *list.Model
	upsertModel *upsert.Model
	deleteModel *delete.Model
}

func NewModel(
	rootCtx context.Context,
	conns []storage.Connection,
	kr *keyring.Keyring,
	st *storage.Storage,
) *Model {
	ctx, cancel := context.WithCancel(rootCtx)
	return &Model{
		ctx:    ctx,
		cancel: cancel,

		state: StateList,
		kr:    kr,
		store: st,

		listModel: list.NewModel(conns),
	}
}

func (m *Model) Init() tea.Cmd {
	return m.listModel.Init()
}
