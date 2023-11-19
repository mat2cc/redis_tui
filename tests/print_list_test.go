package test

import (
	"testing"

	"github.com/mat2cc/redis_tui/tui"
)

func TestPrintList(t *testing.T) {
	node := &tui.Node{
		Value: "",
		Children: []*tui.Node{
			{Value: "foo", Expanded: true, Children: []*tui.Node{
				{Value: "bar", Expanded: true, Children: []*tui.Node{
					{Value: "bat", RedisType: "string"},
					{Value: "baz", RedisType: "string"},
				}},
			}},
			{Value: "test", RedisType: "string"},
			{Value: "testing", Expanded: true, Children: []*tui.Node{
				{Value: "123", RedisType: "string"},
			}},
			{Value: "notexpanded", Expanded: false, Children: []*tui.Node{
				{Value: "321", RedisType: "string"},
			}},
		},
	}

  tpl := tui.TablePrintList{
    List: tui.GeneratePrintList(node, 0),
  }
  rows :=tpl.GetRows()

  shouldEqual := [][]string {
    {"v", "foo"},
    {"v", "  bar"},
    {"",  "    bat [string]"},
    {"",  "    baz [string]"},
    {"",  "test [string]"},
    {"v", "testing"},
    {"",  "  123 [string]"},
    {">", "notexpanded"},
  }

  // comparing the rows, first index is prefix, second is value (including indent and redis type)
  for i, row := range rows {
    if row[0] != shouldEqual[i][0] || row[1] != shouldEqual[i][1] {
      t.Errorf("Expected %v, got %v", shouldEqual[i], row)
    }
  }
}
