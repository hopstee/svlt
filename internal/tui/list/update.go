package list

import (
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/hopstee/svlt/internal/storage"
	"github.com/hopstee/svlt/internal/tui/cmds"
	"github.com/hopstee/svlt/internal/tui/styles"
)

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	if m, c := m.handleSpecialCommands(msg); m != nil {
		return m, c
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m, c := m.handleKeyMsg(msg)
		if m != nil {
			return m, c
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m *Model) handleSpecialCommands(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		m.darkBG = msg.IsDark()
		m.updateAppStyles()
		return m, nil
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.updateAppStyles()
		return m, nil
	}
	return nil, nil
}

func (m *Model) updateAppStyles() {
	h, v := m.styles.App.GetFrameSize()
	m.list.SetSize(m.width-h, m.height-v)

	m.styles = styles.AppStyles(m.darkBG)
	m.list.Styles.Title = m.styles.Title
}

func (m *Model) handleKeyMsg(msg tea.KeyMsg) (*Model, tea.Cmd) {
	if m.list.FilterState() != list.Filtering {
		switch msg.String() {
		case "a":
			return m, cmds.OpenCreate
		case "e":
			if selected, ok := m.list.SelectedItem().(listModelItem); ok {
				return m, cmds.OpenEdit(selected.Connection)
			}
		case "d":
			if selected, ok := m.list.SelectedItem().(listModelItem); ok {
				return m, cmds.OpenDelete(selected.Connection.ID)
			}
		}
	}

	return nil, nil
}

func (m *Model) UpdateConnections(conns []storage.Connection) *Model {
	m.list.SetItems(connectionToItem(conns))

	m.normalizeCursor()

	return m
}
