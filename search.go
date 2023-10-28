package main

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	noStyle      = lipgloss.NewStyle()
)

type Search struct {
	input      textinput.Model
	active     bool
	width      int
	old_search string
}

func NewSearch() *Search {
	ti := textinput.New()
	ti.Placeholder = "Search"
	ti.SetValue("*")
	ti.Width = 100

	return &Search{
		input: ti,
	}
}

func (s *Search) Init() tea.Cmd {
	return nil
}

func (s *Search) createTextMessage() tea.Cmd {
	if !strings.Contains(s.input.Value(), "*") {
		s.input.SetValue(s.input.Value() + "*")
	}
	s.ToggleActive(false)
	return func() tea.Msg {
		return setTextMessage{s.input.Value()}
	}
}

type setSearchString struct {
	text string
}

type setTextMessage struct {
	text string
}

func (s *Search) ToggleActive(active bool) {
	s.active = active
	if active {
        s.old_search = s.input.Value()
        s.input.SetValue("")

		s.input.Focus()
        s.input.PromptStyle = focusedStyle
        s.input.TextStyle = focusedStyle
	} else {
		s.input.Blur()
        s.input.PromptStyle = noStyle
        s.input.TextStyle = noStyle
	}
	//s.input.Cursor.Style = focusedStyle
	//s.input.TextStyle = focusedStyle
	//s.input.PromptStyle = focusedStyle
}

func (s *Search) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
	case setSearchString:
		s.input.SetValue(msg.text)
	case tea.KeyMsg:
		switch {
    case key.Matches(msg, search_keys.Enter):
			cmd = s.createTextMessage()
			return s, cmd
    case key.Matches(msg, search_keys.Esc):
			s.ToggleActive(false)
			s.input.SetValue(s.old_search)
    case key.Matches(msg, search_keys.Quit):
			return s, tea.Quit
		}
		s.input, cmd = s.input.Update(msg)
	}
	return s, cmd
}

func (s *Search) View() string {
	style := lipgloss.
		NewStyle().
		Width(s.width - MARGIN). // subtract 2 for the border
		Border(lipgloss.RoundedBorder())

	text := style.Render(s.input.View())
	return text
}
