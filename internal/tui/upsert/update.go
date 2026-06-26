package upsert

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"github.com/hopstee/svlt/internal/storage"
	"github.com/hopstee/svlt/internal/tui/cmds"
)

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m, c := m.handleKeyMsg(msg)
		if m != nil {
			return m, c
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		return m.processComplete()
	}

	return m, cmd
}

func (m *Model) handleKeyMsg(msg tea.KeyMsg) (*Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		m.cancel()
		return m, cmds.BackToList
	}

	return nil, nil
}

func (m *Model) processComplete() (*Model, tea.Cmd) {
	m.setFieldsDefaults()

	connections, err := m.mutateAndReturnNewItemsList()

	if err != nil {
		m.dbErr = err
		m.form.State = huh.StateNormal
		return m, m.form.Init()
	}

	return m, cmds.RefreshList(connections)
}

func (m *Model) setFieldsDefaults() {
	if m.currentData.Port == "" {
		m.currentData.Port = "22"
	}
	if m.currentData.User == "" {
		m.currentData.User = "root"
	}
}

func (m *Model) mutateAndReturnNewItemsList() ([]storage.Connection, error) {
	var err error
	switch m.mode {
	case CreateMode:
		err = m.store.AddConn(m.ctx, m.currentData)
	case UpdateMode:
		err = m.store.Update(m.ctx, *m.oldID, m.currentData)
	}

	if err != nil {
		return nil, err
	}

	return m.store.GetConns(m.ctx)
}
