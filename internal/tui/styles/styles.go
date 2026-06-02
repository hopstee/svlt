package styles

import (
	"charm.land/lipgloss/v2"
)

type Styles struct {
	App           lipgloss.Style
	Title         lipgloss.Style
	StatusMessage lipgloss.Style
}

func AppStyles(isDark bool) Styles {
	lightDark := lipgloss.LightDark(isDark)
	return Styles{
		App: lipgloss.NewStyle().Padding(0, 1),
		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1),
		StatusMessage: lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color("#04b575"), lipgloss.Color("#04b575"))),
	}
}
