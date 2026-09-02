package cmds

import (
	tea "charm.land/bubbletea/v2"
	"github.com/hopstee/svlt/old_version/internal/storage"
)

type MsgBackToList struct{}
type MsgOpenCreate struct{}
type MsgOpenEdit struct{ Conn storage.Connection }
type MsgOpenDelete struct{ ConnID string }
type MsgRefreshList struct{ Conns []storage.Connection }

func BackToList() tea.Msg {
	return MsgBackToList{}
}

func OpenCreate() tea.Msg {
	return MsgOpenCreate{}
}

func OpenEdit(conn storage.Connection) func() tea.Msg {
	return func() tea.Msg {
		return MsgOpenEdit{Conn: conn}
	}
}

func OpenDelete(connID string) func() tea.Msg {
	return func() tea.Msg {
		return MsgOpenDelete{ConnID: connID}
	}
}

func RefreshList(conns []storage.Connection) func() tea.Msg {
	return func() tea.Msg {
		return MsgRefreshList{Conns: conns}
	}
}
