package tui

import tea "charm.land/bubbletea/v2"

type EditModel struct{}

func NewEditModel() *EditModel {
	return &EditModel{}
}

func (em *EditModel) Init() tea.Cmd {
	return nil
}

func (em *EditModel) Update(msg tea.Msg) (*EditModel, tea.Cmd) {
	return em, nil
}

func (em *EditModel) View() tea.View {
	return tea.NewView("Edit view")
}
