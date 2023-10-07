package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Search struct {
	input  textinput.Model
	active bool
}

func NewSearch() Search {
	return Search{
		input: textinput.New(),
	}
}

func (s Search) Init() tea.Cmd {
	return nil
}

func (s Search) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
        s.input.Update(msg)
	}
	return s, nil
}

func (s Search) View() string {
	style := lipgloss.NewStyle().Width(100).Border(lipgloss.RoundedBorder())
	text := "Search\n"
	text += style.Render(s.input.View())
	return text
}
