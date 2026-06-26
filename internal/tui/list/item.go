package list

import (
	"fmt"

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
