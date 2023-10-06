package main

type Node struct {
	Children []*Node
	Value    string
	expanded bool
}

func (n *Node) AddChild(key []string) {
	if len(key) == 0 {
		return
	}

	for _, child := range n.Children {
		if child.Value == key[0] {
			child.AddChild(key[1:])
			return
		}
	}
	new_node := &Node{Value: key[0]}
	new_node.AddChild(key[1:])
	n.Children = append(n.Children, new_node)
}

func (n *Node) Print(padding int) string {
	str := n.Value + "\n"
	for _, child := range n.Children {
		if child.expanded {
			for i := 0; i < padding; i++ {
				str += " "
			}
			str += child.Print(padding + 2)
		}
	}

	return str
}

type Printable interface {
	Print() string
}

type PrintList struct {
	List   []*PrintItem
	cursor int
}

func (pl *PrintList) Print() string {
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
