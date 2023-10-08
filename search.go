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
	ti := textinput.New()
	ti.Placeholder = "Search"
	ti.Focus()
	ti.Width = 100

	return Search{
		input: ti,
	}
}

func (s Search) Init() tea.Cmd {
	return textinput.Blink
}

type setTextMessage struct {
	text string
}

func (s Search) createTextMessage() tea.Cmd {
	return func() tea.Msg {
		return setTextMessage{s.input.Value()}
	}
}

func (s Search) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			cmd = s.createTextMessage()
			s.active = false
			s.input.SetValue("")
			return s, cmd
		case tea.KeyEscape:
			s.active = false
			s.input.SetValue("")
		case tea.KeyCtrlC:
			return s, tea.Quit
		}
		s.input, cmd = s.input.Update(msg)
	}
	return s, cmd
}

func (s Search) View() string {
	style := lipgloss.NewStyle().Width(100).Border(lipgloss.RoundedBorder())
	text := "Search\n"
	text += style.Render(s.input.View())
	return text
}
