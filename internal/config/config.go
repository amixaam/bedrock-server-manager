package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type WorldDefaults struct {
	LevelName    string `yaml:"level_name"`
	Seed         string `yaml:"seed"`
	Gamemode     string `yaml:"gamemode"`
	Difficulty   string `yaml:"difficulty"`
	AllowList    bool   `yaml:"allow_list"`
	ServerPort   int    `yaml:"server_port"`
	ViewDistance int    `yaml:"view_distance"`
	TickDistance int    `yaml:"tick_distance"`
	MaxPlayers   int    `yaml:"max_players"`
}

type Config struct {
	ServerDirectory  string       `yaml:"server_directory"`
	WorldsDirectory string       `yaml:"worlds_directory"`
	BackupDirectory string       `yaml:"backup_directory"`
	BackupInterval  int          `yaml:"backup_interval"`
	BackupsToKeep   int          `yaml:"backups_to_keep"`
	ServerName      string       `yaml:"server_name"`
	WorldDefaults   WorldDefaults `yaml:"world_defaults"`
}

// LoadConfig loads configuration from the specified file
// Returns default config if file doesn't exist
func LoadConfig(path string) (*Config, error) {
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return GetDefaultConfig(), nil
	}

	// Read the config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return &config, nil
}

	// SaveConfig saves the configuration to the specified file
func (c *Config) SaveConfig(path string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating config directory: %v", err)
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("error marshaling config: %v", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("error writing config file: %v", err)
	}

	return nil
}

// ValidateConfig checks if the configuration is valid
func (c *Config) ValidateConfig() error {
	if c.ServerDirectory == "" {
		return fmt.Errorf("server_directory cannot be empty")
	}
	if c.WorldsDirectory == "" {
		return fmt.Errorf("worlds_directory cannot be empty")
	}
	if c.BackupDirectory == "" {
		return fmt.Errorf("backup_directory cannot be empty")
	}
	if c.BackupInterval < 0 {
		return fmt.Errorf("backup_interval must be non-negative")
	}
	if c.BackupsToKeep < 0 {
		return fmt.Errorf("backups_to_keep must be non-negative")
	}
	return nil
}