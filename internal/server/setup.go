package server

import (
	"fmt"
	"os"
	"path/filepath"

	"bsm/internal/config"
	"bsm/utils"
)

// SetupServer downloads and sets up the Bedrock server
func SetupServer(downloadURL string, cfg *config.Config) error {
	// Create temporary directory for download
	tmpDir, err := os.MkdirTemp("", "bedrock-server")
	if err != nil {
		return fmt.Errorf("error creating temp directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Download server zip
	zipPath := filepath.Join(tmpDir, "server.zip")
	fmt.Println("Downloading server...")
	if err := utils.DownloadFile(downloadURL, zipPath); err != nil {
		return fmt.Errorf("error downloading server: %v", err)
	}

	// Create server directory if it doesn't exist
	if err := os.MkdirAll(cfg.ServerDirectory, 0755); err != nil {
		return fmt.Errorf("error creating server directory: %v", err)
	}

	// Extract zip file
	fmt.Println("Extracting server files...")
	if err := utils.ExtractZip(zipPath, cfg.ServerDirectory); err != nil {
		return fmt.Errorf("error extracting server: %v", err)
	}

	fmt.Printf("Server setup complete! Server installed in: %s\n", cfg.ServerDirectory)
	return nil
}