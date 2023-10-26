package redis_type

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	"github.com/redis/go-redis/v9"
)

type RT string

const (
	STRING = "string"
	LIST   = "list"
	SET    = "set"
	ZSET   = "zset"
	HASH   = "hash"
	STREAM = "stream"
)

type RedisType interface {
	Print(table_width int) string
}

type RedisString struct {
	RedisType RT
	Data      string
}

type RedisHash struct {
	RedisType RT
	Data      map[string]string
}

type RedisList struct {
	RedisType RT
	Data      []string
}

type RedisSet struct {
	RedisType RT
	Data      []string
}

type RedisZSet struct {
	RedisType RT
	Data      []string
}

type RedisStream struct {
	RedisType RT
	Data      []redis.XMessage
}

func tableStyles() table.Styles {
	s := table.DefaultStyles()
	s.Selected = s.Selected.Foreground(lipgloss.NoColor{})
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true)
	return s
}

func StringArrOut(arr *[]string, width int) string {
	out := ""
	var rows []table.Row
	for i, data := range *arr {
		rows = append(rows, table.Row{
			fmt.Sprint(i + 1),
			data,
		})
		out += fmt.Sprintf("%s\n", data)
	}
	t := table.New(
		table.WithColumns(
			[]table.Column{
				{Title: "", Width: 4},
				{Title: "Data", Width: width - 8},
			},
		),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithStyles(tableStyles()),
	)
	return t.View()
}

func (rs *RedisString) Print(table_width int) string {
	return rs.Data
}

func (rh *RedisHash) Print(table_width int) string {
	out := ""
	for key, data := range rh.Data {
		out += fmt.Sprintf("%s: %s\n", key, data)
	}
	return out
}

func (rl *RedisList) Print(table_width int) string {
	return StringArrOut(&rl.Data, table_width)
}

func (rs *RedisSet) Print(table_width int) string {
	return StringArrOut(&rs.Data, table_width)
}

func (gzs *RedisZSet) Print(table_width int) string {
	return StringArrOut(&gzs.Data, table_width)
}

func (rh *RedisStream) Print(table_width int) string {
	out := ""
	for _, data := range rh.Data {
		out += fmt.Sprintf("%s: %s\n", data.ID, data.Values)
	}
	return out
}
