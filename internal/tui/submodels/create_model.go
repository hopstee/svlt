package submodels

import tea "charm.land/bubbletea/v2"

type CreateModel struct{}

func NewCreateModel() *CreateModel {
	return &CreateModel{}
}

func (cm *CreateModel) Init() tea.Cmd {
	return nil
}

func (cm *CreateModel) Update(msg tea.Msg) (*CreateModel, tea.Cmd) {
	return cm, nil
}

func (cm *CreateModel) View() tea.View {
	v := tea.NewView("Create view")
	v.AltScreen = true
	return v
}
