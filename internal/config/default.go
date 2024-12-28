package config

// GetDefaultConfig returns the default configuration
func GetDefaultConfig() *Config {
	return &Config{
		ServerDirectory:  "./bedrock_server",
		BackupDirectory: "./bedrock_server_backups",
		BackupInterval:  1440, // 24 hours in minutes
		BackupsToKeep:   7,
	}
}

// DefaultConfigYAML returns the default configuration as a YAML string
func DefaultConfigYAML() string {
	return `# Bedrock Server Manager
# Configuration file

# Server directory. This is the directory where the server will be installed.
server_directory: ./bedrock_server

# BACKUP SETTINGS
# Directory where backups will be stored
backup_directory: ./bedrock_server_backups

# Backup interval in minutes (default: 1440 = 24 hours)
backup_interval: 1440

# Number of backups to keep (set to 0 to keep all backups)
backups_to_keep: 7
`
}