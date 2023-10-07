package main

import tea "github.com/charmbracelet/bubbletea"

type PrintList struct {
	List   []*PrintItem
	cursor int
}

func (pl *PrintList) View() string {
	str := ""
	for i, item := range pl.List {
		if i == pl.cursor {
			str += ">"
		} else {
			str += " "
		}
		str += item.Print()
	}

	return str
}

type updatePL struct {
	root_node *Node
}

func (pl *PrintList) Init() tea.Cmd {
	return nil
}

func (pl *PrintList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
	str += pi.Node.Value + "\n"

	return str
}

func GeneratePrintList(nodes *Node, depth int) []*PrintItem {
	print_list := []*PrintItem{}

	for _, node := range nodes.Children {
		print_list = append(print_list, &PrintItem{node, depth})
		if node.expanded {
			print_list = append(print_list, GeneratePrintList(node, depth+1)...)
		}
	}

	return print_list
}
