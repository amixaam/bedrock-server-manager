package config

// GetDefaultConfig returns the default configuration
func GetDefaultConfig() *Config {
	return &Config{
		ServerDirectory:  "./server",
		BackupDirectory: "./backups",
		ServerName: "Bedrock Server",
		BackupInterval:  1440, // 24 hours in minutes
		BackupsToKeep:   7,
		WorldsDirectory: "./worlds",
		WorldDefaults: WorldDefaults{
			LevelName:    "default_world",
			Seed:         "",
			Gamemode:     "survival",
			Difficulty:   "normal",
			AllowList:    true,
			ServerPort:   19132,
			ViewDistance: 32,
			TickDistance: 4,
			MaxPlayers:   10,
		},
	}
}

// DefaultConfigYAML returns the default configuration as a YAML string
func DefaultConfigYAML() string {
	return `
# Bedrock Server Manager
# Configuration file

# Server directory. This is the directory where the server will be installed.
server_directory: ./bedrock_server

# Directory where world configurations are stored
worlds_directory: ./worlds

# BACKUP SETTINGS
# Directory where backups will be stored
backup_directory: ./bedrock_server_backups

# Backup interval in minutes (default: 1440 = 24 hours)
backup_interval: 1440

# Number of backups to keep (set to 0 to keep all backups)
backups_to_keep: 7

# Server name used for server-name property
server_name: "Bedrock Server"

# Default world settings
world_defaults:
  # World settings
  level_name: world
  seed: ""
  gamemode: survival
  difficulty: normal
  allow_list: true
  server_port: 19132
  view_distance: 32
  tick_distance: 4
  max_players: 10
`
}