package models

import (
	"github.com/daryn/twitch_chat_tui/internal/config"
	"github.com/daryn/twitch_chat_tui/internal/irc"
	"github.com/daryn/twitch_chat_tui/internal/ui"
)

// App represents the application state
type App struct {
	Config          *config.Config
	Client          *irc.Client
	Layout          *ui.Layout
	ChatPanels      map[string]*ui.ChatPanel
	ActivePaneIndex int
	Width           int
	Height          int
}

// NewApp creates a new application state
func NewApp(cfg *config.Config) *App {
	// Create a new IRC client
	client := irc.NewClient(cfg.Twitch.Username, cfg.Twitch.OAuthToken)

	// Create a new layout
	layout := ui.NewLayout()

	// Create chat panels for each channel
	chatPanels := make(map[string]*ui.ChatPanel)
	for _, channel := range cfg.Channels {
		chatPanels[channel] = ui.NewChatPanel(channel)
	}

	return &App{
		Config:          cfg,
		Client:          client,
		Layout:          layout,
		ChatPanels:      chatPanels,
		ActivePaneIndex: 0,
	}
}

// AddChannel adds a new channel to the application
func (a *App) AddChannel(channel string) {
	// Normalize channel name (remove # if present)
	normalizedChannel := channel
	if len(normalizedChannel) > 0 && normalizedChannel[0] == '#' {
		normalizedChannel = normalizedChannel[1:]
	}

	// Add to config if not already present
	found := false
	for _, ch := range a.Config.Channels {
		if ch == normalizedChannel {
			found = true
			break
		}
	}

	if !found {
		a.Config.Channels = append(a.Config.Channels, normalizedChannel)
	}

	// Join channel in IRC
	a.Client.JoinChannel(normalizedChannel)

	// Create chat panel if it doesn't exist
	if _, exists := a.ChatPanels[normalizedChannel]; !exists {
		a.ChatPanels[normalizedChannel] = ui.NewChatPanel(normalizedChannel)
	}
}

// RemoveChannel removes a channel from the application
func (a *App) RemoveChannel(channel string) {
	// Normalize channel name (remove # if present)
	normalizedChannel := channel
	if len(normalizedChannel) > 0 && normalizedChannel[0] == '#' {
		normalizedChannel = normalizedChannel[1:]
	}

	// Remove from config
	for i, ch := range a.Config.Channels {
		if ch == normalizedChannel {
			a.Config.Channels = append(a.Config.Channels[:i], a.Config.Channels[i+1:]...)
			break
		}
	}

	// Leave channel in IRC
	a.Client.LeaveChannel(normalizedChannel)

	// Remove chat panel
	delete(a.ChatPanels, normalizedChannel)
}

// UpdateSize updates the size of the application
func (a *App) UpdateSize(width, height int) {
	a.Width = width
	a.Height = height
	a.Layout.UpdateSize(width, height)

	// Update each panel's size based on layout
	paneWidth := width / a.Layout.GetVisiblePanes()
	paneHeight := height

	for _, panel := range a.ChatPanels {
		panel.SetSize(paneWidth, paneHeight)
	}
}

// GetVisibleChannels returns the channels currently visible
func (a *App) GetVisibleChannels() []string {
	visibleCount := min(len(a.Config.Channels), a.Layout.GetVisiblePanes())
	return a.Config.Channels[:visibleCount]
}

// GetChannelPanels returns the chat panels for the visible channels
func (a *App) GetChannelPanels() []*ui.ChatPanel {
	visible := a.GetVisibleChannels()
	panels := make([]*ui.ChatPanel, len(visible))

	for i, channel := range visible {
		panels[i] = a.ChatPanels[channel]
	}

	return panels
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ProcessMessage processes an incoming IRC message
func (a *App) ProcessMessage(msg irc.Message) {
	// Get channel name (removing # prefix if present)
	channel := msg.Channel
	if len(channel) > 0 && channel[0] == '#' {
		channel = channel[1:]
	}

	// Check if we have a panel for this channel
	panel, exists := a.ChatPanels[channel]
	if !exists {
		// Create a new panel if it doesn't exist
		panel = ui.NewChatPanel(channel)
		a.ChatPanels[channel] = panel

		// Add to channels list if not already present
		found := false
		for _, ch := range a.Config.Channels {
			if ch == channel {
				found = true
				break
			}
		}

		if !found {
			a.Config.Channels = append(a.Config.Channels, channel)
		}
	}

	// Add message to panel
	panel.AddMessage(msg.ToChatMessage())
}
