# Twitch Chat TUI

A terminal-based Twitch chat client built with Go and the Charm stack (Bubble Tea, Lip Gloss, Bubbles).

## Features

- Multi-panel Twitch chat display
- Dynamic layout that adapts to terminal size
- ASCII art emotes
- Configurable channels and appearance
- Tab-based navigation

## Installation

```bash
# Clone the repository
git clone https://github.com/daryn/twitch_chat_tui.git
cd twitch_chat_tui

# Build the application
go build -o twitch_chat ./cmd/twitch_chat

# Run the application
./twitch_chat
```

## Configuration

Create a configuration file in `~/.config/twitch_chat_tui/config.yaml` with your Twitch credentials and channel preferences.

Example configuration:
```yaml
twitch:
  username: your_username
  oauth_token: oauth:your_token
channels:
  - channel1
  - channel2
  - channel3
```

## Usage

- Use arrow keys to navigate between chats
- Tab to switch between channel tabs
- Ctrl+N to open a new channel
- Ctrl+C to quit

## Development

This project follows the standard Go project layout. 