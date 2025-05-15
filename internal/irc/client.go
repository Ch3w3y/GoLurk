package irc

import (
	"strings"

	"github.com/Ch3w3y/GoLurk/internal/ui"
	"github.com/gempir/go-twitch-irc/v4"
)

// Client represents a Twitch IRC client
type Client struct {
	twitchClient *twitch.Client
	username     string
	oauthToken   string
	channels     map[string]bool
	messagesChan chan Message
}

// NewClient creates a new Twitch IRC client
func NewClient(username, oauthToken string) *Client {
	return &Client{
		username:     username,
		oauthToken:   oauthToken,
		channels:     make(map[string]bool),
		messagesChan: make(chan Message, 100),
	}
}

// Connect connects to the Twitch IRC server
func (c *Client) Connect() error {
	// Create Twitch client
	c.twitchClient = twitch.NewClient(c.username, c.oauthToken)

	// Set up callbacks
	c.twitchClient.OnPrivateMessage(func(message twitch.PrivateMessage) {
		c.messagesChan <- Message{
			Channel:  message.Channel,
			Username: message.User.DisplayName,
			Content:  message.Message,
			IsAction: strings.HasPrefix(message.Message, "\u0001ACTION") && strings.HasSuffix(message.Message, "\u0001"),
			Tags:     message.Tags,
		}
	})

	c.twitchClient.OnClearChatMessage(func(message twitch.ClearChatMessage) {
		// Handle message clears (timeouts, bans)
		if message.TargetUserID != "" {
			// A user was timed out or banned
			c.messagesChan <- Message{
				Channel:  message.Channel,
				Username: "SYSTEM",
				Content:  "A user message was cleared by a moderator",
				IsAction: true,
			}
		} else {
			// Chat was cleared
			c.messagesChan <- Message{
				Channel:  message.Channel,
				Username: "SYSTEM",
				Content:  "Chat was cleared by a moderator",
				IsAction: true,
			}
		}
	})

	// Connect to Twitch
	return c.twitchClient.Connect()
}

// Disconnect disconnects from the Twitch IRC server
func (c *Client) Disconnect() {
	if c.twitchClient != nil {
		c.twitchClient.Disconnect()
	}
}

// JoinChannel joins a Twitch channel
func (c *Client) JoinChannel(channel string) {
	// Normalize channel name
	channel = NormalizeChannel(channel)

	// Check if already joined
	if _, exists := c.channels[channel]; exists {
		return
	}

	// Join the channel
	c.twitchClient.Join(channel)
	c.channels[channel] = true
}

// LeaveChannel leaves a Twitch channel
func (c *Client) LeaveChannel(channel string) {
	// Normalize channel name
	channel = NormalizeChannel(channel)

	// Check if already joined
	if _, exists := c.channels[channel]; !exists {
		return
	}

	// Leave the channel
	c.twitchClient.Depart(channel)
	delete(c.channels, channel)
}

// SendMessage sends a message to a Twitch channel
func (c *Client) SendMessage(channel, message string) {
	// Normalize channel name
	channel = NormalizeChannel(channel)

	// Send message
	c.twitchClient.Say(channel, message)
}

// GetMessageChan returns the message channel
func (c *Client) GetMessageChan() <-chan Message {
	return c.messagesChan
}

// NormalizeChannel ensures the channel name is properly formatted
func NormalizeChannel(channel string) string {
	channel = strings.ToLower(channel)
	if !strings.HasPrefix(channel, "#") {
		channel = "#" + channel
	}
	return channel
}

// ToChatMessage converts an IRC message to a UI chat message
func (m *Message) ToChatMessage() ui.ChatMessage {
	return ui.ChatMessage{
		Username: m.Username,
		Content:  m.Content,
		IsAction: m.IsAction,
	}
}
