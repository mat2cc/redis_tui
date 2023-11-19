package test

import (
	"bytes"
	"context"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mat2cc/redis_tui/tui"
)

func TestScan(t *testing.T) {
	client, err := tui.CreateRedisClient("", "", "", 2)
	client.Set(context.Background(), "foo", "foo", time.Second*10)
	client.Set(context.Background(), "foo:bar", "bar", time.Second*10)
	client.Set(context.Background(), "foo:bar:baz", "baz", time.Second*10)
	client.Set(context.Background(), "foo:bar:bat", "bat", time.Second*10)
	client.Set(context.Background(), "test", "test", time.Second*10)
	client.Set(context.Background(), "testing:123", "123", time.Second*10)

	node := tui.Node{
		Value: "",
		Children: []*tui.Node{
			{Value: "foo", Children: []*tui.Node{
				{Value: "bar", Children: []*tui.Node{
					{Value: "bat"},
					{Value: "baz"},
				}},
			}},
			{Value: "test"},
			{Value: "testing", Children: []*tui.Node{
				{Value: "123"},
			}},
		},
	}

	if err != nil {
		t.Error(err)
	}
  var buf bytes.Buffer

	model := tui.InitialModel(client, 10, true)
	p := tea.NewProgram(
		model,
    tea.WithOutput(&buf),
	)
	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			if model.Node.Children != nil {
				p.Quit()
				return
			}
		}
	}()

	_, err = p.Run()
	if err != nil {
		t.Fatal(err)
	}
	if !recursivelyCompareTrees(model.Node, node) {
		t.Error("Nodes are not equal")
	}
}

func recursivelyCompareTrees(a, b tui.Node) bool {
	if a.Value != b.Value {
		return false
	}

	if len(a.Children) != len(b.Children) {
		return false
	}

	if len(a.Children) == 0 && len(b.Children) == 0 {
		return true
	}

	// scan results are not ordered, so this is necessary
	for i := range a.Children {
		for j := range b.Children {
			if a.Children[i].Value == b.Children[j].Value {
				return recursivelyCompareTrees(*a.Children[i], *b.Children[j])
			}
		}
	}

	return false
}
