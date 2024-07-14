package cli

import "github.com/charmbracelet/lipgloss"

var (
	appStyle          = lipgloss.NewStyle().Padding(1, 2)
	textStyle         = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemTitle = lipgloss.NewStyle().
				Border(lipgloss.Border{Left: "➤"}, false, false, false, true).
				BorderForeground(lipgloss.Color("#00FF00")).
				Foreground(lipgloss.Color("#00FF00")).
				PaddingLeft(1).
				Bold(true)
	selectedItemDesc = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#558855")).
				PaddingLeft(2)
	filterMatch = lipgloss.NewStyle().
			Background(lipgloss.Color("6"))
)
