package main

import (
	"bsm/internal/config"
	"bsm/internal/server"
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

	// Get default config
	defaultConfig := config.GetDefaultConfig()

	// Save config file
	if err := defaultConfig.SaveConfig(configPath); err != nil {
		fmt.Printf("Error creating config file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Config file created at %s\n", configPath)
}

func handleServer() {
	serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
	serverCmd.Parse(os.Args[2:])

	if serverCmd.NArg() < 1 {
		fmt.Println("Usage: bsm server [setup|update]")
		os.Exit(1)
	}

	subcommand := serverCmd.Arg(0)

	switch subcommand {
	case "setup":
		if serverCmd.NArg() < 2 {
			fmt.Println("Usage: bsm server setup [version]")
			fmt.Println("Example: bsm server setup 1.21.51.02")
			os.Exit(1)
		}

		version := serverCmd.Arg(1)
		downloadURL := fmt.Sprintf("https://www.minecraft.net/bedrockdedicatedserver/bin-linux/bedrock-server-%s.zip", version)
		
		// Load config
		cfg, err := config.LoadConfig("config.yaml")
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Setting up server version %s...\n", version)
		if err := server.SetupServer(downloadURL, cfg); err != nil {
			fmt.Printf("Error setting up server: %v\n", err)
			os.Exit(1)
		}
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

func printUsage() {
	fmt.Println(`Usage: bsm [command]

Commands:
  config                   Generate config file
  server setup {version}   Setup new server
  server update {version}  Update server using download URL`)
}