package main

import (
	"bsm/internal/backup"
	"bsm/internal/config"
	"bsm/internal/server"
	"bsm/internal/worlds"
	"flag"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "config":
		handleConfig()
	case "server":
		handleServer()
	case "world":
		handleWorlds()
	case "backup":
		handleBackup()
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func handleConfig() {
	configPath := "config.yaml"
	
	// Check if config file already exists
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("Config file already exists at %s\n", configPath)
		return
	}

	// Write the default config template
	if err := os.WriteFile(configPath, []byte(config.DefaultConfigYAML()), 0644); err != nil {
		fmt.Printf("Error creating config file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Config file created at %s\n", configPath)
}

func handleServer() {
	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	serverCmd.Parse(os.Args[2:])

	if serverCmd.NArg() < 1 {
		fmt.Println("Usage: bsm server [setup|start|stop|status|update]")
		os.Exit(1)
	}

	// Load config
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	subcommand := serverCmd.Arg(0)
	sm := server.NewServerManager(cfg.ServerDirectory)

	switch subcommand {
	case "setup":
		if serverCmd.NArg() < 2 {
			fmt.Println("Usage: bsm server setup [version]")
			fmt.Println("Example: bsm server setup 1.21.51.02")
			os.Exit(1)
		}

		version := serverCmd.Arg(1)
		downloadURL := fmt.Sprintf("https://www.minecraft.net/bedrockdedicatedserver/bin-linux/bedrock-server-%s.zip", version)
		
		fmt.Printf("Setting up server version %s...\n", version)
		if err := server.SetupServer(downloadURL, cfg); err != nil {
			fmt.Printf("Error setting up server: %v\n", err)
			os.Exit(1)
		}

	case "start":
		fmt.Println("Starting Bedrock server...")
		if err := sm.Start(); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Server started successfully")

	case "stop":
		fmt.Println("Stopping Bedrock server...")
		if err := sm.Stop(); err != nil {
			fmt.Printf("Error stopping server: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Server stopped successfully")

	case "status":
		status, err := sm.Status()
		if err != nil {
			fmt.Printf("Error getting server status: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Server status: %s\n", status)

	case "update":
		if serverCmd.NArg() < 2 {
			fmt.Println("Usage: bsm server update [download_url]")
			os.Exit(1)
		}
		downloadURL := serverCmd.Arg(1)
		fmt.Printf("Updating server from URL: %s\n", downloadURL)
		// TODO: Implement server update

	default:
		fmt.Printf("Unknown server subcommand: %s\n", subcommand)
		os.Exit(1)
	}
}

func handleWorlds() {
	worldsCmd := flag.NewFlagSet("world", flag.ExitOnError)
	worldsCmd.Parse(os.Args[2:])

	if worldsCmd.NArg() < 1 {
		fmt.Println("Usage: bsm world [list|switch|create]")
		os.Exit(1)
	}

	// Load config
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	wm := worlds.NewWorldManager(cfg.ServerDirectory, cfg.WorldsDirectory, cfg.WorldDefaults, cfg.ServerName)
	subcommand := worldsCmd.Arg(0)

	switch subcommand {
	case "list":
		worlds, err := wm.ListWorlds()
		if err != nil {
			fmt.Printf("Error listing worlds: %v\n", err)
			os.Exit(1)
		}

		activeWorld, err := wm.GetActiveWorld()
		if err != nil {
			fmt.Printf("Error getting active world: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Available worlds:")
		for _, world := range worlds {
			if world.Name == activeWorld {
				fmt.Printf("* %s (active)\n", world.Name)
			} else {
				fmt.Printf("  %s\n", world.Name)
			}
		}

	case "switch":
		if worldsCmd.NArg() < 2 {
			fmt.Println("Usage: bsm worlds switch [world_name]")
			os.Exit(1)
		}

		worldName := worldsCmd.Arg(1)
		if err := wm.SwitchWorld(worldName); err != nil {
			fmt.Printf("Error switching world: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Switched to world: %s\n", worldName)

	case "create":
		if err := wm.CreateWorld(); err != nil {
			fmt.Printf("Error creating world: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown worlds subcommand: %s\n", subcommand)
		os.Exit(1)
	}
}

func handleBackup() {
	backupCmd := flag.NewFlagSet("backup", flag.ExitOnError)
	backupCmd.Parse(os.Args[2:])

	if backupCmd.NArg() < 1 {
		fmt.Println("Usage: bsm backup [list|create|restore]")
		os.Exit(1)
	}

	// Load config
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	bm := backup.NewBackupManager(cfg)
	subcommand := backupCmd.Arg(0)

	switch subcommand {
	case "list":
		backups, err := bm.ListBackups()
		if err != nil {
			fmt.Printf("Error listing backups: %v\n", err)
			os.Exit(1)
		}

		if len(backups) == 0 {
			fmt.Println("No backups found")
			return
		}

		for _, wb := range backups {
			fmt.Printf("\nWorld: %s\n", wb.WorldName)
			fmt.Printf("Total backups: %d (%.2f MB)\n", wb.BackupCount, float64(wb.TotalSize)/(1024*1024))
			if len(wb.Backups) > 0 {
				fmt.Println("Recent backups:")
				for _, b := range wb.Backups {
					fmt.Printf("  %s (%.2f MB)\n", 
						b.CreatedAt.Format("2006-01-02 15:04:05"),
						float64(b.Size)/(1024*1024))
				}
			}
		}

	case "create":
		if backupCmd.NArg() < 2 {
			fmt.Println("Usage: bsm backup create [world_name]")
			os.Exit(1)
		}
		worldName := backupCmd.Arg(1)
		if err := bm.CreateBackup(worldName); err != nil {
			fmt.Printf("Error creating backup: %v\n", err)
			os.Exit(1)
		}

	case "restore":
		if backupCmd.NArg() < 2 {
			fmt.Println("Usage: bsm backup restore [world_name]")
			os.Exit(1)
		}
		worldName := backupCmd.Arg(1)
		if err := bm.RestoreBackup(worldName); err != nil {
			fmt.Printf("Error restoring backup: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Printf("Unknown backup subcommand: %s\n", subcommand)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`
Usage: bsm [command]

Commands:
  config                   Generate config file
  server setup {version}   Setup new server
  server start            Start the Bedrock server
  server stop             Stop the Bedrock server
  server status           Check server status
  server update {version}  Update server using download URL
  world list               List all worlds
  world switch {name}      Switch to world {name}
  world create {name}      Create a new world
  backup list              List all backups
  backup create {name}     Create backup {name}
  backup restore {name}    Restore backup {name}
  `)
}