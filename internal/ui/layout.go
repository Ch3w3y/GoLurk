package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Layout manages the arrangement of chat panels based on screen size
type Layout struct {
	width           int
	height          int
	maxVisiblePanes int
	minPaneWidth    int
}

// NewLayout creates a new layout manager
func NewLayout() *Layout {
	return &Layout{
		minPaneWidth: 40, // Minimum width for a chat panel to be readable
	}
}

// UpdateSize updates the layout dimensions
func (l *Layout) UpdateSize(width, height int) {
	l.width = width
	l.height = height

	// Calculate how many panes can fit based on the minimum width
	l.maxVisiblePanes = max(1, l.width/l.minPaneWidth)
}

// GetVisiblePanes returns the number of panes that can be displayed
func (l *Layout) GetVisiblePanes() int {
	return l.maxVisiblePanes
}

// ArrangePanes arranges chat panel views in a horizontal layout
func (l *Layout) ArrangePanes(paneViews []string) string {
	// Determine number of panes to show (limited by available panes and max visible)
	numVisible := min(len(paneViews), l.maxVisiblePanes)
	if numVisible == 0 {
		return "No chat panels available"
	}

	// Calculate width for each panel
	paneWidth := l.width / numVisible

	// Style each panel with the calculated width
	styledPanes := make([]string, numVisible)
	for i := 0; i < numVisible; i++ {
		styledPanes[i] = lipgloss.NewStyle().
			Width(paneWidth).
			Height(l.height).
			Render(paneViews[i])
	}

	// Join all styled panes horizontally
	return lipgloss.JoinHorizontal(lipgloss.Top, styledPanes...)
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
