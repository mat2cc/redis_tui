package tui

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

func (dm *Details) Reset() {
	dm.key = ""
	dm.data = ""
	dm.redis_type = ""
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

	if dm.key == "" && dm.redis_type == "" {
		return style.
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Bold(true).
			Render("Select a key to view details")
	}
	header := lipgloss.NewStyle().
		Width(dm.width).
		Border(lipgloss.DoubleBorder(), false, false, true).
		Render(fmt.Sprintf("Key: %s\nType: %s", dm.key, dm.redis_type))
	out := lipgloss.JoinVertical(lipgloss.Top, header, dm.data)
	return style.Render(out)
}
