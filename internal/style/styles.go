package style

import "github.com/charmbracelet/lipgloss"

var (
	GreenTextStyle = lipgloss.NewStyle().
			Bold(false).
			Foreground(lipgloss.Color("#00ff00"))

	ErrorTextStyle = lipgloss.NewStyle().
			Bold(false).
			Foreground(lipgloss.Color("#cc0000"))

	WarningTextStyle = lipgloss.NewStyle().
				Bold(false).
				Foreground(lipgloss.Color("#eed202"))
)
