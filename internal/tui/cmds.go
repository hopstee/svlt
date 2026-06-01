package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/hopstee/svlt/internal/storage"
)

type MsgBackToList struct{}
type MsgOpenCreate struct{}
type MsgOpenEdit struct{ conn storage.Connection }
type MsgOpenDelete struct{ connName string }
type MsgRefreshList struct{ conns []storage.Connection }

func backToList() tea.Msg {
	return MsgBackToList{}
}

func openCreate() tea.Msg {
	return MsgOpenCreate{}
}

func openEdit(conn storage.Connection) func() tea.Msg {
	return func() tea.Msg {
		return MsgOpenEdit{conn: conn}
	}
}

func openDelete(connName string) func() tea.Msg {
	return func() tea.Msg {
		return MsgOpenDelete{connName: connName}
	}
}

func refreshList(conns []storage.Connection) func() tea.Msg {
	return func() tea.Msg {
		return MsgRefreshList{conns: conns}
	}
}
