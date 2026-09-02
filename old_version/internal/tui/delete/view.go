package delete

import tea "charm.land/bubbletea/v2"

func (m *Model) View() tea.View {
	v := tea.NewView(m.form.View())
	v.AltScreen = true
	return v
}
