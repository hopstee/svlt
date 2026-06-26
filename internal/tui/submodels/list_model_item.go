package submodels

import (
	"fmt"

	"charm.land/bubbles/v2/list"
	"github.com/hopstee/svlt/internal/storage"
)

type listModelItem struct {
	storage.Connection
}

func (i listModelItem) Title() string {
	return i.Label
}

func (i listModelItem) Description() string {
	return fmt.Sprintf("%s:%s", i.Host, i.Port)
}

func (i listModelItem) FilterValue() string {
	return fmt.Sprintf("%s:%s", i.Host, i.Label)
}

func connectionToItem(conns []storage.Connection) []list.Item {
	items := make([]list.Item, len(conns))
	for i, c := range conns {
		items[i] = listModelItem{Connection: c}
	}
	return items
}
