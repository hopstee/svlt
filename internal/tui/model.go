package tui

import (
	"github.com/hopstee/svlt/internal/keyring"
	"github.com/hopstee/svlt/internal/storage"

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
	state       RootState
	connections []storage.Connection
	kr          *keyring.Keyring
	dataPath    string

	// TODO: create submodels
	listModel   *ListModel
	createModel *CreateModel
	editModel   *EditModel
	deleteModel *DeleteModel
}

func NewRootModel(conns []storage.Connection, kr *keyring.Keyring, dataPath string) *RootModel {
	return &RootModel{
		state:       StateList,
		connections: conns,
		kr:          kr,
		dataPath:    dataPath,

		// TODO: init submodels with data
		listModel: NewListModel(conns),
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
	case MsgBackToList:
		rm.state = StateList
		return rm, nil
	case MsgOpenCreate:
		rm.createModel = NewCreateModel()
		rm.state = StateCreate
		return rm, rm.createModel.Init()
	case MsgOpenEdit:
		rm.editModel = NewEditModel()
		rm.state = StateEdit
		return rm, rm.editModel.Init()
	case MsgOpenDelete:
		rm.deleteModel = NewDeleteModel()
		rm.state = StateDelete
		return rm, rm.deleteModel.Init()
	case MsgRefreshList:
		rm.connections = msg.conns
		rm.listModel = rm.listModel.UpdateConnections(msg.conns)
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
	case StateCreate:
		rm.createModel, cmd = rm.createModel.Update(msg)
		cmds = append(cmds, cmd)
	case StateEdit:
		rm.editModel, cmd = rm.editModel.Update(msg)
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
	case StateCreate:
	// TODO: return rm.createModel.View()
	case StateEdit:
	// TODO: return rm.editModel.View()
	case StateDelete:
		// TODO: return rm.deleteModel.View()
	}
	return tea.NewView("Unknown state")
}
