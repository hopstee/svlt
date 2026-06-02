package submodels

import tea "charm.land/bubbletea/v2"

type DeleteModel struct{}

func NewDeleteModel() *DeleteModel {
	return &DeleteModel{}
}

func (dm *DeleteModel) Init() tea.Cmd {
	return nil
}

func (dm *DeleteModel) Update(msg tea.Msg) (*DeleteModel, tea.Cmd) {
	return dm, nil
}

func (dm *DeleteModel) View() tea.View {
	v := tea.NewView("Delete view")
	v.AltScreen = true
	return v
}
