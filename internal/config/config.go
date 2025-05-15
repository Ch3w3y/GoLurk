package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	Twitch   TwitchConfig `yaml:"twitch"`
	Channels []string     `yaml:"channels"`
	UI       UIConfig     `yaml:"ui"`
}

// TwitchConfig holds Twitch-specific settings
type TwitchConfig struct {
	Username   string `yaml:"username"`
	OAuthToken string `yaml:"oauth_token"`
}

// UIConfig holds UI-specific settings
type UIConfig struct {
	Theme             string `yaml:"theme"`
	UseAsciiEmotes    bool   `yaml:"use_ascii_emotes"`
	MaxMessagesBuffer int    `yaml:"max_messages_buffer"`
}

// DefaultConfig creates a default configuration
func DefaultConfig() *Config {
	return &Config{
		Twitch: TwitchConfig{
			Username:   "",
			OAuthToken: "",
		},
		Channels: []string{},
		UI: UIConfig{
			Theme:             "default",
			UseAsciiEmotes:    true,
			MaxMessagesBuffer: 500,
		},
	}
}

// LoadConfig loads the configuration file from the default location
func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return DefaultConfig(), err
	}

	configDir := filepath.Join(homeDir, ".config", "twitch_chat_tui")
	configPath := filepath.Join(configDir, "config.yaml")

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	return LoadConfigFromFile(configPath)
}

// LoadConfigFromFile loads the configuration from the specified file
func LoadConfigFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return DefaultConfig(), err
	}

	config := DefaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return config, err
	}

	return config, nil
}

// SaveConfig saves the configuration to the default location
func SaveConfig(config *Config) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configDir := filepath.Join(homeDir, ".config", "twitch_chat_tui")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.yaml")
	return SaveConfigToFile(config, configPath)
}

// SaveConfigToFile saves the configuration to the specified file
func SaveConfigToFile(config *Config, path string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
