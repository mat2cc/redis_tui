package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Details struct {
    key string
    open bool
}

func (dm *Details) Init() tea.Cmd {
	return nil 
}


func (dm *Details) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
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
    style := lipgloss.NewStyle().Border(lipgloss.RoundedBorder())
    return style.Render("Details")
}
