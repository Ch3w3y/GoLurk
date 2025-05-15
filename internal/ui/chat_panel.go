package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ChatMessage represents a message in the chat
type ChatMessage struct {
	Username string
	Content  string
	IsAction bool
}

// ChatPanel represents a single Twitch chat panel
type ChatPanel struct {
	channel   string
	messages  []ChatMessage
	viewport  viewport.Model
	width     int
	height    int
	focused   bool
	style     lipgloss.Style
	maxBuffer int
}

// NewChatPanel creates a new chat panel for a specific channel
func NewChatPanel(channel string) *ChatPanel {
	vp := viewport.New(80, 20)
	vp.KeyMap = viewport.KeyMap{} // Disable default keybindings

	return &ChatPanel{
		channel:   channel,
		messages:  []ChatMessage{},
		viewport:  vp,
		maxBuffer: 500, // Maximum number of messages to keep in buffer
		style:     lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()),
	}
}

// SetSize updates the size of the chat panel
func (c *ChatPanel) SetSize(width, height int) {
	c.width = width
	c.height = height
	c.viewport.Width = width - 2   // Account for borders
	c.viewport.Height = height - 3 // Account for header and borders

	// Re-render content after resize
	c.renderContent()
}

// SetFocused sets the focus state of the panel
func (c *ChatPanel) SetFocused(focused bool) {
	c.focused = focused
	if focused {
		c.style = c.style.BorderForeground(lipgloss.Color("62"))
	} else {
		c.style = c.style.BorderForeground(lipgloss.Color("240"))
	}
}

// AddMessage adds a new message to the chat panel
func (c *ChatPanel) AddMessage(msg ChatMessage) {
	c.messages = append(c.messages, msg)

	// Trim buffer if needed
	if len(c.messages) > c.maxBuffer {
		c.messages = c.messages[len(c.messages)-c.maxBuffer:]
	}

	c.renderContent()
	c.viewport.GotoBottom() // Auto-scroll to bottom
}

// Update handles messages and user input
func (c *ChatPanel) Update(msg tea.Msg) (ChatPanel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if c.focused {
			switch msg.String() {
			case "up":
				c.viewport.LineUp(1)
			case "down":
				c.viewport.LineDown(1)
			case "page_up":
				c.viewport.HalfViewUp()
			case "page_down":
				c.viewport.HalfViewDown()
			}
		}
	}

	c.viewport, cmd = c.viewport.Update(msg)
	return *c, cmd
}

// View renders the chat panel
func (c *ChatPanel) View() string {
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Width(c.width - 2).
		Align(lipgloss.Center).
		Render("#" + c.channel)

	content := c.viewport.View()

	return c.style.Width(c.width).Height(c.height).Render(
		lipgloss.JoinVertical(lipgloss.Left, header, content),
	)
}

// renderContent updates the viewport content with formatted messages
func (c *ChatPanel) renderContent() {
	var sb strings.Builder

	for _, msg := range c.messages {
		username := lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Render(msg.Username)
		if msg.IsAction {
			// Format action messages differently
			actionText := lipgloss.NewStyle().Foreground(lipgloss.Color("63")).Render(
				username + " " + msg.Content,
			)
			sb.WriteString(actionText + "\n")
		} else {
			// Format normal messages
			sb.WriteString(username + ": " + msg.Content + "\n")
		}
	}

	c.viewport.SetContent(sb.String())
}
