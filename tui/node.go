package tui

import (
	"log"
	"os"
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

type type_builder struct {
	pipe  redis.Pipeliner
	nodes []*Node
}

func new_type_builder(client *redis.Client) *type_builder {
	pipe := client.Pipeline()
	return &type_builder{
		pipe: pipe,
	}
}

func (tb *type_builder) add_type(node *Node, redis_type string) {
	tb.nodes = append(tb.nodes, node)
	tb.pipe.Type(ctx, node.FullKey)
}

func (n *Node) GenNodes(keys []string, client *redis.Client, search string, opts ModelOptions) {
	tb := new_type_builder(client)
	for _, key := range keys {
		split := strings.Split(key, opts.Delimiter)
		n.AddChild(split, key, client, search, tb)
	}

	if opts.IncludeTypes {
		cmds, err := tb.pipe.Exec(ctx)
		if err != nil {
			log.Fatal(err)
		}

		errLogger := log.New(os.Stderr, "Error getting type for node: ", 0)
		for i, cmd := range cmds {
			rt, err := cmd.(*redis.StatusCmd).Result()
			if err != nil {
				errLogger.Println(tb.nodes[i].FullKey)
			}
			tb.nodes[i].RedisType = rt
		}
	}
}

// recursively add a child node to the tree
func (n *Node) AddChild(key []string, full string, client *redis.Client, search_string string, tb *type_builder) {
	if len(key) == 0 {
		return
	}

	// if you find an existing node with the same name, recurse into that node without adding a new node
	for _, child := range n.Children {
		if child.Value == key[0] {
			child.AddChild(key[1:], full, client, search_string, tb)
			return
		}
	}

	new_node := &Node{Value: key[0], FullKey: full}

	// only get the type if it's a leaf node
	if len(key) == 1 {
		tb.add_type(new_node, full)
	}

	// expand if the full key contains the search string
	if search_string != "" && strings.Contains(full, search_string) {
		new_node.Expanded = true
	}

	new_node.AddChild(key[1:], full, client, search_string, tb)
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
