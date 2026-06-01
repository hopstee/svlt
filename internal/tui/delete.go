package tui

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
	return tea.NewView("Delete view")
}
