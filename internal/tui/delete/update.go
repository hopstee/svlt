package delete

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
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
		if *m.confirm {
			err := m.store.DeleteConn(m.ctx, m.connID)

			connections, err := m.store.GetConns(m.ctx)

			if err != nil {
				m.form.State = huh.StateNormal
				return m, m.form.Init()
			}

			return m, cmds.RefreshList(connections)
		}
		return m, cmds.BackToList
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
