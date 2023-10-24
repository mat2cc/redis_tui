package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type TablePrintList struct {
	width int
	table table.Model
	List  []*PrintItem
}

func (pl *TablePrintList) Init() tea.Cmd {
	return nil
}

func (pl *TablePrintList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
  case updatePL:
		msg.root_node.expanded = !msg.root_node.expanded
		pl.List = GeneratePrintList(msg.root_node, 0)
    pl.table.SetRows(pl.GetRows())
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
		case "x":
			return pl, tea.Batch(
				tea.Printf("Let's go to %s!", pl.table.SelectedRow()[1]),
			)
		}
	}
	pl.table, cmd = pl.table.Update(msg)

	return pl, cmd
}

func (pl *TablePrintList) GetRows() []table.Row {
  rows := []table.Row{}
  for _, item := range pl.List {
		prexfix := ""
		postfix := ""
		if len(item.Node.Children) > 0 {
			if item.Node.expanded {
				prexfix += "v"
			} else {
				prexfix += ">"
			}
		} else {
			postfix = fmt.Sprintf(" [%s]", item.Node.RedisType)
			prexfix += ""
		}
    rows = append(rows, table.Row{prexfix, item.Print() + postfix})
  }
  return rows
}

func (pl *TablePrintList) View() string {
	return pl.table.View()
}

func NewTable() *TablePrintList {
	columns := []table.Column{
		{Title: "", Width: 1},
		{Title: "", Width: 80},
	}
	rows := []table.Row{}

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
