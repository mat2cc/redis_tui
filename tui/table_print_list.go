package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TablePrintList struct {
	width  int
	height int
	table  table.Model
	List   []*PrintItem
}

func (pl *TablePrintList) Init() tea.Cmd {
	return nil
}

type resetCursor struct{}

func (pl *TablePrintList) ResetCursor() tea.Msg {
	return resetCursor{}
}

func (pl *TablePrintList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		pl.table.SetColumns(createTableCols(pl.width))
		pl.table.SetWidth(pl.width)
		pl.table.SetHeight(pl.height - MARGIN)

	case updatePL:
		msg.root_node.Expanded = !msg.root_node.Expanded
		pl.List = GeneratePrintList(msg.root_node, 0)
		pl.table.SetRows(pl.GetRows())

	case resetCursor:
		pl.table.SetCursor(0)
	}
	pl.table, cmd = pl.table.Update(msg)

	return pl, cmd
}

// return a list of rows to be displayed in the table
func (pl *TablePrintList) GetRows() []table.Row {
	rows := []table.Row{}
	for _, item := range pl.List {
		prexfix := ""
		postfix := ""
		if len(item.Node.Children) > 0 {
			if item.Node.Expanded {
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

// return a pointer to the current node
func (pl *TablePrintList) GetCurrent() *Node {
	cursor := pl.table.Cursor()
	if cursor < 0 || cursor > len(pl.List) {
		return nil
	}
	item := pl.List[cursor]
	if item == nil {
		return nil
	}
	return item.Node
}

func (pl *TablePrintList) View() string {
	style := lipgloss.
		NewStyle().
		Width(pl.width).
		Height(pl.height).
		Border(lipgloss.RoundedBorder())
	if pl.List == nil || len(pl.List) == 0 {
		return style.
			AlignVertical(lipgloss.Center).
			AlignHorizontal(lipgloss.Center).
			Bold(true).
			Render("Whoops, no keys found for that search value")
	}
	return style.Render(pl.table.View())
}

func createTableCols(width int) []table.Column {
	return []table.Column{
		{Title: "", Width: 1},
		{Title: "", Width: width - MARGIN * 2},
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
