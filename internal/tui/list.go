package tui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/hopstee/svlt/internal/storage"
)

type ListModel struct {
	allConns    []storage.Connection
	filterd     []storage.Connection
	searchInput textinput.Model
	cursor      int
}

func NewListModel(conns []storage.Connection) *ListModel {
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.Focus()
	return &ListModel{
		allConns:    conns,
		filterd:     conns,
		searchInput: ti,
	}
}

func (lm *ListModel) Init() tea.Cmd {
	return textinput.Blink
}

func (lm *ListModel) UpdateConnections(conns []storage.Connection) *ListModel {
	lm.allConns = conns
	lm.filterd = conns
	return lm
}

func (lm *ListModel) Update(msg tea.Msg) (*ListModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m, c := lm.handleKeyMsg(msg)
		if m != nil {
			return m, c
		}
	}

	var cmd tea.Cmd
	lm.searchInput, cmd = lm.searchInput.Update(msg)

	lm.applyFilter()

	return lm, cmd
}

func (lm *ListModel) handleKeyMsg(msg tea.KeyMsg) (*ListModel, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		return lm, tea.Quit
	case "up":
		if lm.cursor > 0 {
			lm.cursor--
		}
		return lm, nil
	case "down":
		if lm.cursor < len(lm.filterd)-1 {
			lm.cursor++
		}
		return lm, nil
	case "enter":
		// TODO: call ssh connect method, will be created later
		return lm, nil
	case "a":
		if !lm.searchInput.Focused() {
			return lm, openCreate
		}
	case "e":
		if !lm.searchInput.Focused() && len(lm.filterd) > 0 {
			return lm, openEdit(lm.filterd[lm.cursor])
		}
	case "d":
		if !lm.searchInput.Focused() && len(lm.filterd) > 0 {
			return lm, openDelete(lm.filterd[lm.cursor].Label)
		}
	}

	return nil, nil
}

func (lm *ListModel) applyFilter() {
	query := strings.ToLower(lm.searchInput.Value())
	if query == "" {
		lm.filterd = lm.allConns
		return
	}

	var res []storage.Connection
	for _, c := range lm.allConns {
		normLabel := strings.ToLower(strings.TrimSpace(c.Label))
		normHost := strings.ToLower(strings.TrimSpace(c.Host))
		if strings.Contains(normLabel, query) || strings.Contains(normHost, query) {
			res = append(res, c)
		}
	}
	lm.filterd = res
}

func (lm *ListModel) View() tea.View {
	var connectionsToRender = []string{"===SSH MANAGER==="}
	for i, c := range lm.filterd {
		prefix := "[ ]"
		if i == lm.cursor {
			prefix = "[x]"
		}
		connectionsToRender = append(connectionsToRender, fmt.Sprintf("%s %s", prefix, c.Label))
	}
	return tea.NewView(strings.Join(connectionsToRender, "\n"))
}
