package list

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/hopstee/svlt/old_version/internal/storage"
	"github.com/hopstee/svlt/old_version/internal/tui/styles"
)

type Model struct {
	list list.Model

	darkBG bool
	width  int
	height int
	styles styles.Styles
}

func NewModel(conns []storage.Connection) *Model {
	items := connectionToItem(conns)
	// TODO: create custom delegate
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "SSH Connections"

	return &Model{
		list: l,
	}
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}
