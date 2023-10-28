package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Details struct {
	key        string
	data       string
	redis_type string
	width      int
	height     int
}

func (dm *Details) Init() tea.Cmd {
	return nil
}

type setDetailsMessage struct {
	key        string
	redis_type string
	data       string
}

func (dm *Details) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case setDetailsMessage:
		dm.key = msg.key
		dm.redis_type = msg.redis_type
		dm.data = msg.data
	}
	return dm, nil
}

func (dm *Details) View() string {
	style := lipgloss.
		NewStyle().
    Height(dm.height).
		Width(dm.width).
		Border(lipgloss.RoundedBorder())

	out := fmt.Sprintf("Key: %s\nType: %s\n\n%s", dm.key, dm.redis_type, dm.data)
	return style.Render(out)
}
