package delete

import (
	"charm.land/huh/v2"
)

func generateForm(confirm *bool) *huh.Form {
	return huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Are you sure?").
				Affirmative("Yes!").
				Negative("No.").
				Value(confirm),
		),
	).WithShowErrors(true).WithShowHelp(true)
}
