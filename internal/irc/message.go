package irc

import (
	"time"
)

// Message represents an IRC message from Twitch
type Message struct {
	Channel      string
	Username     string
	Content      string
	IsAction     bool
	Tags         map[string]string
	TimeReceived time.Time
}

// NewMessage creates a new IRC message
func NewMessage(channel, username, content string, isAction bool) Message {
	return Message{
		Channel:      channel,
		Username:     username,
		Content:      content,
		IsAction:     isAction,
		Tags:         make(map[string]string),
		TimeReceived: time.Now(),
	}
}
