package main

type Node struct {
	Children []*Node
	Value    string
  FullKey string
	expanded bool
}

func (n *Node) AddChild(key []string, full string) {
	if len(key) == 0 {
		return
	}

	for _, child := range n.Children {
		if child.Value == key[0] {
			child.AddChild(key[1:], full)
			return
		}
	}
	new_node := &Node{Value: key[0], FullKey: full}
	new_node.AddChild(key[1:], full)
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

