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

func (n *Node) GenNodes(keys []string, client *redis.Client, search string) {
	pipe := client.Pipeline()
	for _, key := range keys {
		split := strings.Split(key, ":")
		n.AddChild(split, key, client, search, pipe)
	}
	cmds, err := pipe.Exec(ctx)
	if err != nil {
		log.Fatal(err)
	}

	errLogger := log.New(os.Stderr, "Error: ", 0)
	goodLogger := log.New(os.Stderr, "Success: ", 0)
	for _, cmd := range cmds {
		rt, err := cmd.(*redis.StatusCmd).Result()
		if err != nil {
			errLogger.Println(err)
		}
		goodLogger.Println(rt)
	}
}
func (n *Node) ddChild(full string, client *redis.Client) {
	client.Type(ctx, full)
}

// recursively add a child node to the tree
func (n *Node) AddChild(key []string, full string, client *redis.Client, search_string string, pipe redis.Pipeliner) {
	if len(key) == 0 {
		return
	}

	for _, child := range n.Children {
		if child.Value == key[0] {
			child.AddChild(key[1:], full, client, search_string, pipe)
			return
		}
	}
	pipe.Type(ctx, full)
	// rt, err := client.Type(ctx, full).Result()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	new_node := &Node{Value: key[0], FullKey: full, RedisType: "string"}
	// expand if the full key contains the search string
	if search_string != "" && strings.Contains(full, search_string) {
		new_node.Expanded = true
	}
	new_node.AddChild(key[1:], full, client, search_string, pipe)
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
