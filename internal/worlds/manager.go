package worlds

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"bsm/internal/config"
	"bsm/utils"
)

type World struct {
	Name           string
	PropertiesPath string
}

type WorldManager struct {
	ServerDir     string
	WorldsDir     string
	ActiveWorld   string
	Defaults      config.WorldDefaults
	ServerName    string
}

func NewWorldManager(serverDir, worldsDir string, defaults config.WorldDefaults, serverName string) *WorldManager {
	return &WorldManager{
		ServerDir:  serverDir,
		WorldsDir:  worldsDir,
		Defaults:   defaults,
		ServerName: serverName,
	}
}

// promptString asks for user input with a default value
func promptString(prompt string, defaultValue string) string {
	fmt.Printf("%s [%s]: ", prompt, defaultValue)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}
	return input
}

// promptInt asks for numeric user input with a default value
func promptInt(prompt string, defaultValue int) int {
	input := promptString(prompt, strconv.Itoa(defaultValue))
	value, err := strconv.Atoi(input)
	if err != nil {
		return defaultValue
	}
	return value
}

// promptBool asks for boolean user input with a default value
func promptBool(prompt string, defaultValue bool) bool {
	defaultStr := "no"
	if defaultValue {
		defaultStr = "yes"
	}
	input := strings.ToLower(promptString(prompt, defaultStr))
	return input == "yes" || input == "y" || (input == "" && defaultValue)
}

// ListWorlds returns a list of available worlds
func (wm *WorldManager) ListWorlds() ([]World, error) {
	// Create worlds directory if it doesn't exist
	if err := os.MkdirAll(wm.WorldsDir, 0755); err != nil {
		return nil, fmt.Errorf("error creating worlds directory: %v", err)
	}

	entries, err := os.ReadDir(wm.WorldsDir)
	if err != nil {
		return nil, fmt.Errorf("error reading worlds directory: %v", err)
	}

	var worlds []World
	for _, entry := range entries {
		if entry.IsDir() {
			propPath := filepath.Join(wm.WorldsDir, entry.Name(), "server.properties")
			if _, err := os.Stat(propPath); err == nil {
				worlds = append(worlds, World{
					Name:           entry.Name(),
					PropertiesPath: propPath,
				})
			}
		}
	}

	return worlds, nil
}

// GetActiveWorld reads the current level-name from server.properties
func (wm *WorldManager) GetActiveWorld() (string, error) {
	serverProps := filepath.Join(wm.ServerDir, "server.properties")
	data, err := os.ReadFile(serverProps)
	if err != nil {
		return "", fmt.Errorf("error reading server.properties: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "level-name=") {
			return strings.TrimPrefix(line, "level-name="), nil
		}
	}

	return "", fmt.Errorf("level-name not found in server.properties")
}

// SwitchWorld changes the active world
func (wm *WorldManager) SwitchWorld(worldName string) error {
	worldDir := filepath.Join(wm.WorldsDir, worldName)
	
	// Check if world exists
	if _, err := os.Stat(worldDir); err != nil {
		return fmt.Errorf("world %s not found", worldName)
	}

	// Copy server.properties
	serverProps := filepath.Join(wm.ServerDir, "server.properties")
	worldProps := filepath.Join(worldDir, "server.properties")
	if err := utils.CopyFile(worldProps, serverProps); err != nil {
		return fmt.Errorf("error copying properties: %v", err)
	}

	// Copy allowlist.json
	serverAllowlist := filepath.Join(wm.ServerDir, "allowlist.json")
	worldAllowlist := filepath.Join(worldDir, "allowlist.json")
	if err := utils.CopyFile(worldAllowlist, serverAllowlist); err != nil {
		return fmt.Errorf("error copying allowlist: %v", err)
	}

	return nil
}

// CreateWorld creates a new world with custom properties
func (wm *WorldManager) CreateWorld() error {
	// Get world settings from user
	levelName := promptString("Enter world name", wm.Defaults.LevelName)
	seed := promptString("Enter seed (leave empty for random)", wm.Defaults.Seed)
	gamemode := promptString("Enter gamemode (survival/creative/adventure)", wm.Defaults.Gamemode)
	difficulty := promptString("Enter difficulty (peaceful/easy/normal/hard)", wm.Defaults.Difficulty)
	allowList := promptBool("Enable allow list? (yes/no)", wm.Defaults.AllowList)
	serverPort := promptInt("Enter server port", wm.Defaults.ServerPort)
	viewDistance := promptInt("Enter view distance", wm.Defaults.ViewDistance)
	tickDistance := promptInt("Enter tick distance", wm.Defaults.TickDistance)
	maxPlayers := promptInt("Enter max players", wm.Defaults.MaxPlayers)

	// Create world directory
	worldDir := filepath.Join(wm.WorldsDir, levelName)
	if err := os.MkdirAll(worldDir, 0755); err != nil {
		return fmt.Errorf("error creating world directory: %v", err)
	}

	// Create server.properties
	props := map[string]string{
		"server-name":    fmt.Sprintf("%s - %s", wm.ServerName, levelName),
		"level-name":     levelName,
		"server-port":    strconv.Itoa(serverPort),
		"gamemode":       gamemode,
		"difficulty":     difficulty,
		"allow-cheats":   "false",
		"view-distance":  strconv.Itoa(viewDistance),
		"tick-distance":  strconv.Itoa(tickDistance),
		"max-players":    strconv.Itoa(maxPlayers),
		"allow-list":     strconv.FormatBool(allowList),
	}
	if seed != "" {
		props["level-seed"] = seed
	}

	// Create properties file
	if err := wm.createPropertiesFile(filepath.Join(worldDir, "server.properties"), props); err != nil {
		return fmt.Errorf("error creating properties file: %v", err)
	}

	// Create empty allowlist.json
	allowlistPath := filepath.Join(worldDir, "allowlist.json")
	if err := os.WriteFile(allowlistPath, []byte("[]"), 0644); err != nil {
		return fmt.Errorf("error creating allowlist.json: %v", err)
	}

	fmt.Printf("Created world '%s' in %s\n", levelName, worldDir)
	return nil
}

func (wm *WorldManager) createPropertiesFile(path string, props map[string]string) error {
	// First read the template properties file from the server directory
	templatePath := filepath.Join(wm.ServerDir, "server.properties")
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("error reading template properties: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	newLines := make([]string, 0, len(lines))

	for _, line := range lines {
		if line == "" || strings.HasPrefix(line, "#") {
			newLines = append(newLines, line)
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			newLines = append(newLines, line)
			continue
		}

		key := strings.TrimSpace(parts[0])
		if newValue, exists := props[key]; exists {
			newLines = append(newLines, key + "=" + newValue)
		} else {
			newLines = append(newLines, line)
		}
	}

	return os.WriteFile(path, []byte(strings.Join(newLines, "\n")), 0644)
}