package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ArefinAlter/zeus-tui/internal/client"
	"github.com/ArefinAlter/zeus-tui/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	version = "0.1.0-alpha.1"
)

func main() {
	daemonURL := flag.String("daemon-url", "http://127.0.0.1:18540", "Zeus daemon IPC HTTP endpoint")
	demoMode := flag.Bool("demo", false, "Run in interactive simulation demo mode")
	showVersion := flag.Bool("version", false, "Show zeus-tui version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("zeus-tui version %s\n", version)
		os.Exit(0)
	}

	c := client.New(*daemonURL, *demoMode)
	m := ui.NewModel(c)

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running zeus-tui: %v\n", err)
		os.Exit(1)
	}
}
