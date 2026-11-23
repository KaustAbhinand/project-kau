package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "kau/services"
)

func main() {
	fmt.Println("Kau - the installation automation tool")
	fmt.Println("Version: 0.0.1")
	fmt.Println("Your commands: kau <commmand>")
	fmt.Println("help - show this help message")
	fmt.Println("install - install and configure the entire project")
	fmt.Println("version - show the version number")
	fmt.Println("exit - exit the program")
	fmt.Println("Please enter a command:")

	// Read command.
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	command := strings.TrimSpace(scanner.Text())

	switch command {
	case "help":
		fmt.Println("Available commands: help, install, version, exit")
	case "install":
		fmt.Println("Installing...")
		absPath, err := filepath.Abs(".")
		if err != nil {
			fmt.Println("Error resolving path:", err)
			os.Exit(1)
		}

		// Traverse the repo ----> traverse.go
		tasks, err := traverseRepo(absPath)
		if err != nil {
			fmt.Println("Error scanning repository:", err)
			os.Exit(1)
		}

		// Run installation. ----> installer.go
		if err := runInstallations(tasks); err != nil {
			fmt.Println("Error installing dependencies:", err)
			os.Exit(1)
		}

		// Run configuration. --------> configure.go
		if err := runConfigurations(tasks); err != nil {
			fmt.Println("Error configuring project:", err)
			os.Exit(1)
		}

		fmt.Println("Installed and configured!")

	case "version":
		fmt.Println("Version: 0.0.1")

	case "exit":
		os.Exit(0)

	default:
		fmt.Println("Invalid command.")
	}
}
