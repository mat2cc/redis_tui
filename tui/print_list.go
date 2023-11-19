package tui

type updatePL struct {
	root_node *Node
}

type PrintItem struct {
	Node  *Node
	Depth int
}

func (pi *PrintItem) Print() string {
	str := ""
	for i := 0; i < pi.Depth; i++ {
		str += "  "
	}
	str += pi.Node.Value

	return str
}

// convert the tree to a list of PrintItems
func GeneratePrintList(root_node *Node, depth int) []*PrintItem {
	print_list := []*PrintItem{}
	for _, node := range root_node.Children {
		print_list = append(print_list, &PrintItem{node, depth})
		if node.Expanded {
			print_list = append(print_list, GeneratePrintList(node, depth+1)...)
		}
	}

	return print_list
}
