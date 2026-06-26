package root

import tea "charm.land/bubbletea/v2"

func (m *Model) View() tea.View {
	switch m.state {
	case StateList:
		return m.listModel.View()
	case StateCreate, StateEdit:
		return m.upsertModel.View()
	case StateDelete:
		return m.deleteModel.View()
	}
	return tea.NewView("Unknown state")
}
