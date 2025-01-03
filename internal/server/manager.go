package server

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

type ServerManager struct {
	serverDir string
	pidFile   string
}

func NewServerManager(serverDir string) *ServerManager {
	return &ServerManager{
		serverDir: serverDir,
		pidFile:   filepath.Join(serverDir, "server.pid"),
	}
}

// Start launches the Bedrock server in the background
func (sm *ServerManager) Start() error {
	// Check if server is already running
	if pid, _ := sm.getServerPID(); pid > 0 {
		return fmt.Errorf("server is already running with PID %d", pid)
	}

	// Check if bedrock_server exists
	serverPath := filepath.Join(sm.serverDir, "bedrock_server")
	if _, err := os.Stat(serverPath); err != nil {
		return fmt.Errorf("bedrock_server not found in %s", sm.serverDir)
	}

	// Make sure the server file is executable
	if err := os.Chmod(serverPath, 0755); err != nil {
		return fmt.Errorf("failed to make server executable: %v", err)
	}

	// Start the server process
	cmd := exec.Command(serverPath)
	cmd.Dir = sm.serverDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	// Write PID file
	if err := os.WriteFile(sm.pidFile, []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0644); err != nil {
		return fmt.Errorf("failed to write PID file: %v", err)
	}

	// Start a goroutine to clean up PID file when process exits
	go func() {
		cmd.Wait()
		os.Remove(sm.pidFile)
	}()

	return nil
}

// Stop gracefully stops the Bedrock server
func (sm *ServerManager) Stop() error {
	pid, err := sm.getServerPID()
	if err != nil {
		return fmt.Errorf("failed to get server PID: %v", err)
	}

	if pid <= 0 {
		return fmt.Errorf("server is not running")
	}

	// Try to terminate gracefully first
	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %v", err)
	}

	// Send SIGTERM for graceful shutdown
	if err := process.Signal(syscall.SIGTERM); err != nil {
		return fmt.Errorf("failed to send termination signal: %v", err)
	}

	// Wait for up to 30 seconds for the server to shut down
	for i := 0; i < 30; i++ {
		if !sm.IsRunning() {
			return nil
		}
		time.Sleep(time.Second)
	}

	// Force kill if still running
	if err := process.Kill(); err != nil {
		return fmt.Errorf("failed to force kill server: %v", err)
	}

	// Clean up PID file
	os.Remove(sm.pidFile)
	return nil
}

// Status returns the current status of the server
func (sm *ServerManager) Status() (string, error) {
	pid, err := sm.getServerPID()
	if err != nil {
		return "", fmt.Errorf("failed to get server status: %v", err)
	}

	if pid <= 0 {
		return "stopped", nil
	}

	if sm.IsRunning() {
		return fmt.Sprintf("running (PID: %d)", pid), nil
	}

	// Clean up stale PID file
	os.Remove(sm.pidFile)
	return "stopped", nil
}

// IsRunning checks if the server process is currently running
func (sm *ServerManager) IsRunning() bool {
	pid, err := sm.getServerPID()
	if err != nil || pid <= 0 {
		return false
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Send signal 0 to check if process exists
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// getServerPID reads the PID from the PID file
func (sm *ServerManager) getServerPID() (int, error) {
	if _, err := os.Stat(sm.pidFile); os.IsNotExist(err) {
		return 0, nil
	}

	data, err := os.ReadFile(sm.pidFile)
	if err != nil {
		return 0, err
	}

	var pid int
	if _, err := fmt.Sscanf(string(data), "%d", &pid); err != nil {
		return 0, err
	}

	return pid, nil
}