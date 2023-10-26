package main

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TablePrintList struct {
	width int
	table table.Model
	List  []*PrintItem
}

func (pl *TablePrintList) Init() tea.Cmd {
	return nil
}

func (pl *TablePrintList) ToggleExpand() {
	n := pl.GetCurrent()
	n.expanded = !n.expanded
}

func (pl *TablePrintList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
    pl.table.SetColumns(createTableCols(pl.width))
		pl.table.SetWidth(pl.width)
		pl.table.SetHeight(msg.Height - 8)

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

func (pl *TablePrintList) GetCurrent() *Node {
  return pl.List[pl.table.Cursor()].Node
}

func (pl *TablePrintList) View() string {
	style := lipgloss.
		NewStyle().
    Width(pl.width). // subtract 2 for the border
		Border(lipgloss.RoundedBorder())
	return style.Render(pl.table.View())
}

func createTableCols(width int) []table.Column{
	return []table.Column{
		{Title: "", Width: 1},
		{Title: "", Width: width - 4 },
	}
}

func NewTable() *TablePrintList {
	rows := []table.Row{}

	t := table.New(
		table.WithColumns(createTableCols(1)),
		table.WithRows(rows),
		table.WithFocused(true),
	)

	return &TablePrintList{
		table: t,
	}
}
