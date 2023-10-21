package main

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type TablePrintList struct {
	width int
	table table.Model
}

func (pl *TablePrintList) Init() tea.Cmd {
	return nil
}

func (pl *TablePrintList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
    switch msg := msg.(type) {
    case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if pl.table.Focused() {
				pl.table.Blur()
			} else {
				pl.table.Focus()
			}
		case "q", "ctrl+c":
			return pl, tea.Quit
		case "enter":
			return pl, tea.Batch(
				tea.Printf("Let's go to %s!", pl.table.SelectedRow()[1]),
			)
		}
	}
	pl.table, cmd = pl.table.Update(msg)

	return pl, cmd
}

func (pl *TablePrintList) View() string {
	return pl.table.View()
}

func NewTable() *TablePrintList {
	columns := []table.Column{
		{Title: "", Width: 2},
		{Title: "", Width: 80},
	}
	rows := []table.Row{
		{"> ", "foo"},
		{"  ", "bar"},
		{"v ", "baz"},
		{"  ", "  bip"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithWidth(80),
		table.WithFocused(true),
	)

	return &TablePrintList{
		table: t,
		width: 80,
	}
}
