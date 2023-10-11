package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PrintList struct {
	List   []*PrintItem
	width  int
	cursor int
}

func (pl *PrintList) View() string {
	str := ""
	for i, item := range pl.List {
		s := ""
		if len(item.Node.Children) > 0 {
			if item.Node.expanded {
				s += "v "
			} else {
				s += "> "
			}
		} else {
			s += "  "
		}
		str += get_style(i == pl.cursor).Render(s+item.Print()) + "\n"
	}

	style := lipgloss.
		NewStyle().
    Width(pl.width / 2 - 10). // subtract 2 for the border
		Border(lipgloss.RoundedBorder())
	return style.Render(str)
}

func get_style(on_cursor bool) lipgloss.Style {
	if on_cursor {
		return lipgloss.NewStyle().Background(lipgloss.Color("#000000")).Foreground(lipgloss.Color("#ffffff"))
	} else {
		return lipgloss.NewStyle()
	}
}

type updatePL struct {
	root_node *Node
}

func (pl *PrintList) Init() tea.Cmd {
	return nil
}

func (pl *PrintList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
  case tea.WindowSizeMsg:
    pl.width = msg.Width
	case updatePL:
		msg.root_node.expanded = !msg.root_node.expanded
		pl.List = GeneratePrintList(msg.root_node, 0)

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if pl.cursor > 0 {
				pl.cursor--
			}
		case "down", "j":
			if pl.cursor < len(pl.List)-1 {
				pl.cursor++
			}
		}
	}

	return pl, nil
}

func (pl *PrintList) ToggleExpand() {
	n := pl.List[pl.cursor].Node
	n.expanded = !n.expanded
}

type PrintItem struct {
	Node  *Node
	depth int
}

func (pi *PrintItem) Print() string {
	str := ""
	for i := 0; i < pi.depth; i++ {
		str += "  "
	}
	str += pi.Node.Value

	return str
}

func GeneratePrintList(root_node *Node, depth int) []*PrintItem {
	print_list := []*PrintItem{}
	for _, node := range root_node.Children {
		print_list = append(print_list, &PrintItem{node, depth})
		if node.expanded {
			print_list = append(print_list, GeneratePrintList(node, depth+1)...)
		}
	}

	return print_list
}
