package root

import (
	tea "charm.land/bubbletea/v2"
	"github.com/hopstee/svlt/internal/storage"
	"github.com/hopstee/svlt/internal/tui/cmds"
	"github.com/hopstee/svlt/internal/tui/delete"
	"github.com/hopstee/svlt/internal/tui/upsert"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m, c := m.handleCustomCommands(msg); m != nil {
		return m, c
	}

	cmds := m.delegateEventToActiveSubmodel(msg)

	return m, tea.Batch(cmds...)
}

func (m *Model) handleCustomCommands(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case cmds.MsgBackToList:
		return m.processBackToListMsg()
	case cmds.MsgOpenCreate:
		return m.processOpenCreateMsg()
	case cmds.MsgOpenEdit:
		return m.processOpenEditMsg(msg)
	case cmds.MsgOpenDelete:
		return m.processOpenDeleteMsg(msg)
	case cmds.MsgRefreshList:
		return m.processRefreshListMsg(msg)
	}

	return nil, nil
}

func (m *Model) processBackToListMsg() (*Model, tea.Cmd) {
	m.state = StateList
	return m, nil
}

func (m *Model) processOpenCreateMsg() (*Model, tea.Cmd) {
	m.upsertModel = upsert.NewModel(
		m.ctx,
		m.store,
		storage.UpsertConnectionDto{},
		nil,
		upsert.CreateMode,
	)
	m.state = StateCreate
	return m, m.upsertModel.Init()
}

func (m *Model) processOpenEditMsg(msg cmds.MsgOpenEdit) (*Model, tea.Cmd) {
	upsertConnectionDto := storage.ConnectionToUpsertConnectionDto(msg.Conn)
	m.upsertModel = upsert.NewModel(
		m.ctx,
		m.store,
		upsertConnectionDto,
		&msg.Conn.ID,
		upsert.UpdateMode,
	)
	m.state = StateEdit
	return m, m.upsertModel.Init()
}

func (m *Model) processOpenDeleteMsg(msg cmds.MsgOpenDelete) (*Model, tea.Cmd) {
	m.deleteModel = delete.NewModel(m.ctx, msg.ConnID, m.store)
	m.state = StateDelete
	return m, m.deleteModel.Init()
}

func (m *Model) processRefreshListMsg(msg cmds.MsgRefreshList) (*Model, tea.Cmd) {
	m.listModel = m.listModel.UpdateConnections(msg.Conns)
	m.state = StateList
	return m, nil
}

func (m *Model) delegateEventToActiveSubmodel(msg tea.Msg) []tea.Cmd {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch m.state {
	case StateList:
		m.listModel, cmd = m.listModel.Update(msg)
		cmds = append(cmds, cmd)
	case StateCreate, StateEdit:
		m.upsertModel, cmd = m.upsertModel.Update(msg)
		cmds = append(cmds, cmd)
	case StateDelete:
		m.deleteModel, cmd = m.deleteModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return cmds
}
