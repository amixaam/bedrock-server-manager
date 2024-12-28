package backup

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"bsm/internal/config"
	"bsm/utils"
)

type Backup struct {
	Name      string
	Path      string
	Size      int64
	CreatedAt time.Time
}

type WorldBackups struct {
	WorldName     string
	Backups      []Backup
	TotalSize    int64
	BackupCount  int
}

type BackupManager struct {
	ServerDir     string
	BackupDir     string
	MaxBackups    int
}

func NewBackupManager(cfg *config.Config) *BackupManager {
	return &BackupManager{
		ServerDir:  cfg.ServerDirectory,
		BackupDir:  cfg.BackupDirectory,
		MaxBackups: cfg.BackupsToKeep,
	}
}

// ListBackups returns a list of all backups grouped by world
func (bm *BackupManager) ListBackups() ([]WorldBackups, error) {
	var worldBackups []WorldBackups
	
	// Read backup directory
	entries, err := os.ReadDir(bm.BackupDir)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("error reading backup directory: %v", err)
	}

	// Process each world's backups
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		worldName := entry.Name()
		worldBackupDir := filepath.Join(bm.BackupDir, worldName)
		backups, totalSize, err := bm.getWorldBackups(worldBackupDir)
		if err != nil {
			return nil, fmt.Errorf("error getting backups for %s: %v", worldName, err)
		}

		// Sort backups by creation time (newest first)
		sort.Slice(backups, func(i, j int) bool {
			return backups[i].CreatedAt.After(backups[j].CreatedAt)
		})

		// Take only the 5 newest backups for display
		displayBackups := backups
		if len(displayBackups) > 5 {
			displayBackups = displayBackups[:5]
		}

		worldBackups = append(worldBackups, WorldBackups{
			WorldName:    worldName,
			Backups:     displayBackups,
			TotalSize:   totalSize,
			BackupCount: len(backups),
		})
	}

	return worldBackups, nil
}

// CreateBackup creates a backup of the specified world
func (bm *BackupManager) CreateBackup(worldName string) error {
	worldPath := filepath.Join(bm.ServerDir, "worlds", worldName)
	if _, err := os.Stat(worldPath); err != nil {
		return fmt.Errorf("world '%s' not found in server directory. Run the server to generate the world first", worldName)
	}

	// Create backup directory for this world
	worldBackupDir := filepath.Join(bm.BackupDir, worldName)
	if err := os.MkdirAll(worldBackupDir, 0755); err != nil {
		return fmt.Errorf("error creating backup directory: %v", err)
	}

	// Create backup filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupPath := filepath.Join(worldBackupDir, fmt.Sprintf("%s_%s.zip", worldName, timestamp))

	// Create zip file
	if err := utils.ZipDirectory(worldPath, backupPath); err != nil {
		return fmt.Errorf("error creating backup: %v", err)
	}

	// Clean up old backups if needed
	if err := bm.cleanOldBackups(worldName); err != nil {
		fmt.Printf("Warning: error cleaning old backups: %v\n", err)
	}

	fmt.Printf("Created backup of '%s' at %s\n", worldName, backupPath)
	return nil
}

// RestoreBackup restores a world from a backup
func (bm *BackupManager) RestoreBackup(worldName string) error {
	// First check if we have write permissions to the worlds directory
	testPath := filepath.Join(bm.ServerDir, "worlds", ".test_write")
	if err := os.WriteFile(testPath, []byte("test"), 0644); err != nil {
		return fmt.Errorf("insufficient permissions to modify worlds directory. Please run with appropriate permissions")
	}
	os.Remove(testPath)

	worldBackupDir := filepath.Join(bm.BackupDir, worldName)
	backups, _, err := bm.getWorldBackups(worldBackupDir)
	if err != nil {
		return fmt.Errorf("error getting backups: %v", err)
	}

	if len(backups) == 0 {
		return fmt.Errorf("no backups found for world '%s'", worldName)
	}

	// Sort backups by creation time (newest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].CreatedAt.After(backups[j].CreatedAt)
	})

	// Display available backups
	fmt.Printf("Available backups for '%s':\n", worldName)
	for i, backup := range backups {
		fmt.Printf("[%d] %s (%.2f MB)\n", i+1, 
			backup.CreatedAt.Format("2006-01-02 15:04:05"),
			float64(backup.Size)/(1024*1024))
	}

	// Get user selection
	var selection int
	fmt.Print("\nEnter backup number to restore (0 to cancel): ")
	fmt.Scanln(&selection)

	if selection == 0 {
		return fmt.Errorf("backup restoration cancelled")
	}
	if selection < 1 || selection > len(backups) {
		return fmt.Errorf("invalid backup selection")
	}

	selectedBackup := backups[selection-1]

	// Confirm restoration
	fmt.Printf("\nWARNING: This will replace the current world '%s' with the backup from %s\n", 
		worldName, selectedBackup.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Print("Are you sure you want to continue? (y/n): ")
	
	var confirm string
	fmt.Scanln(&confirm)
	confirm = strings.ToLower(confirm)
	if confirm != "yes" && confirm != "y" {
		return fmt.Errorf("backup restoration cancelled")
	}

	// Perform restoration
	worldPath := filepath.Join(bm.ServerDir, "worlds", worldName)

	// Remove existing world if it exists
	if _, err := os.Stat(worldPath); err == nil {
		if err := os.RemoveAll(worldPath); err != nil {
			return fmt.Errorf("error removing existing world: %v", err)
		}
	}

	// Extract backup
	if err := utils.UnzipFile(selectedBackup.Path, filepath.Join(bm.ServerDir, "worlds"), worldName); err != nil {
		return fmt.Errorf("error restoring backup: %v", err)
	}

	fmt.Printf("Successfully restored '%s' from backup\n", worldName)
	return nil
}

// Helper function to get backups for a specific world
func (bm *BackupManager) getWorldBackups(worldBackupDir string) ([]Backup, int64, error) {
	var backups []Backup
	var totalSize int64

	entries, err := os.ReadDir(worldBackupDir)
	if err != nil {
		if os.IsNotExist(err) {
			return backups, 0, nil
		}
		return nil, 0, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".zip") {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			backups = append(backups, Backup{
				Name:      entry.Name(),
				Path:      filepath.Join(worldBackupDir, entry.Name()),
				Size:      info.Size(),
				CreatedAt: info.ModTime(),
			})
			totalSize += info.Size()
		}
	}

	return backups, totalSize, nil
}

// Helper function to clean up old backups
func (bm *BackupManager) cleanOldBackups(worldName string) error {
	if bm.MaxBackups <= 0 {
		return nil // Keep all backups
	}

	worldBackupDir := filepath.Join(bm.BackupDir, worldName)
	backups, _, err := bm.getWorldBackups(worldBackupDir)
	if err != nil {
		return err
	}

	if len(backups) <= bm.MaxBackups {
		return nil
	}

	// Sort by creation time (oldest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].CreatedAt.Before(backups[j].CreatedAt)
	})

	// Remove oldest backups
	for i := 0; i < len(backups)-bm.MaxBackups; i++ {
		if err := os.Remove(backups[i].Path); err != nil {
			return fmt.Errorf("error removing old backup %s: %v", backups[i].Name, err)
		}
	}

	return nil
}