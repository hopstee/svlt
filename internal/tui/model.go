package tui

import (
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
	state       RootState
	connections []storage.Connection
	kr          *keyring.Keyring
	dataPath    string

	listModel   *submodels.ListModel
	createModel *submodels.CreateModel
	editModel   *submodels.EditModel
	deleteModel *submodels.DeleteModel
}

func NewRootModel(conns []storage.Connection, kr *keyring.Keyring, dataPath string) *RootModel {
	return &RootModel{
		state:       StateList,
		connections: conns,
		kr:          kr,
		dataPath:    dataPath,

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
		rm.createModel = submodels.NewCreateModel()
		rm.state = StateCreate
		return rm, rm.createModel.Init()
	case cmds.MsgOpenEdit:
		rm.editModel = submodels.NewEditModel()
		rm.state = StateEdit
		return rm, rm.editModel.Init()
	case cmds.MsgOpenDelete:
		rm.deleteModel = submodels.NewDeleteModel()
		rm.state = StateDelete
		return rm, rm.deleteModel.Init()
	case cmds.MsgRefreshList:
		rm.connections = msg.Conns
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
		return rm.createModel.View()
	case StateEdit:
		return rm.editModel.View()
	case StateDelete:
		return rm.deleteModel.View()
	}
	return tea.NewView("Unknown state")
}
