package tui

import (
	"context"

	"github.com/hopstee/svlt/internal/keyring"
	"github.com/hopstee/svlt/internal/storage"
	"github.com/hopstee/svlt/internal/tui/cmds"
	"github.com/hopstee/svlt/internal/tui/submodels"

	tea "charm.land/bubbletea/v2"
)

type RootState int

const (
	StateList RootState = iota
	StateCreate
	StateEdit
	StateDelete
)

type RootModel struct {
	ctx    context.Context
	cancel context.CancelFunc

	state RootState
	kr    *keyring.Keyring
	store *storage.Storage

	listModel   *submodels.ListModel
	upsertModel *submodels.UpsertModel
	deleteModel *submodels.DeleteModel
}

func NewRootModel(
	rootCtx context.Context,
	conns []storage.Connection,
	kr *keyring.Keyring,
	st *storage.Storage,
) *RootModel {
	ctx, cancel := context.WithCancel(rootCtx)
	return &RootModel{
		ctx:    ctx,
		cancel: cancel,

		state: StateList,
		kr:    kr,
		store: st,

		listModel: submodels.NewListModel(conns),
	}
}

func (rm *RootModel) Init() tea.Cmd {
	return rm.listModel.Init()
}

func (rm *RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m, c := rm.handleCustomCommands(msg); m != nil {
		return m, c
	}

	cmds := rm.delegateEventToActiveSubmodel(msg)

	return rm, tea.Batch(cmds...)
}

func (rm *RootModel) handleCustomCommands(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case cmds.MsgBackToList:
		rm.state = StateList
		return rm, nil
	case cmds.MsgOpenCreate:
		rm.upsertModel = submodels.NewUpsertModel(
			rm.ctx,
			rm.store,
			storage.UpsertConnectionDto{},
			nil,
			submodels.UpsertCreateMode,
		)
		rm.state = StateCreate
		return rm, rm.upsertModel.Init()
	case cmds.MsgOpenEdit:
		upsertConnectionDto := storage.ConnectionToUpsertConnectionDto(msg.Conn)
		rm.upsertModel = submodels.NewUpsertModel(
			rm.ctx,
			rm.store,
			upsertConnectionDto,
			&msg.Conn.ID,
			submodels.UpsertUpdateMode,
		)
		rm.state = StateEdit
		return rm, rm.upsertModel.Init()
	case cmds.MsgOpenDelete:
		rm.deleteModel = submodels.NewDeleteModel()
		rm.state = StateDelete
		return rm, rm.deleteModel.Init()
	case cmds.MsgRefreshList:
		rm.listModel = rm.listModel.UpdateConnections(msg.Conns)
		rm.state = StateList
		return rm, nil
	}

	return nil, nil
}

func (rm *RootModel) delegateEventToActiveSubmodel(msg tea.Msg) []tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch rm.state {
	case StateList:
		rm.listModel, cmd = rm.listModel.Update(msg)
		cmds = append(cmds, cmd)
	case StateCreate, StateEdit:
		rm.upsertModel, cmd = rm.upsertModel.Update(msg)
		cmds = append(cmds, cmd)
	case StateDelete:
		rm.deleteModel, cmd = rm.deleteModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return cmds
}

func (rm *RootModel) View() tea.View {
	switch rm.state {
	case StateList:
		return rm.listModel.View()
	case StateCreate, StateEdit:
		return rm.upsertModel.View()
	case StateDelete:
		return rm.deleteModel.View()
	}
	return tea.NewView("Unknown state")
}
