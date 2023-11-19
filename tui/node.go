package tui

import (
	"log"
	"strings"

	"github.com/redis/go-redis/v9"
)

type Node struct {
	Children  []*Node
	Value     string
	FullKey   string
	RedisType string
	Expanded  bool
}

// recursively add a child node to the tree
func (n *Node) AddChild(key []string, full string, redis *redis.Client, search_string string) {
	if len(key) == 0 {
		return
	}

	for _, child := range n.Children {
		if child.Value == key[0] {
			child.AddChild(key[1:], full, redis, search_string)
			return
		}
	}
	rt, err := redis.Type(ctx, full).Result()
	if err != nil {
		log.Fatal(err)
	}
	new_node := &Node{Value: key[0], FullKey: full, RedisType: rt}
  // expand if the full key contains the search string
	if search_string != "" && strings.Contains(full, search_string) {
		new_node.Expanded = true
	}
	new_node.AddChild(key[1:], full, redis, search_string)
	n.Children = append(n.Children, new_node)
}

func (n *Node) Print(padding int) string {
	str := n.Value + "\n"
	for _, child := range n.Children {
		if child.Expanded {
      // left pad expanded nodes
			for i := 0; i < padding; i++ {
				str += " "
			}
			str += child.Print(padding + 2)
		}
	}

	return str
}
