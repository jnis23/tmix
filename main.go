package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"tmix/internal/config"
	"tmix/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
)

type Flags struct {
	ConfigPath *string
}

func parseFlags() *Flags {
	flags := &Flags{
		ConfigPath: flag.String(
			"configPath",
			"",
			"Path to your config file. This will default to $HOME/.config/tmix/config.toml",
		),
	}
	flag.Parse()
	return flags
}

func main() {
	flags := parseFlags()
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("Failed:", err)
		os.Exit(1)
	}
	defer f.Close()

	config := config.LoadConfig(flags.ConfigPath)
	model, err := tui.New(config)
	if err != nil {
		log.Fatalf("Failed to start UI: %v", err)
	}

	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
	)

	if _, err := p.Run(); err != nil {
		log.Fatalf("Start tui failed: %v", err)
	}
}
