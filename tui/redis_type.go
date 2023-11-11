package tui

import (
	"encoding/json"
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
    PrettyPrintJson bool
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
	for i, data := range *arr {
		out += fmt.Sprintf("%v\t%s\n", i+1, data)
	}
	return out
}

func (rs *RedisString) Print(table_width int) string {
	if rs.PrettyPrintJson && json.Valid([]byte(rs.Data)) {
		var data interface{}
		json.Unmarshal([]byte(rs.Data), &data)
		out, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			return rs.Data
		}
		return string(out)
	}
    return rs.Data
}

func (rh *RedisHash) Print(table_width int) string {
	out := ""
	for key, data := range rh.Data {
		out += fmt.Sprintf("%s:\t%s\n", key, data)
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
