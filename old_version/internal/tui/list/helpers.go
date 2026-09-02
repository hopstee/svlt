package list

import (
	"charm.land/bubbles/v2/list"
	"github.com/hopstee/svlt/old_version/internal/storage"
)

func connectionToItem(conns []storage.Connection) []list.Item {
	items := make([]list.Item, len(conns))
	for i, c := range conns {
		items[i] = listModelItem{Connection: c}
	}
	return items
}

func (m *Model) normalizeCursor() {
	cursor := m.list.Cursor()
	items := m.list.Items()

	if cursor > len(items) {
		m.list.Select(max(0, len(items)-1))
	}
}
