package submodels

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/hopstee/svlt/internal/storage"
	"github.com/hopstee/svlt/internal/tui/cmds"
	"github.com/hopstee/svlt/internal/tui/styles"
)

type ListModel struct {
	list list.Model

	darkBG bool
	width  int
	height int
	styles styles.Styles
}

func NewListModel(conns []storage.Connection) *ListModel {
	items := connectionToItem(conns)
	// TODO: create custom delegate
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "SSH Connections"

	return &ListModel{
		list: l,
	}
}

func (lm *ListModel) Init() tea.Cmd {
	return textinput.Blink
}

func (lm *ListModel) UpdateConnections(conns []storage.Connection) *ListModel {
	lm.list.SetItems(connectionToItem(conns))
	return lm
}

func (lm *ListModel) Update(msg tea.Msg) (*ListModel, tea.Cmd) {
	if m, c := lm.handleSpecialCommands(msg); m != nil {
		return m, c
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		m, c := lm.handleKeyMsg(msg)
		if m != nil {
			return m, c
		}
	}

	var cmd tea.Cmd
	lm.list, cmd = lm.list.Update(msg)

	return lm, cmd
}

func (lm *ListModel) handleSpecialCommands(msg tea.Msg) (*ListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.BackgroundColorMsg:
		lm.darkBG = msg.IsDark()
		lm.updateAppStyles()
		return lm, nil
	case tea.WindowSizeMsg:
		lm.width, lm.height = msg.Width, msg.Height
		lm.updateAppStyles()
		return lm, nil
	}
	return nil, nil
}

func (lm *ListModel) updateAppStyles() {
	h, v := lm.styles.App.GetFrameSize()
	lm.list.SetSize(lm.width-h, lm.height-v)

	lm.styles = styles.AppStyles(lm.darkBG)
	lm.list.Styles.Title = lm.styles.Title
}

func (lm *ListModel) handleKeyMsg(msg tea.KeyMsg) (*ListModel, tea.Cmd) {
	if lm.list.FilterState() != list.Filtering {
		switch msg.String() {
		case "a":
			return lm, cmds.OpenCreate
		case "e":
			if selected, ok := lm.list.SelectedItem().(listModelItem); ok {
				return lm, cmds.OpenEdit(selected.Connection)
			}
		}
	}

	return nil, nil
}

func (lm *ListModel) View() tea.View {
	v := tea.NewView(lm.styles.App.Render(lm.list.View()))
	v.AltScreen = true
	return v
}
