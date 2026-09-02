package list

import tea "charm.land/bubbletea/v2"

func (m *Model) View() tea.View {
	v := tea.NewView(m.styles.App.Render(m.list.View()))
	v.AltScreen = true
	return v
}
