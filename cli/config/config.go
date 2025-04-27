package config

import (
	"encoding/json"
	"os"
)

// Config defines the structure for GoBot CLI configuration.
// It specifies the locations of important files and directories used by the CLI.
type Config struct {
	// CommandsDir is the directory containing Discord bot command files
	CommandsDir string `json:"commandsDir"`

	// MainFile is the path to the main entry point of the bot
	MainFile string `json:"mainFile"`

	// BotFile is the path to the core bot implementation file
	BotFile string `json:"botFile"`
}

// DefaultConfig provides the default configuration values used when
// no config file is present or when config values are missing.
var DefaultConfig = Config{
	CommandsDir: "commands",
	MainFile:    "main.go",
	BotFile:     "bot/bot.go",
}

// LoadConfig attempts to load configuration from Gobot.config.json.
// If no config file is found, it returns the default configuration.
func LoadConfig() (*Config, error) {
	config := DefaultConfig

	// Look for config file variations
	configFiles := []string{
		"Gobot.config.json",
		"gobot.config.json",
		".gobot.config.json",
	}

	var configFile string
	for _, file := range configFiles {
		if _, err := os.Stat(file); err == nil {
			configFile = file
			break
		}
	}

	if configFile == "" {
		return &config, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return &config, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return &config, err
	}

	return &config, nil
}
