package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Details struct {
	key   string
	data  string
	open  bool
	width int
}

func (dm *Details) Init() tea.Cmd {
	return nil
}

type setDetailsMessage struct {
	key   string
	data  string
}

func (dm *Details) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
  case setDetailsMessage:
    dm.key = msg.key
    dm.data = msg.data
	case tea.WindowSizeMsg:
		dm.width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "d":
			dm.open = !dm.open
			return dm, nil
		}
	}
	return dm, nil
}

func (dm *Details) View() string {
	style := lipgloss.
		NewStyle().
		Width(dm.width/2 - 10).
		Border(lipgloss.RoundedBorder())

  out := fmt.Sprintf("Key: %s\nData: %s", dm.key, dm.data)
	return style.Render(out)
}
