package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// DocStyle is the style for the entire app
	DocStyle = lipgloss.NewStyle().
			Margin(1, 2)

	// FocusedStyle is for elements that are focused
	FocusedStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62"))

	// BlurredStyle is for elements that are not focused
	BlurredStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))

	// HeaderStyle is for panel headers
	HeaderStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Align(lipgloss.Center)

	// UsernameStyle is for chat usernames
	UsernameStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("63")).
			Bold(true)

	// ActionStyle is for chat action messages
	ActionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("63")).
			Italic(true)

	// TabStyle is the style for tabs
	TabStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// ActiveTabStyle is for the active tab
	ActiveTabStyle = TabStyle.Copy().
			Foreground(lipgloss.Color("62")).
			Bold(true).
			Underline(true)

	// TabGap is the gap between tabs
	TabGap = lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render("│")
)
